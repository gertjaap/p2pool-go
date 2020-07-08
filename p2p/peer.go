package p2p

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	p2poolnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/util"
	"github.com/gertjaap/p2pool-go/wire"
)

type Peer struct {
	Connection *wire.P2PoolConnection
	RemoteIP   net.IP
	RemotePort int
	Network    p2poolnet.Network

	newPeers    chan []wire.Addr
	sharesChan  chan []wire.Share
	versionInfo *wire.MsgVersion
}

func NewPeer(ip net.IP, port int, n p2poolnet.Network, newPeers chan []wire.Addr, closed chan bool, sharesChan chan []wire.Share) (*Peer, error) {
	p := Peer{Network: n, newPeers: newPeers, sharesChan: sharesChan}
	p.RemoteIP = ip
	var err error
	p.Connection, err = wire.NewP2PoolClient(ip, port, n)
	if err != nil {
		return nil, err
	}

	err = p.Handshake()
	if err != nil {
		p.Connection.Close()
		return nil, err
	}

	go func() {
		<-p.Connection.Disconnected
		closed <- true
	}()

	go p.IncomingLoop()
	go p.PingLoop()

	return &p, nil
}

func (p *Peer) BestShare() *chainhash.Hash {
	return p.versionInfo.BestShareHash
}

func (p *Peer) PingLoop() {
	for {
		time.Sleep(time.Second * 15)
		p.Connection.Outgoing <- &wire.MsgPing{}
	}
}

func (p *Peer) IncomingLoop() {
	for msg := range p.Connection.Incoming {
		switch t := msg.(type) {
		case *wire.MsgAddrs:
			p.newPeers <- t.Addresses
		case *wire.MsgShares:
			p.sharesChan <- t.Shares
		case *wire.MsgShareReply:
			p.sharesChan <- t.Shares
		}
	}
}

func (p *Peer) AskNewAddresses(count int32) {
	p.Connection.Outgoing <- &wire.MsgGetAddrs{
		Count: count,
	}
}

func (p *Peer) Handshake() error {
	myIP, err := util.GetMyPublicIP()
	if err != nil {
		panic(err)
	}
	p.Connection.Outgoing <- &wire.MsgVersion{
		Version:  1800,
		Services: 0,
		AddrTo: wire.P2PoolAddress{
			Services: 0,
			Address:  p.RemoteIP,
			Port:     int16(p.Network.P2PPort),
		},
		AddrFrom: wire.P2PoolAddress{
			Services: 0,
			Address:  myIP,
			Port:     int16(p.Network.P2PPort),
		},
		Nonce:      int64(rand.Uint64()),
		SubVersion: "p2pool-go/0.0.1",
		Mode:       1,
	}
	select {
	case msg := <-p.Connection.Incoming:
		var ok bool
		p.versionInfo, ok = msg.(*wire.MsgVersion)
		if !ok {
			return fmt.Errorf("First message received from peer was not version message")
		}
	case <-time.After(5 * time.Second):
		return fmt.Errorf("Timeout waiting for version message from peer")
	}
	return nil
}
