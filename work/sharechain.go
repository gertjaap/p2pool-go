package work

import (
	"sync"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/gertjaap/p2pool-go/wire"
)

type ShareChain struct {
	SharesChannel    chan []wire.Share
	NeedShareChannel chan *chainhash.Hash
	Tip              *wire.Share
	AllShares        map[string]*wire.Share

	disconnectedShares    []*wire.Share
	disconnectedShareLock sync.Mutex
	allSharesLock         sync.Mutex
}

type Share struct {
	PreviousShare *Share
}

func NewShareChain() *ShareChain {
	sc := &ShareChain{disconnectedShares: make([]*wire.Share, 0), allSharesLock: sync.Mutex{}, disconnectedShareLock: sync.Mutex{}, SharesChannel: make(chan []wire.Share, 10)}
	go sc.ReadShareChan()
	return sc
}

func (sc *ShareChain) ReadShareChan() {
	for s := range sc.SharesChannel {
		sc.AddShares(s)
	}
}

func (sc *ShareChain) Resolve() {
	if len(sc.disconnectedShares) == 0 {
		return
	}

	extended := false
	for {
		if !extended {
			break
		}
	}
}

func (sc *ShareChain) AddShares(s []wire.Share) {
	// Decode

	sc.disconnectedShareLock.Lock()
	for _, sh := range s {
		sc.disconnectedShares = append(sc.disconnectedShares, &sh)
	}
	sc.disconnectedShareLock.Unlock()

	sc.Resolve()
}
