package wire

var _ P2PoolMessage = &MsgPing{}

type MsgPing struct {
}

func (m *MsgPing) FromBytes(b []byte) error {
	return nil
}

func (m *MsgPing) ToBytes() ([]byte, error) {
	return []byte{}, nil
}

func (m *MsgPing) Command() string {
	return "ping"
}
