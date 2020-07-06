package main

import (
	"encoding/hex"
	"time"

	"github.com/gertjaap/p2pool-go/logging"
	p2pnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/p2p"
	"github.com/gertjaap/p2pool-go/wire"
	"github.com/gertjaap/p2pool-go/work"
)

func main() {
	logging.SetLogLevel(int(logging.LogLevelDebug))

	m := &wire.MsgShares{}
	b, _ := hex.DecodeString("0111fd0b01fe00000020eb400edde9c93272c3e3c74c58ca47ee4280fc68a4a6af658d270394978b7d55dc13035fbaae001b1029d544c58ddfe53843917f5491185695d0a6143f2fd5e364141557df971a164c1684c60403020f1510f6db71531e9692d87cf03c0340e3b0491150d752c2d5d44700f90295000000008f0200110000000000000000000000000000000000000000000000000000000000000000000000381419829e6b838121edf10fd958f5354cc03b960dc385f09b28f1a565891ab6115b051dc5b32d1ce513035ff2b52d0093ae67bd4f06080500000000000000000000000000b6f00000cf4dbba0589d15681e43b17de1fd2079f3e8d58e2336c81b563575bfc268e6adfdad0200")
	err := m.FromBytes(b)
	if err != nil {
		panic(err)
	}

	sc := work.NewShareChain()
	pm := p2p.NewPeerManager(p2pnet.Vertcoin(), sc.SharesChannel)
	pm.AskForShare(m.Shares[0].ShareInfo.ShareData.PreviousShareHash)
	for {
		logging.Debugf("Number of active peers: %d", pm.GetPeerCount())
		time.Sleep(time.Second * 5)
	}
}
