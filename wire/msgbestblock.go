package wire

import (
	"bytes"

	"github.com/btcsuite/btcd/wire"
)

var _ P2PoolMessage = &MsgBestBlock{}

type MsgBestBlock struct {
	BestBlock *wire.BlockHeader
}

func (m *MsgBestBlock) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	m.BestBlock = wire.NewBlockHeader(0, nullHash, nullHash, 0, 0)
	return m.BestBlock.Deserialize(r)
}

func (m *MsgBestBlock) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	err := m.BestBlock.Serialize(&buf)
	return buf.Bytes(), err
}

func (m *MsgBestBlock) Command() string {
	return "bestblock"
}
