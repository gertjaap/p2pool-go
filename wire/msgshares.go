package wire

import (
	"bytes"
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
	LastTxOutNonce int64
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
	Nonce             int32
	PubKeyHash        []byte
	PubKeyHashVersion int8
	Subsidy           int64
	Donation          int16
	StaleInfo         StaleInfo
	DesiredVersion    uint64
}

type StaleInfo struct {
}

type SegwitData struct {
	TXIDMerkleLink  []*chainhash.Hash
	WTXIDMerkleRoot *chainhash.Hash
}

func (m *MsgShares) FromBytes(b []byte) error {
	r := bytes.NewReader(b)
	m.Shares = make([]Share, 0)
	count, err := ReadVarInt(r)
	if err != nil {
		return err
	}
	log.Printf("Deserializing %d shares", count)
	for i := uint64(0); i < count; i++ {
		s := Share{}
		s.Type, err = ReadVarInt(r)
		if err != nil {
			return err
		}

		m.Shares = append(m.Shares, s)
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

	}
	return buf.Bytes(), nil
}

func (m *MsgShares) Command() string {
	return "shares"
}
