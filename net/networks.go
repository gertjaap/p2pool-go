package net

import (
	"encoding/hex"

	"github.com/adamcollier1/lyra2rev3"
)

var ActiveNetwork Network

type Network struct {
	MessagePrefix []byte
	Identifier    []byte
	P2PPort       int
	SeedHosts     []string
	ChainLength   int
	POWHash       func([]byte) []byte
}

func Vertcoin() Network {
	n := Network{P2PPort: 9346}
	n.MessagePrefix, _ = hex.DecodeString("7c3614a6bcdcf784")
	n.Identifier, _ = hex.DecodeString("a06a81c827cab983")
	n.ChainLength = 5100
	n.SeedHosts = []string{"localhost", "p2proxy.vertcoin.org", "vtc.alwayshashing.com", "crypto.office-on-the.net", "pool.vtconline.org"}
	n.POWHash = func(b []byte) []byte {
		res, _ := lyra2rev3.SumV3(b)
		return res
	}
	return n
}
