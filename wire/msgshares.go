package wire

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var _ P2PoolMessage = &MsgShares{}

type MsgShares struct {
	Shares []Share
}

type Share struct {
	Type           uint64
	MinHeader      SmallBlockHeader
	ShareInfo      ShareInfo
	RefMerkleLink  []*chainhash.Hash
	LastTxOutNonce uint64
	HashLink       HashLink
	MerkleLink     []*chainhash.Hash
}

type HashLink struct {
	State  string
	Length uint64
}

type SmallBlockHeader struct {
	Version       uint64
	PreviousBlock *chainhash.Hash
	Timestamp     int32
	Bits          int32
	Nonce         int32
}

type ShareInfo struct {
	ShareData            ShareData
	SegwitData           SegwitData
	NewTransactionHashes []*chainhash.Hash
	TransactionHashRefs  []TransactionHashRef
	FarShareHash         *chainhash.Hash
	MaxBits              int32
	Bits                 int32
	Timestamp            int32
	AbsHeight            int32
	AbsWork              *big.Int
}

type TransactionHashRef struct {
	ShareCount uint64
	TxCount    uint64
}

type ShareData struct {
	PreviousShareHash *chainhash.Hash
	CoinBase          string
	Nonce             uint32
	PubKeyHash        []byte
	PubKeyHashVersion uint8
	Subsidy           uint64
	Donation          uint16
	StaleInfo         StaleInfo
	DesiredVersion    uint64
}

type StaleInfo uint8

const (
	StaleInfoNone   = StaleInfo(0)
	StaleInfoOrphan = StaleInfo(253)
	StaleInfoDOA    = StaleInfo(254)
)

type SegwitData struct {
	TXIDMerkleLink  []*chainhash.Hash
	WTXIDMerkleRoot *chainhash.Hash
}

func ReadShares(r io.Reader) ([]Share, error) {
	shares := make([]Share, 0)
	count, err := ReadVarInt(r)
	if err != nil {
		return shares, err
	}
	log.Printf("Deserializing %d shares", count)
	for i := uint64(0); i < count; i++ {
		s := Share{}
		s.Type, err = ReadVarInt(r)
		if err != nil {
			return shares, err
		}

		log.Printf("Type is %d", s.Type)

		s.MinHeader, err = ReadSmallBlockHeader(r)
		if err != nil {
			return shares, err
		}

		log.Printf("Minheader is Prevblock: %s", s.MinHeader.PreviousBlock.String())

		s.ShareInfo, err = ReadShareInfo(r, s.Type >= 17)
		if err != nil {
			return shares, err
		}

		log.Printf("Read shareinfo. MaxBits %d, Bits %d, AbsHeight %d, AbsWork: %x", s.ShareInfo.MaxBits, s.ShareInfo.Bits, s.ShareInfo.AbsHeight, s.ShareInfo.AbsWork.Bytes())

		s.RefMerkleLink, err = ReadChainHashList(r)
		if err != nil {
			return shares, err
		}

		err = binary.Read(r, binary.LittleEndian, &s.LastTxOutNonce)
		if err != nil {
			return shares, err
		}

		log.Printf("Read lasttxoutnonce: %d", s.LastTxOutNonce)

		s.HashLink, err = ReadHashLink(r)
		if err != nil {
			return shares, err
		}

		s.MerkleLink, err = ReadChainHashList(r)
		if err != nil {
			return shares, err
		}

		shares = append(shares, s)
	}
	return shares, nil
}

func (m *MsgShares) FromBytes(b []byte) error {
	var err error
	r := bytes.NewReader(b)
	m.Shares, err = ReadShares(r)
	if err != nil {
		return err
	}
	log.Printf("Deserialized %d shares", len(m.Shares))
	return nil
}

func (m *MsgShares) ToBytes() ([]byte, error) {
	var buf bytes.Buffer

	err := WriteVarInt(&buf, uint64(len(m.Shares)))
	if err != nil {
		return nil, err
	}
	for _, s := range m.Shares {
		err = WriteVarInt(&buf, s.Type)
		if err != nil {
			return nil, err
		}
		err = WriteSmallBlockHeader(&buf, s.MinHeader)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (m *MsgShares) Command() string {
	return "shares"
}
