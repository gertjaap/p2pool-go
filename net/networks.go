package net

import "encoding/hex"

type Network struct {
	MessagePrefix []byte
	P2PPort       int
	SeedHosts     []string
}

func Vertcoin() Network {
	n := Network{P2PPort: 9346}
	n.MessagePrefix, _ = hex.DecodeString("7c3614a6bcdcf784")
	n.SeedHosts = []string{"localhost", "p2proxy.vertcoin.org", "vtc.alwayshashing.com", "crypto.office-on-the.net", "pool.vtconline.org"}
	return n
}
