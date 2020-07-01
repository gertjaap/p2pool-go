package work

import (
	"sync"

	"github.com/gertjaap/p2pool-go/wire"
)

type ShareChain struct {
	SharesChannel chan []wire.Share
	Tip           *wire.Share

	disconnectedShares    []*wire.Share
	disconnectedShareLock sync.Mutex
}

func NewShareChain() *ShareChain {
	sc := &ShareChain{disconnectedShares: make([]*wire.Share, 0), disconnectedShareLock: sync.Mutex{}, SharesChannel: make(chan []wire.Share, 10)}
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
