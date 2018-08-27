package p2p

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)


type P2PMessage interface {
	GetType() P2PMessageType
}

type HandshakeMessage struct {
	// net_plugin/protocol.hpp handshake_message
	NetworkVersion           uint16             `json:"network_version"`
	ChainID                  common.ChainIDType `json:"chain_id"`
	NodeID                   common.NodeIDType  `json:"node_id"` // sha256
	Key                      ecc.PublicKey      `json:"key"`     // can be empty, producer key, or peer key
	Time                     common.Tstamp      `json:"time"`    // time?!
	Token                    common.Sha256      `json:"token"`   // digest of time to prove we own the private `key`
	Signature                ecc.Signature      `json:"sig"`     // can be empty if no key, signature of the digest above
	P2PAddress               string             `json:"p2p_address"`
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIDType `json:"last_irreversible_block_id"`
	HeadNum                  uint32             `json:"head_num"`
	HeadID                   common.BlockIDType `json:"head_id"`
	OS                       string             `json:"os"`
	Agent                    string             `json:"agent"`
	Generation               uint16             `json:"generation"`
}

func (m *HandshakeMessage) GetType() P2PMessageType {
	return HandshakeMessageType
}

func (m *HandshakeMessage) String() string {
	// return fmt.Sprintf("Handshake: Head [%d] Last Irreversible [%d] Time [%s]", m.HeadNum, m.LastIrreversibleBlockNum, m.Time)
	return "handshakemessage"
}

type ChainSizeMessage struct {
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIDType `json:"last_irreversible_block_id"`
	HeadNum                  uint32             `json:"head_num"`
	HeadID                   common.BlockIDType `json:"head_id"`
}

func (m *ChainSizeMessage) GetType() P2PMessageType {
	return ChainSizeType
}

type GoAwayReason uint32

const (
	GoAwayNoReason = uint8(iota)
	GoAwaySelfConnect
	GoAwayDuplicate
	GoAwayWrongChain
	GoAwayWrongVersion
	GoAwayForked
	GoAwayUnlinkable
	GoAwayBadTransaction
	GoAwayValidation
	GoAwayAuthentication
	GoAwayFatalOther
	GoAwayBenignOther
	GoAwayCrazy
)

type GoAwayMessage struct {
	Reason GoAwayReason  `json:"reason"`
	NodeID common.Sha256 `json:"node_id"`
}

func (m *GoAwayMessage) GetType() P2PMessageType {
	return GoAwayMessageType
}

type TimeMessage struct {
	Origin      common.Tstamp `json:"org"`
	Receive     common.Tstamp `json:"rec"`
	Transmit    common.Tstamp `json:"xmt"`
	Destination common.Tstamp `json:"dst"`
}

func (m *TimeMessage) GetType() P2PMessageType {
	return TimeMessageType
}

func (t *TimeMessage) String() string {
	return fmt.Sprintf("Origin [%s], Receive [%s], Transmit [%s], Destination [%s]", t.Origin, t.Receive, t.Transmit, t.Destination)
}

type IDListMode uint32

const (
	none IDListMode = iota
	catch_up
	last_irr_catch_up
	normal
)

type OrderedTransactionIDs struct {
	// Unknown [3]byte             `json:"-"` ///// WWUUuuuuuuuuuuuutzthat ?
	Mode    IDListMode                  `json:"mode"`
	Pending uint32                      `json:"pending"`
	IDs     []*common.TransactionIDType `json:"ids"`
}
type OrderedBlockIDs struct {
	// Unknown [3]byte             `json:"-"` ///// wuuttzthat?
	Mode    IDListMode            `json:"mode"`
	Pending uint32                `json:"pending"`
	IDs     []*common.BlockIDType `json:"ids"`
}

type NoticeMessage struct {
	KnownTrx    OrderedTransactionIDs `json:"known_trx"`
	KnownBlocks OrderedBlockIDs       `json:"known_blocks"`
}

func (m *NoticeMessage) GetType() P2PMessageType {
	return NoticeMessageType
}

type SyncRequestMessage struct {
	StartBlock uint32 `json:"start_block"`
	EndBlock   uint32 `json:"end_block"`
}

func (m *SyncRequestMessage) GetType() P2PMessageType {
	return SyncRequestMessageType
}
func (m *SyncRequestMessage) String() string {
	return fmt.Sprintf("SyncRequest: Start Block [%d] End Block [%d]", m.StartBlock, m.EndBlock)
}

type RequestMessage struct {
	ReqTrx    OrderedTransactionIDs `json:"req_trx"`
	ReqBlocks OrderedBlockIDs       `json:"req_blocks"`
}

func (m *RequestMessage) GetType() P2PMessageType {
	return RequestMessageType
}

type SignedBlockMessage struct {
	types.SignedBlock
}

func (s *SignedBlockMessage) GetType() P2PMessageType {
	return SignedBlockType
}

type PackedTransactionMessage struct {
	types.PackedTransaction
}

func (m *PackedTransactionMessage) GetType() P2PMessageType {
	return PackedTransactionMessageType
}
