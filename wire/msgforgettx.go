package wire

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var _ P2PoolMessage = &MsgForgetTx{}

type MsgForgetTx struct {
	TXHashes []*chainhash.Hash
}

func (m *MsgForgetTx) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	m.TXHashes = make([]*chainhash.Hash, 0)
	count, err := ReadVarInt(r)
	if err != nil {
		return err
	}

	for i := uint64(0); i < count; i++ {
		h, err := ReadChainHash(r)
		if err != nil {
			return err
		}
		m.TXHashes = append(m.TXHashes, h)
	}
	return nil
}

func (m *MsgForgetTx) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := WriteVarInt(&buf, uint64(len(m.TXHashes)))
	if err != nil {
		return nil, err
	}
	for _, h := range m.TXHashes {
		err = WriteChainHash(&buf, h)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (m *MsgForgetTx) Command() string {
	return "forget_tx"
}
