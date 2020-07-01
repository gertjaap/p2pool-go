package wire

import (
	"bytes"
	"encoding/binary"
)

var _ P2PoolMessage = &MsgGetAddrs{}

type MsgGetAddrs struct {
	Count int32
}

func (m *MsgGetAddrs) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	err := binary.Read(r, binary.LittleEndian, &m.Count)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsgGetAddrs) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, m.Count)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *MsgGetAddrs) Command() string {
	return "getaddrs"
}
