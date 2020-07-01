package wire

import (
	"fmt"
	"net"

	p2pnet "github.com/gertjaap/p2pool-go/net"
)

func NewP2PoolClient(host string, network p2pnet.Network) (*P2PoolConnection, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, network.P2PPort))
	if err != nil {
		return nil, err
	}
	return NewP2PoolConnection(conn, network), nil
}
