package net_plugin

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type NetMessageType byte

const (
	HandshakeMessageType NetMessageType = iota // 0
	ChainSizeType
	GoAwayMessageType
	TimeMessageType
	NoticeMessageType // 4
	RequestMessageType
	SyncRequestMessageType
	SignedBlockType
	PackedTransactionMessageType //8
)

type MessageReflectTypes struct {
	Name        string
	ReflectType reflect.Type
}

var messageAttributes = []MessageReflectTypes{
	{Name: "Handshake", ReflectType: reflect.TypeOf(HandshakeMessage{})},
	{Name: "ChainSize", ReflectType: reflect.TypeOf(ChainSizeMessage{})},
	{Name: "GoAway", ReflectType: reflect.TypeOf(GoAwayMessage{})},
	{Name: "Time", ReflectType: reflect.TypeOf(TimeMessage{})},
	{Name: "Notice", ReflectType: reflect.TypeOf(NoticeMessage{})},
	{Name: "Request", ReflectType: reflect.TypeOf(RequestMessage{})},
	{Name: "SyncRequest", ReflectType: reflect.TypeOf(SyncRequestMessage{})},
	{Name: "SignedBlock", ReflectType: reflect.TypeOf(SignedBlockMessage{})},
	{Name: "PackedTransaction", ReflectType: reflect.TypeOf(PackedTransactionMessage{})},
}

func (t NetMessageType) isValid() bool {
	index := byte(t)
	return int(index) < len(messageAttributes)
}

func (t NetMessageType) Name() (string, bool) {
	index := byte(t)
	if !t.isValid() {
		return "Unknown", false
	}

	attr := messageAttributes[index]
	return attr.Name, true
}

func (t NetMessageType) reflectTypes() (MessageReflectTypes, bool) {
	index := byte(t)
	if !t.isValid() {
		return MessageReflectTypes{}, false
	}

	attr := messageAttributes[index]
	return attr, true
}

type NetMessage interface {
	GetType() NetMessageType
	String() string
}

type HandshakeMessage struct {
	NetworkVersion           uint16             `json:"network_version"` // incremental value above a computed base
	ChainID                  common.ChainIdType `json:"chain_id"`        // used to identify chain
	NodeID                   common.NodeIdType  `json:"node_id"`         // used to identify peers and prevent self-connect
	Key                      ecc.PublicKey      `json:"key"`             // authentication key; may be a producer or peer key, or empty
	Time                     common.TimePoint   `json:"time"`
	Token                    crypto.Sha256      `json:"token"` // digest of time to prove we own the private `key`
	Signature                ecc.Signature      `json:"sig"`   // signature for the digest
	P2PAddress               string             `json:"p2p_address"`
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIdType `json:"last_irreversible_block_id"`
	HeadNum                  uint32             `json:"head_num"`
	HeadID                   common.BlockIdType `json:"head_id"`
	OS                       string             `json:"os"`
	Agent                    string             `json:"agent"`
	Generation               uint16             `json:"generation"`
}

func (h *HandshakeMessage) GetType() NetMessageType {
	return HandshakeMessageType
}
func (h *HandshakeMessage) String() string {
	bytes, _ := json.Marshal(h)
	return string(bytes)
}

type ChainSizeMessage struct {
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIdType `json:"last_irreversible_block_id"`
	HeadNum                  uint32             `json:"head_num"`
	HeadID                   common.BlockIdType `json:"head_id"`
}

func (c *ChainSizeMessage) GetType() NetMessageType {
	return ChainSizeType
}
func (c *ChainSizeMessage) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

type GoAwayReason uint32

const (
	noReason       = GoAwayReason(iota) //no reason to go away
	selfConnect                         //the connection is to itself
	duplicate                           //the connection is redundant
	wrongChain                          //the peer's chain id doesn't match
	wrongVersion                        //the peer's network version doesn't match
	forked                              //the peer's irreversible blocks are different
	unlinkable                          //the peer sent a block we couldn't use
	badTransaction                      //the peer sent a transaction that failed verification
	validation                          //the peer sent a block that failed validation
	benignOther                         //reasons such as a timeout. not fatal but warrant resetting
	fatalOther                          //a catch-all for errors we don't have discriminated
	authentication                      //peer failed authenicatio
	crazy                               //some crazy reason
)

var ReasonStr = map[GoAwayReason]string{
	noReason:       "no reason",
	selfConnect:    "self connect",
	duplicate:      "duplicate",
	wrongChain:     "wrong chain",
	wrongVersion:   "wrong version",
	forked:         "chain is forked",
	unlinkable:     "unlinkable block received",
	badTransaction: "bad transaction",
	validation:     "invalid block",
	benignOther:    "some other non-fatal condition",
	fatalOther:     "some other failure",
	authentication: "authentication failure",
	crazy:          "some crazy reason",
}

type GoAwayMessage struct {
	Reason GoAwayReason      `json:"reason"`
	NodeID common.NodeIdType `json:"node_id"` //for duplicate notification
}

func (g *GoAwayMessage) GetType() NetMessageType {
	return GoAwayMessageType
}
func (g *GoAwayMessage) String() string {
	bytes, _ := json.Marshal(g)
	return string(bytes)
}

