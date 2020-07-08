package main

import (
	"time"

	"github.com/gertjaap/p2pool-go/logging"
	p2pnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/p2p"
	"github.com/gertjaap/p2pool-go/work"
)

func main() {
	logging.SetLogLevel(int(logging.LogLevelDebug))
	p2pnet.ActiveNetwork = p2pnet.Vertcoin()

	sc := work.NewShareChain()
	err := sc.Load()
	if err != nil {
		panic(err)
	}

	//return
	pm := p2p.NewPeerManager(p2pnet.ActiveNetwork, sc)

	go func() {
		for s := range sc.NeedShareChannel {
			pm.AskForShare(s)
		}
	}()

	for {
		logging.Debugf("Number of active peers: %d", pm.GetPeerCount())
		time.Sleep(time.Second * 5)
	}
}
