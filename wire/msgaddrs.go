package wire

import (
	"bytes"
	"encoding/binary"
)

var _ P2PoolMessage = &MsgAddrs{}

type MsgAddrs struct {
	Addresses []Addr
}

type Addr struct {
	Timestamp int64
	Address   P2PoolAddress
}

func (m *MsgAddrs) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	m.Addresses = make([]Addr, 0)
	count, err := ReadVarInt(r)
	if err != nil {
		return err
	}

	for i := uint64(0); i < count; i++ {
		var a Addr
		err = binary.Read(r, binary.LittleEndian, &a.Timestamp)
		if err != nil {
			return err
		}

		a.Address, err = ReadP2PoolAddress(r)
		if err != nil {
			return err
		}

		m.Addresses = append(m.Addresses, a)
	}
	return nil
}

func (m *MsgAddrs) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := WriteVarInt(&buf, uint64(len(m.Addresses)))
	if err != nil {
		return nil, err
	}
	for _, a := range m.Addresses {
		err = binary.Write(&buf, binary.LittleEndian, a.Timestamp)
		if err != nil {
			return nil, err
		}

		b, err := a.Address.ToBytes()
		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (m *MsgAddrs) Command() string {
	return "addrs"
}
