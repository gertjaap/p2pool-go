package wire

import (
	"bytes"
	"encoding/binary"
)

var _ P2PoolMessage = &MsgAddrMe{}

type MsgAddrMe struct {
	Port int16
}

func (m *MsgAddrMe) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	err := binary.Read(r, binary.LittleEndian, &m.Port)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsgAddrMe) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, m.Port)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *MsgAddrMe) Command() string {
	return "addrme"
}
