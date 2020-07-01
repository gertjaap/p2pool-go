package wire

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type P2PoolAddress struct {
	Services int64
	Address  net.IP
	Port     int16
}

func (p P2PoolAddress) ToBytes() ([]byte, error) {
	var b []byte
	buf := bytes.NewBuffer(b)

	err := binary.Write(buf, binary.LittleEndian, p.Services)
	if err != nil {
		return nil, err
	}
	err = WriteIPAddr(buf, p.Address)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, p.Port)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ReadP2PoolAddress(r io.Reader) (P2PoolAddress, error) {
	a := P2PoolAddress{}
	err := binary.Read(r, binary.LittleEndian, &a.Services)
	if err != nil {
		return a, err
	}

	a.Address, err = ReadIPAddr(r)
	if err != nil {
		return a, err
	}

	err = binary.Read(r, binary.LittleEndian, &a.Port)
	if err != nil {
		return a, err
	}

	return a, nil
}

func P2PoolAddressFromBytes(b []byte) (P2PoolAddress, error) {
	r := bytes.NewReader(b)
	return ReadP2PoolAddress(r)
}
