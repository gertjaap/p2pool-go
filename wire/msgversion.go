package wire

import (
	"bytes"
	"encoding/binary"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var _ P2PoolMessage = &MsgVersion{}

type MsgVersion struct {
	Version       int32
	Services      int64
	AddrTo        P2PoolAddress
	AddrFrom      P2PoolAddress
	Nonce         int64
	SubVersion    string
	Mode          int32
	BestShareHash *chainhash.Hash
}

func (m *MsgVersion) FromBytes(b []byte) error {
	buf := bytes.NewReader(b)

	err := binary.Read(buf, binary.LittleEndian, &m.Version)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &m.Services)
	if err != nil {
		return err
	}

	m.AddrTo, err = ReadP2PoolAddress(buf)
	if err != nil {
		return err
	}

	m.AddrFrom, err = ReadP2PoolAddress(buf)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &m.Nonce)
	if err != nil {
		return err
	}

	m.SubVersion, err = ReadVarString(buf)
	if err != nil {
		return err
	}

	err = binary.Read(buf, binary.LittleEndian, &m.Mode)
	if err != nil {
		return err
	}
	m.BestShareHash, err = ReadChainHash(buf)

	return nil
}

func (m *MsgVersion) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, m.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, m.Services)
	if err != nil {
		return nil, err
	}
	addr, err := m.AddrTo.ToBytes()
	if err != nil {
		return nil, err
	}
	buf.Write(addr)
	addr, err = m.AddrFrom.ToBytes()
	if err != nil {
		return nil, err
	}
	buf.Write(addr)
	err = binary.Write(&buf, binary.LittleEndian, m.Nonce)
	if err != nil {
		return nil, err
	}
	err = WriteVarString(&buf, m.SubVersion)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, m.Mode)
	if err != nil {
		return nil, err
	}
	err = WriteChainHash(&buf, m.BestShareHash)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *MsgVersion) Command() string {
	return "version"
}
