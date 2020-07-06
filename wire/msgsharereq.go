package wire

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var _ P2PoolMessage = &MsgShareReq{}

type MsgShareReq struct {
	ID      *chainhash.Hash
	Hashes  []*chainhash.Hash
	Parents uint64
	Stops   []*chainhash.Hash
}

func (m *MsgShareReq) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	var err error
	m.ID, err = ReadChainHash(r)
	if err != nil {
		return err
	}
	m.Hashes, err = ReadChainHashList(r)
	if err != nil {
		return err
	}
	m.Parents, err = ReadVarInt(r)
	if err != nil {
		return err
	}
	m.Stops, err = ReadChainHashList(r)
	if err != nil {
		return err
	}

	return nil
}

func (m *MsgShareReq) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	var err error

	err = WriteChainHash(&buf, m.ID)
	if err != nil {
		return nil, err
	}
	err = WriteChainHashList(&buf, m.Hashes)
	if err != nil {
		return nil, err
	}
	err = WriteVarInt(&buf, m.Parents)
	if err != nil {
		return nil, err
	}
	err = WriteChainHashList(&buf, m.Stops)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *MsgShareReq) Command() string {
	return "sharereq"
}
