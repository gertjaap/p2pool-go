// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var nullHash *chainhash.Hash

func ReadVarString(r io.Reader) (string, error) {
	len, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	b := make([]byte, len)
	rl, err := r.Read(b)
	if rl != int(len) {
		return "", fmt.Errorf("Could not read all string bytes")
	}
	return string(b), nil
}

func ReadVarInt(r io.Reader) (uint64, error) {
	var discriminant uint8
	err := binary.Read(r, binary.LittleEndian, &discriminant)
	if err != nil {
		return 0, err
	}

	var rv uint64
	switch discriminant {
	case 0xff:
		err = binary.Read(r, binary.LittleEndian, &rv)
		if err != nil {
			return 0, err
		}

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x100000000)
		if rv < min {
			return 0, fmt.Errorf("Varint not canonically packed")
		}
	case 0xfe:
		var sv uint32
		binary.Read(r, binary.LittleEndian, &sv)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x10000)
		if rv < min {
			return 0, fmt.Errorf("Varint not canonically packed")
		}
	case 0xfd:
		var sv uint32
		binary.Read(r, binary.LittleEndian, &sv)
		if err != nil {
			return 0, err
		}
		rv = uint64(sv)

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0xfd)
		if rv < min {
			return 0, fmt.Errorf("Varint not canonically packed")
		}
	default:
		rv = uint64(discriminant)
	}

	return rv, nil
}

func ReadIPAddr(r io.Reader) (net.IP, error) {
	b := make([]byte, 16)
	i, err := r.Read(b)
	if i != 16 {
		return nil, fmt.Errorf("Unable to read IP address")
	}
	if err != nil {
		return nil, err
	}
	return net.IP(b), nil
}

func WriteIPAddr(w io.Writer, val net.IP) error {
	b := make([]byte, 16)
	copy(b, val.To16())
	i, err := w.Write(b)
	if i != 16 {
		return fmt.Errorf("Unable to write IP address")
	}
	return err
}

func WriteVarString(w io.Writer, val string) error {
	err := WriteVarInt(w, uint64(len(val)))
	if err != nil {
		return err
	}
	num, err := w.Write([]byte(val))
	if num != len(val) {
		return fmt.Errorf("Not all bytes could be written")
	}
	return nil
}

func WriteVarInt(w io.Writer, val uint64) error {
	if val < 0xfd {
		return binary.Write(w, binary.LittleEndian, uint8(val))
	}

	if val <= math.MaxUint16 {
		err := binary.Write(w, binary.LittleEndian, uint8(0xfd))
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(val))
	}

	if val <= math.MaxUint32 {
		err := binary.Write(w, binary.LittleEndian, uint8(0xfe))
		if err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(val))
	}

	err := binary.Write(w, binary.LittleEndian, uint8(0xff))
	if err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, val)
}

func WriteBigInt256(w io.Writer, i *big.Int) error {
	b := make([]byte, 32)
	numBytes := i.Bytes()
	b = append(b, numBytes...)
	l, err := w.Write(b[len(b)-32:])
	if l != 32 {
		return fmt.Errorf("Couldn't write 32 bytes for big.int")
	}
	return err
}

func ReadBigInt256(r io.Reader) (*big.Int, error) {
	b := make([]byte, 32)
	i, err := r.Read(b)
	if i != 32 {
		return nil, fmt.Errorf("Couldn't read 32 bytes for big.int")
	}
	if err != nil {
		return nil, err
	}
	for b[0] == 0x00 {
		b = b[1:]
	}
	return big.NewInt(0).SetBytes(b), nil
}

func WriteChainHash(w io.Writer, i *chainhash.Hash) error {
	if i == nil {
		i = nullHash
	}
	l, err := w.Write(i.CloneBytes())
	if l != 32 {
		return fmt.Errorf("Couldn't write 32 bytes for chainhash")
	}
	return err
}

func ReadChainHash(r io.Reader) (*chainhash.Hash, error) {
	b := make([]byte, 32)
	i, err := r.Read(b)
	if i != 32 {
		return nil, fmt.Errorf("Couldn't read 32 bytes for chainhash")
	}
	if err != nil {
		return nil, err
	}
	return chainhash.NewHash(b)
}

func init() {
	nullHash, _ = chainhash.NewHashFromStr("0000000000000000000000000000000000000000000000000000000000000000")
}
