package net

import "encoding/hex"

type Network struct {
	MessagePrefix []byte
	P2PPort       int
}

func Vertcoin() Network {
	n := Network{P2PPort: 9346}
	n.MessagePrefix, _ = hex.DecodeString("7c3614a6bcdcf784")
	return n
}
