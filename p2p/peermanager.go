package p2p

import (
	"net"
	"sync"
	"time"

	"github.com/gertjaap/p2pool-go/util"
	"github.com/gertjaap/p2pool-go/work"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gertjaap/p2pool-go/logging"
	p2poolnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/wire"
)

type PeerManager struct {
	Network           p2poolnet.Network
	peers             []*Peer
	possiblePeers     []wire.Addr
	shareChain        *work.ShareChain
	askSharesChan     chan *chainhash.Hash
	peersLock         sync.Mutex
	possiblePeersLock sync.Mutex
}

func NewPeerManager(n p2poolnet.Network, sc *work.ShareChain) *PeerManager {
	p := &PeerManager{
		Network:           n,
		peers:             make([]*Peer, 0),
		possiblePeers:     make([]wire.Addr, 0),
		peersLock:         sync.Mutex{},
		possiblePeersLock: sync.Mutex{},
		shareChain:        sc,
		askSharesChan:     make(chan *chainhash.Hash, 100),
	}

	for _, h := range n.SeedHosts {
		addrs, err := net.LookupIP(h)
		if err == nil {
			a := wire.Addr{
				Address: wire.P2PoolAddress{
					Address: addrs[0],
					Port:    int16(n.P2PPort),
				},
			}
			p.possiblePeers = append(p.possiblePeers, a)
		}
	}
	go p.MonitorPeerCount()
	go p.ShareAskLoop()
	return p
}

func (p *PeerManager) MonitorPeerCount() {
	for {
		if len(p.peers) > 0 {
			time.Sleep(time.Second * 10)
		}
		for len(p.peers) < 1 {
			tryPeer := p.GetPossiblePeer()
			if tryPeer.Timestamp == -1 {
				logging.Debugf("Not enough peers, and no possible peers to try. Asking existing peers for new peers")
				// No peers left to try. Ask for more.
				for _, peer := range p.peers {
					peer.AskNewAddresses(10)
				}
				break
			}
			peerAddress := tryPeer.Address.Address
			logging.Debugf("Trying peer %s", peerAddress.String())

			err := p.AddPeerWithPort(peerAddress, int(tryPeer.Address.Port))
			if err != nil {
				logging.Warnf("Peer %s failed: %s", peerAddress.String(), err.Error())
				p.RemovePossiblePeer(tryPeer)
			}
		}
	}
}

func (p *PeerManager) ShareAskLoop() {
	for {
		if len(p.peers) > 0 {
			h := <-p.askSharesChan
			for _, pr := range p.peers {
				stops := make([]*chainhash.Hash, 0)
				tip := p.shareChain.GetTipHash()
				if tip != nil {
					stops = append(stops, tip)
				}
				pr.Connection.Outgoing <- &wire.MsgShareReq{
					ID:      util.GetRandomId(),
					Parents: 1000,
					Stops:   stops,
					Hashes:  []*chainhash.Hash{h},
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func (p *PeerManager) AskForShare(h *chainhash.Hash) {
	p.askSharesChan <- h
}

func (p *PeerManager) GetPossiblePeer() wire.Addr {
	for _, pos := range p.possiblePeers {
		alreadyAPeer := false
		for _, pr := range p.peers {
			if pr.RemoteIP.String() == pos.Address.Address.String() {
				alreadyAPeer = true
				break
			}
		}
		if !alreadyAPeer {
			return pos
		}
	}
	return wire.Addr{Timestamp: -1}
}

func (p *PeerManager) RemovePossiblePeer(addr wire.Addr) {
	p.possiblePeersLock.Lock()
	newPossiblePeers := make([]wire.Addr, 0)
	for _, p := range p.possiblePeers {
		if p.Address.Address.String() != addr.Address.Address.String() {
			newPossiblePeers = append(newPossiblePeers, p)
		}
	}
	p.possiblePeers = newPossiblePeers
	p.possiblePeersLock.Unlock()
}

func (p *PeerManager) AddPeer(ip net.IP) error {
	return p.AddPeerWithPort(ip, 0)
}

func (p *PeerManager) AddPeerWithPort(ip net.IP, port int) error {
	newPeers := make(chan []wire.Addr, 10)
	closed := make(chan bool, 1)
	peer, err := NewPeer(ip, port, p.Network, newPeers, closed, p.shareChain.SharesChannel)
	if err != nil {
		return err
	}
	p.peersLock.Lock()
	p.peers = append(p.peers, peer)
	p.peersLock.Unlock()

	stops := make([]*chainhash.Hash, 0)
	tip := p.shareChain.GetTipHash()
	skipAsk := false
	if tip != nil {
		stops = append(stops, tip)
		if tip.IsEqual(peer.versionInfo.BestShareHash) {
			skipAsk = true
		}
	}

	if !skipAsk {
		peer.Connection.Outgoing <- &wire.MsgShareReq{
			ID:      util.GetRandomId(),
			Parents: 1000,
			Stops:   stops,
			Hashes:  []*chainhash.Hash{peer.versionInfo.BestShareHash},
		}
	}

	go p.NewPeersHandler(newPeers)
	go p.ClosedHandler(peer, closed)
	return nil
}

func (p *PeerManager) NewPeersHandler(c chan []wire.Addr) {
	for a := range c {
		p.possiblePeersLock.Lock()
		p.possiblePeers = append(p.possiblePeers, a...)
		p.possiblePeersLock.Unlock()
	}

}

func (p *PeerManager) ClosedHandler(peer *Peer, c chan bool) {
	<-c
	p.peersLock.Lock()
	newPeers := make([]*Peer, 0)
	for _, p := range p.peers {
		if p != peer {
			newPeers = append(newPeers, p)
		}
	}
	p.peers = newPeers
	p.peersLock.Unlock()
}

func (p *PeerManager) GetPeerCount() int {
	return len(p.peers)
}
