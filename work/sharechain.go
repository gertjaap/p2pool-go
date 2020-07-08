package work

import (
	"sync"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gertjaap/p2pool-go/logging"
	"github.com/gertjaap/p2pool-go/wire"
)

type ShareChain struct {
	SharesChannel    chan []wire.Share
	NeedShareChannel chan *chainhash.Hash
	Tip              *ChainShare
	AllShares        map[string]*ChainShare

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
	sc := &ShareChain{disconnectedShares: make([]*wire.Share, 0), allSharesLock: sync.Mutex{}, AllShares: map[string]*ChainShare{}, disconnectedShareLock: sync.Mutex{}, SharesChannel: make(chan []wire.Share, 10)}
	go sc.ReadShareChan()
	return sc
}

func (sc *ShareChain) ReadShareChan() {
	for s := range sc.SharesChannel {
		sc.AddShares(s)
	}
}

func (sc *ShareChain) Resolve() {
	logging.Debugf("Resolving sharechain")
	if len(sc.disconnectedShares) == 0 {
		return
	}

	if sc.Tip == nil {
		logging.Debugf("Setting tip to first disconnected share")
		sc.disconnectedShareLock.Lock()
		newChainShare := &ChainShare{Share: sc.disconnectedShares[0]}
		sc.Tip = newChainShare
		sc.disconnectedShares = sc.disconnectedShares[1:]
		sc.disconnectedShareLock.Unlock()

		sc.allSharesLock.Lock()
		sc.AllShares[newChainShare.Share.Hash.String()] = newChainShare
		sc.allSharesLock.Unlock()
	}

	for {
		extended := false
		sc.disconnectedShareLock.Lock()
		newDisconnectedShares := make([]*wire.Share, 0)
		for _, s := range sc.disconnectedShares {
			es, ok := sc.AllShares[s.ShareInfo.ShareData.PreviousShareHash.String()]
			if ok {
				logging.Debugf("Found connecting share after existing one")
				newChainShare := &ChainShare{Share: s, Previous: es}
				es.Next = newChainShare
				if es.Share.Hash.IsEqual(sc.Tip.Share.Hash) {
					sc.Tip = newChainShare
				}
				sc.allSharesLock.Lock()
				sc.AllShares[newChainShare.Share.Hash.String()] = newChainShare
				sc.allSharesLock.Unlock()
				extended = true
			} else {
				newDisconnectedShares = append(newDisconnectedShares, s)
			}
		}

		sc.disconnectedShares = newDisconnectedShares

		usedIndices := make([]int, 0)
		for _, es := range sc.AllShares {
			for i, s := range sc.disconnectedShares {
				if s.Hash.IsEqual(es.Share.ShareInfo.ShareData.PreviousShareHash) {
					logging.Debugf("Found connecting share before existing one")
					newChainShare := &ChainShare{Share: s, Next: es}
					es.Previous = newChainShare
					sc.allSharesLock.Lock()
					sc.AllShares[newChainShare.Share.Hash.String()] = newChainShare
					sc.allSharesLock.Unlock()
					extended = true
					usedIndices = append(usedIndices, i)
				}
			}
		}

		newDisconnectedShares = make([]*wire.Share, 0)
		for i, s := range sc.disconnectedShares {
			found := false
			for _, idx := range usedIndices {
				if i == idx {
					found = true
				}
			}
			if !found {
				newDisconnectedShares = append(newDisconnectedShares, s)
			}
		}

		sc.disconnectedShares = newDisconnectedShares

		logging.Debugf("Tip is now %s - disconnected: %d", sc.Tip.Share.Hash.String(), len(sc.disconnectedShares))

		sc.disconnectedShareLock.Unlock()
		if !extended || len(sc.disconnectedShares) == 0 {
			break
		}
	}

	if len(sc.disconnectedShares) > 0 {
		logging.Debugf("Shares in the chain:")
		for _, s := range sc.AllShares {
			logging.Debugf("H: %s P: %s Hght: %d", s.Share.Hash.String()[56:], s.Share.ShareInfo.ShareData.PreviousShareHash.String()[56:], s.Share.ShareInfo.AbsHeight)
		}

		logging.Debugf("Still have %d unconnected shares:", len(sc.disconnectedShares))
		for _, s := range sc.disconnectedShares {
			logging.Debugf("H: %s P: %s Hght: %d", s.Hash.String()[56:], s.ShareInfo.ShareData.PreviousShareHash.String()[56:], s.ShareInfo.AbsHeight)
		}
	}
}

func (sc *ShareChain) AddShares(s []wire.Share) {
	// Decode

	sc.disconnectedShareLock.Lock()
	for i := range s {
		sc.disconnectedShares = append(sc.disconnectedShares, &s[i])
	}
	sc.disconnectedShareLock.Unlock()

	sc.Resolve()
}

func (sc *ShareChain) GetTipHash() *chainhash.Hash {
	if sc.Tip != nil {
		return sc.Tip.Share.Hash
	}
	return nil
}
