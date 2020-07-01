package wire

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btcwire "github.com/btcsuite/btcd/wire"
)

var _ P2PoolMessage = &MsgHaveTx{}

type MsgRememberTx struct {
	TXHashes []*chainhash.Hash
	TXs      []*btcwire.MsgTx
}

func (m *MsgRememberTx) FromBytes(b []byte) error {
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

	count, err = ReadVarInt(r)
	if err != nil {
		return err
	}

	for i := uint64(0); i < count; i++ {
		tx := btcwire.NewMsgTx(1)
		err = tx.Deserialize(r)
		if err != nil {
			return err
		}
		m.TXs = append(m.TXs, tx)
	}

	return nil
}

func (m *MsgRememberTx) ToBytes() ([]byte, error) {
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
	err = WriteVarInt(&buf, uint64(len(m.TXs)))
	if err != nil {
		return nil, err
	}
	for _, tx := range m.TXs {
		err = tx.Serialize(&buf)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (m *MsgRememberTx) Command() string {
	return "remember_tx"
}
