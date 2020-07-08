package main

import (
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/gertjaap/p2pool-go/logging"
	p2pnet "github.com/gertjaap/p2pool-go/net"
	"github.com/gertjaap/p2pool-go/p2p"
	"github.com/gertjaap/p2pool-go/wire"
	"github.com/gertjaap/p2pool-go/work"
)

func main() {
	logging.SetLogLevel(int(logging.LogLevelDebug))

	m := &wire.MsgShareReply{}
	b, _ := hex.DecodeString("1d729566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e001679000111fd4d01fe00000020418674e935503fc5c9a90052a080b1d48994f66a9274ac0002c9f0f1fdbd235ad9de045f60ae001b0b0f2d17ee1437ab8a21126d8cd9ccc791bde69fd6cc29cd41c0711d32aa188f891287f40403161215d1a09b94531e9692d87cf03c0340e3b0491150d752c2d5d447c85d0395000000008f020011015355e4eb7b4526e6d54cb954f2d013e2d041cc0b4de0245bb3c305e7733f1bc06faccb80657760be3660e48fdff22058b72ba1ad30e080a2806e6d1c82406e88000102000b99b3abc583ae5ff76a5df1e82b106a0bd1fa199c6a019c53d2d6b46a7b76e7d6ce061d36183a1ce1de045f46d42d00866a726f1ba90805000000000000000000000000000823010047985fd8cd08f7cf62167bd55afff7bbb64f5d166001cd9a367f7282dc44e78cfdad02011fbf5e8a4effa65ebbeb7b3b6c9aac7f6e9e5e7f9ba3a2b55fbdebeba44185ac")
	err := m.FromBytes(b)
	if err != nil {
		panic(err)
	}

	/*

		previous_share_hash:  f48712898f18aa321d71c041cd29ccd69fe6bd91c7ccd98c6d12218aab3714ee -
		GenTXHash:  a89c94f17b50c6732e7c1f5227110efd4f5566bd208633117924fc77e61a70f8 -
		MerkleRoot:  cbadc884ebfd2814bad396f2f4955a349a9f5230bac495b7aab48bb6760b727b -
		Hash:  3433aee492c68904de82d5d0f8a8e375e44a5a872f009eefb529743dc8726061

	*/

	logging.Debugf("previous_share_hash:  %s - GenTXHash:  %s - MerkleRoot:  %s - Hash:  %s", m.Shares[0].ShareInfo.ShareData.PreviousShareHash.String(), m.Shares[0].GenTXHash.String(), m.Shares[0].MerkleRoot.String(), m.Shares[0].Hash())

	return

	sc := work.NewShareChain()
	p2pnet.ActiveNetwork = p2pnet.Vertcoin()
	pm := p2p.NewPeerManager(p2pnet.ActiveNetwork, sc.SharesChannel)

	h, _ := chainhash.NewHashFromStr("3433aee492c68904de82d5d0f8a8e375e44a5a872f009eefb529743dc8726061")
	pm.AskForShare(h)
	for {
		logging.Debugf("Number of active peers: %d", pm.GetPeerCount())
		time.Sleep(time.Second * 5)
	}
}
