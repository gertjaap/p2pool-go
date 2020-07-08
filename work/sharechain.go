package work

import (
	"fmt"
	"os"
	"sync"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gertjaap/p2pool-go/logging"
	p2pnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/wire"
)

type ShareChain struct {
	SharesChannel    chan []wire.Share
	NeedShareChannel chan *chainhash.Hash
	Tip              *ChainShare
	Tail             *ChainShare
	AllShares        map[string]*ChainShare
	AllSharesByPrev  map[string]*ChainShare

	disconnectedShares    []*wire.Share
	disconnectedShareLock sync.Mutex
	allSharesLock         sync.Mutex
}

type ChainShare struct {
	Share    *wire.Share
	Previous *ChainShare
	Next     *ChainShare
}

func NewShareChain() *ShareChain {
	sc := &ShareChain{disconnectedShares: make([]*wire.Share, 0), allSharesLock: sync.Mutex{}, AllSharesByPrev: map[string]*ChainShare{}, AllShares: map[string]*ChainShare{}, disconnectedShareLock: sync.Mutex{}, SharesChannel: make(chan []wire.Share, 10), NeedShareChannel: make(chan *chainhash.Hash, 10)}
	go sc.ReadShareChan()
	return sc
}

func (sc *ShareChain) ReadShareChan() {
	for s := range sc.SharesChannel {
		sc.AddShares(s)
	}
}

func (sc *ShareChain) AddChainShare(newChainShare *ChainShare) {
	sc.allSharesLock.Lock()
	sc.AllShares[newChainShare.Share.Hash.String()] = newChainShare
	sc.AllSharesByPrev[newChainShare.Share.ShareInfo.ShareData.PreviousShareHash.String()] = newChainShare
	sc.allSharesLock.Unlock()
}

func (sc *ShareChain) Resolve(skipCommit bool) {
	logging.Debugf("Resolving sharechain")
	if len(sc.disconnectedShares) == 0 {
		return
	}

	if sc.Tip == nil {
		sc.disconnectedShareLock.Lock()
		newChainShare := &ChainShare{Share: sc.disconnectedShares[0]}
		sc.Tip = newChainShare
		sc.disconnectedShares = sc.disconnectedShares[1:]
		sc.disconnectedShareLock.Unlock()
		sc.AddChainShare(newChainShare)
		sc.Tail = sc.Tip
	}

	for {
		extended := false
		sc.disconnectedShareLock.Lock()
		newDisconnectedShares := make([]*wire.Share, 0)
		for _, s := range sc.disconnectedShares {
			_, ok := sc.AllShares[s.Hash.String()]
			if ok {
				// Duplicate
				continue
			}

			es, ok := sc.AllShares[s.ShareInfo.ShareData.PreviousShareHash.String()]
			if ok {
				newChainShare := &ChainShare{Share: s, Previous: es}
				es.Next = newChainShare
				if es.Share.Hash.IsEqual(sc.Tip.Share.Hash) {
					sc.Tip = newChainShare
				}
				sc.AddChainShare(newChainShare)
				extended = true
			} else {
				es, ok := sc.AllSharesByPrev[s.Hash.String()]
				if ok {
					newChainShare := &ChainShare{Share: s, Next: es}
					es.Previous = newChainShare
					if es.Share.Hash.IsEqual(sc.Tail.Share.Hash) {
						sc.Tail = newChainShare
					}
					sc.AddChainShare(newChainShare)
					extended = true
				} else {
					newDisconnectedShares = append(newDisconnectedShares, s)
				}
			}
		}

		sc.disconnectedShares = newDisconnectedShares

		sc.disconnectedShareLock.Unlock()
		if !extended || len(sc.disconnectedShares) == 0 {
			break
		}
	}

	logging.Debugf("Tip is now %s - disconnected: %d - Length: %d", sc.Tip.Share.Hash.String(), len(sc.disconnectedShares), len(sc.AllShares))

	if len(sc.AllShares) < p2pnet.ActiveNetwork.ChainLength {
		sc.NeedShareChannel <- sc.Tail.Share.ShareInfo.ShareData.PreviousShareHash
	}
	if !skipCommit {
		sc.Commit()
	}
}

func (sc *ShareChain) Commit() error {
	sc.allSharesLock.Lock()

	shares := make([]wire.Share, 0)
	i := 0
	s := sc.Tip
	for s != nil {
		shares = append(shares, *(s.Share))
		s = s.Previous
		i++
	}
	f, err := os.Create("sharechain-new.dat")
	if err != nil {
		return err
	}

	wire.WriteShares(f, shares)

	f.Close()

	if _, err := os.Stat("sharechain.dat"); err == nil {
		err = os.Remove("sharechain.dat")
		if err != nil {
			return err
		}
	}

	os.Rename("sharechain-new.dat", "sharechain.dat")

	sc.allSharesLock.Unlock()
	return nil
}

func (sc *ShareChain) Load() error {

	if _, err := os.Stat("sharechain.dat"); os.IsNotExist(err) {
		return nil // Sharechain data absent, no need to do anything then.
	}

	f, err := os.Open("sharechain.dat")
	if err != nil {
		return err
	}
	shares, err := wire.ReadShares(f)
	if err != nil {
		return err
	}

	for _, s := range shares {
		if !s.IsValid() {
			return fmt.Errorf("Invalid share found")
		}
	}

	sc.disconnectedShareLock.Lock()
	sc.disconnectedShares = make([]*wire.Share, len(shares))
	for i := range shares {
		sc.disconnectedShares[i] = &shares[i]
	}
	sc.disconnectedShareLock.Unlock()

	logging.Debugf("Loaded %d shares from disk", len(sc.disconnectedShares))

	sc.Resolve(true)

	return nil
}

func (sc *ShareChain) AddShares(s []wire.Share) {
	// Decode

	sc.disconnectedShareLock.Lock()
	for i := range s {
		if s[i].IsValid() {
			_, ok := sc.AllShares[s[i].Hash.String()]
			if !ok {
				sc.disconnectedShares = append(sc.disconnectedShares, &s[i])
			}
		} else {
			logging.Warnf("Ignoring invalid share %s", s[i].Hash.String())
		}
	}
	sc.disconnectedShareLock.Unlock()

	sc.Resolve(false)
}

func (sc *ShareChain) GetTipHash() *chainhash.Hash {
	if sc.Tip != nil {
		return sc.Tip.Share.Hash
	}
	return nil
}