type TimeMessage struct {
	Org common.TimePoint `json:"org"` //origin timestamp
	Rec common.TimePoint `json:"rec"` //receive timestamp
	Xmt common.TimePoint `json:"xmt"` //transmit timestamp
	Dst common.TimePoint `json:"dst"` //destination timestamp
}

func (t *TimeMessage) GetType() NetMessageType {
	return TimeMessageType
}
func (t *TimeMessage) String() string {
	return fmt.Sprintf("Origin [%s], Receive [%s], Transmit [%s], Destination [%s]", t.Org, t.Rec, t.Xmt, t.Dst)
}

type IdListMode uint32

const (
	none IdListMode = iota
	catchUp
	lastIrrCatchUp
	normal
)

var modeStr = map[IdListMode]string{
	none:           "none",
	catchUp:        "catch up",
	lastIrrCatchUp: "last irreversible",
	normal:         "normal",
}

type OrderedTransactionIDs struct {
	Mode    IdListMode                 `json:"mode"`
	Pending uint32                     `json:"pending"`
	IDs     []common.TransactionIdType `json:"ids"`
}

func (t *OrderedTransactionIDs) empty() bool {
	return t.Mode == none || len(t.IDs) == 0
}

type OrderedBlockIDs struct {
	Mode    IdListMode           `json:"mode"`
	Pending uint32               `json:"pending"`
	IDs     []common.BlockIdType `json:"ids"`
}

func (b *OrderedBlockIDs) empty() bool {
	return b.Mode == none || len(b.IDs) == 0
}

type NoticeMessage struct {
	KnownTrx    OrderedTransactionIDs `json:"known_trx"`
	KnownBlocks OrderedBlockIDs       `json:"known_blocks"`
}

func (n *NoticeMessage) GetType() NetMessageType {
	return NoticeMessageType
}
func (n *NoticeMessage) String() string {
	bytes, _ := json.Marshal(n)
	return string(bytes)
}

type SyncRequestMessage struct {
	StartBlock uint32 `json:"start_block"`
	EndBlock   uint32 `json:"end_block"`
}

func (s *SyncRequestMessage) GetType() NetMessageType {
	return SyncRequestMessageType
}
func (s *SyncRequestMessage) String() string {
	return fmt.Sprintf("SyncRequest: Start Block [%d] End Block [%d]", s.StartBlock, s.EndBlock)
}

type RequestMessage struct {
	ReqTrx    OrderedTransactionIDs `json:"req_trx"`
	ReqBlocks OrderedBlockIDs       `json:"req_blocks"`
}

func (r *RequestMessage) GetType() NetMessageType {
	return RequestMessageType
}
func (r *RequestMessage) String() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}
func (r RequestMessage) IsEmpty() bool {
	return r.ReqTrx.empty() && r.ReqBlocks.empty()
}

type SignedBlockMessage struct {
	types.SignedBlock
}

func (s *SignedBlockMessage) GetType() NetMessageType {
	return SignedBlockType
}
func (s *SignedBlockMessage) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

type PackedTransactionMessage struct {
	types.PackedTransaction
}

func (p *PackedTransactionMessage) GetType() NetMessageType {
	return PackedTransactionMessageType
}
func (p *PackedTransactionMessage) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

/**
Goals of Network Code
1. low latency to minimize missed blocks and potentially reduce block interval
2. minimize redundant data between blocks and transactions.
3. enable rapid sync of a new node
4. update to new boost / fc

State:
   All nodes know which blocks and transactions they have
   All nodes know which blocks and transactions their peers have
   A node knows which blocks and transactions it has requested
   All nodes know when they learned of a transaction

   send hello message
   write loop (true)
      if peer knows the last irreversible block {
         if peer does not know you know a block or transactions
            send the ids you know (so they don't send it to you)
            yield continue
         if peer does not know about a block
            send transactions in block peer doesn't know then send block summary
            yield continue
         if peer does not know about new public endpoints that you have verified
            relay new endpoints to peer
            yield continue
         if peer does not know about transactions
            sends the oldest transactions that is not known by the remote peer
            yield continue
         wait for new validated block, transaction, or peer signal from network fiber
      } else {
         we assume peer is in sync mode in which case it is operating on a
         request / response basis

         wait for notice of sync from the read loop
      }


    read loop
      if hello message
         verify that peers Last Ir Block is in our state or disconnect, they are on fork
         verify peer network protocol

      if notice message update list of transactions known by remote peer
      if trx message then insert into global state as unvalidated
      if blk summary message then insert into global state *if* we know of all dependent transactions
         else close connection


    if my head block < the LIB of a peer and my head block age > block interval * round_size/2 then
    enter sync mode...
        divide the block numbers you need to fetch among peers and send fetch request
        if peer does not respond to request in a timely manner then make request to another peer
        ensure that there is a constant queue of requests in flight and everytime a request is filled
        send of another request.

     Once you have caught up to all peers, notify all peers of your head block so they know that you
     know the LIB and will start sending you real time transactions

parallel fetches, request in groups


only relay transactions to peers if we don't already know about it.

send a notification rather than a transaction if the txn is > 3mtu size.

*/
