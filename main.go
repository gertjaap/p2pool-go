package main

import (
	"math/rand"
	sysnet "net"

	"github.com/gertjaap/p2pool-go/logging"
	p2pnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/util"
	"github.com/gertjaap/p2pool-go/wire"
)

func main() {

	net := p2pnet.Vertcoin()

	logging.SetLogLevel(int(logging.LogLevelDebug))

	addrs, err := sysnet.LookupIP("p2proxy.vertcoin.org")
	if err != nil {
		panic(err)
	}

	c, err := wire.NewP2PoolClient("p2proxy.vertcoin.org", net)
	if err != nil {
		panic(err)
	}

	myIP, err := util.GetMyPublicIP()
	if err != nil {
		panic(err)
	}
	c.Outgoing <- &wire.MsgVersion{
		Version:  1800,
		Services: 0,
		AddrTo: wire.P2PoolAddress{
			Services: 0,
			Address:  addrs[0],
			Port:     int16(net.P2PPort),
		},
		AddrFrom: wire.P2PoolAddress{
			Services: 0,
			Address:  myIP,
			Port:     int16(net.P2PPort),
		},
		Nonce:      int64(rand.Uint64()),
		SubVersion: "0.7.2",
		Mode:       1,
	}

	c.Outgoing <- &wire.MsgGetAddrs{
		Count: 10,
	}

	for msg := range c.Incoming {
		logging.Debugf("Received incoming message [%s]", msg.Command())
		switch t := msg.(type) {
		case *wire.MsgVersion:
			logging.Debugf("Received version message - Version [%d] - Best Share [%s]", t.Version, t.BestShareHash.String())
		case *wire.MsgAddrs:
			logging.Debugf("Received addresses:")
			for _, a := range t.Addresses {
				logging.Debugf("Timestamp [%d] - IP [%s]", a.Timestamp, a.Address.Address.String())
			}
		}
	}
}
