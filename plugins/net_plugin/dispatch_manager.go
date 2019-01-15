package net_plugin

import (
	"encoding/binary"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
	"unsafe"
)

type BlockRequest struct {
	id         common.BlockIdType
	localRetry bool
}
type blockOrigin struct {
	id   common.BlockIdType
	conn *Connection
}
type transactionOrigin struct {
	id   common.TransactionIdType
	conn *Connection
}

type dispatchManager struct {
	justSendItMax uint32
	regBlks       []BlockRequest
	reqTrx        []common.TransactionIdType
	//receivedBlocks       map[common.BlockIdType]*Peer
	//receivedTransactions map[common.TransactionIdType]*Peer

	receivedBlocks       map[common.BlockIdType][]*Connection
	receivedTransactions map[common.TransactionIdType][]*Connection
	myImpl               *netPluginIMpl
}

func NewDispatchManager(impl *netPluginIMpl) *dispatchManager {
	return &dispatchManager{
		myImpl: impl,
	}
}

func (d *dispatchManager) bcastBlock(myImpl *netPluginIMpl, bsum *types.SignedBlock) {
	skips := map[*Connection]int{}

	bid := bsum.BlockID()
	bnum := bsum.BlockNumber()
	_, ok := d.receivedBlocks[bid]
	if ok {
		for i, p := range d.receivedBlocks[bid] {
			skips[p] = i
		}
	}
	delete(d.receivedBlocks, bid)

	msg := SignedBlockMessage{*bsum}
	packsize, _ := rlp.EncodeSize(msg)
	msgsiz := uint32(packsize) + uint32(unsafe.Sizeof(packsize))
	pendingNotify := NoticeMessage{}
	pendingNotify.KnownBlocks.Mode = normal
	pendingNotify.KnownBlocks.IDs = append(pendingNotify.KnownBlocks.IDs, &bid)
	pendingNotify.KnownTrx.Mode = none

	pbstate := PeerBlockState{
		ID:            bid,
		BlockNum:      bnum,
		IsKnown:       false,
		IsNoticed:     true,
		RequestedTime: common.TimePoint(0),
	}
	// skip will be empty if our producer emitted this block so just send it
	if (largeMsgNotify && msgsiz > d.justSendItMax) && len(skips) > 0 {
		fcLog.Info("block_size is %d ,sending notify", msgsiz)
		myImpl.sendAll(&pendingNotify, func(p *Connection) bool {
			_, ok := skips[p]
			if ok || !p.current() {
				return false
			}
			unknown := p.addPeerBlock(&pbstate)
			if !unknown {
				netLog.Error("%s already has knowledge of block %d", p.peerAddr, pbstate.BlockNum)
			}
			return unknown
		})
	} else {
		pbstate.IsKnown = true
		//for _,_ =range my.peers{
		p := Connection{}
		_, ok := skips[&p]
		if ok || !p.current() {
			//continue
		}
		p.addPeerBlock(&pbstate)
		p.write(&msg)
		//}
	}
}

func (d *dispatchManager) rejectedBlock(id *common.BlockIdType) {
	fcLog.Debug("not sending rejected block %s", id)
	_, ok := d.receivedBlocks[*id]
	if ok {
		delete(d.receivedBlocks, *id)
	}
}

func (d *dispatchManager) recvBlock(p *Connection, id *common.BlockIdType, bnum uint32) {
	d.receivedBlocks[*id] = append(d.receivedBlocks[*id], p)
	IdsCount := len(p.lastReq.ReqBlocks.IDs)
	if p != nil && p.lastReq != nil && p.lastReq.ReqBlocks.Mode != none && IdsCount > 0 && p.lastReq.ReqBlocks.IDs[IdsCount-1] == id {
		p.lastReq = &RequestMessage{} //TODO
	}

	pbs := PeerBlockState{
		ID:            *id,
		BlockNum:      bnum,
		IsKnown:       false,
		IsNoticed:     true,
		RequestedTime: common.TimePoint(0),
	}
	p.addPeerBlock(&pbs)
	fcLog.Debug("canceling wait on %s", p.peerAddr)
	p.cancelWait()
}

func (d *dispatchManager) bcastTransaction(trx *types.PackedTransaction) { // TODO impl
	skips := map[*Connection]int{}
	id := trx.ID()

	peers, ok := d.receivedTransactions[id]
	if ok {
		for i := 0; i < len(peers); i++ {
			skips[peers[i]] = i
		}
	}
	delete(d.receivedTransactions, id)

	for i, _ := range d.reqTrx {
		if d.reqTrx[i] == id {
			d.reqTrx = append(d.reqTrx[:i], d.reqTrx[i+1:]...) //TODO req_trx.erase(ref)
			break
		}
	}

	if !d.myImpl.localTxns.GetById().Find(id).IsEnd() { //found
		fcLog.Info("found trx id in local_trxs")
		return
	}

	trxExpiration := trx.Expiration()
	msg := PackedTransactionMessage{*trx}
	packedTrxBuf, _ := rlp.EncodeToBytes(msg)

	packSize := uint32(len(packedTrxBuf))
	bufsize := packSize + uint32(unsafe.Sizeof(packSize))
	buffer := make([]byte, bufsize)

	binary.LittleEndian.PutUint32(buffer, uint32(unsafe.Sizeof(packSize))) // binary.LittleEndian.PutUint32(buffer,4)
	buffer = append(buffer, packedTrxBuf...)
	nts := NodeTransactionState{
		ID:            id,
		Expires:       trxExpiration,
		PackedTxn:     *trx,
		SerializedTxn: buffer,
		BlockNum:      0,
		TrueBlock:     0,
		Requests:      0,
	}
	d.myImpl.localTxns.Insert(nts)

	if !largeMsgNotify || bufsize <= d.justSendItMax {
		packedTrx := PackedTransactionMessage{*trx}
		d.myImpl.sendAll(&packedTrx, func(c *Connection) bool {
			_, ok := skips[c]
			if ok || c.syncing {
				return false
			}

			bs := c.trxState.GetById().Find(id)
			unknown := bs.IsEnd()

			if unknown {
				c.trxState.Insert(TransactionState{id, false, true, 0, trxExpiration, common.TimePoint(0)})
				fcLog.Debug("sending notice to  %s", c.peerAddr)
			} else {
				c.trxState.Modify(bs, func(state *TransactionState) {
					(*state).Expires = trxExpiration
				})
			}
			return unknown
		})
	} else {
		pendingNotify := NoticeMessage{}
		pendingNotify.KnownTrx.Mode = normal
		pendingNotify.KnownTrx.IDs = append(pendingNotify.KnownTrx.IDs, &id)
		pendingNotify.KnownBlocks.Mode = none
		d.myImpl.sendAll(&pendingNotify, func(c *Connection) bool {
			_, ok := skips[c]
			if ok || c.syncing {
				return false
			}

			bs := c.trxState.GetById().Find(id)
			unknown := bs.IsEnd()
			if unknown {
				fcLog.Debug("sending notice to  %s", c.peerAddr)
				c.trxState.Insert(TransactionState{id, false, true, 0, trxExpiration, common.TimePoint(0)})
			} else {
				c.trxState.Modify(bs, func(state *TransactionState) {
					(*state).Expires = trxExpiration
				})
			}
			return unknown
		})
	}
}

func (d *dispatchManager) rejectedTransaction(id *common.TransactionIdType) {
	fcLog.Debug("not sending rejected transaction %s", id)
	_, ok := d.receivedTransactions[*id]
	if ok {
		delete(d.receivedTransactions, *id)
	}
}

func (d *dispatchManager) recvTransaction(p *Connection, id *common.TransactionIdType) {
	d.receivedTransactions[*id] = append(d.receivedTransactions[*id], p)
	idsCount := len(p.lastReq.ReqTrx.IDs)
	if p != nil && p.lastReq != nil && p.lastReq.ReqTrx.Mode != none && idsCount > 0 && p.lastReq.ReqTrx.IDs[idsCount-1] == id { //TODO c && c->last_req
		//p.lastReq.reset()
		p.lastReq = &RequestMessage{}
	}
	fcLog.Debug("canceling wait on %s", p.peerAddr)
	p.cancelWait()

}

func (d *dispatchManager) recvNotice(c *Connection, msg *NoticeMessage, generated bool) {
	req := RequestMessage{}
	req.ReqTrx.Mode = none
	req.ReqBlocks.Mode = none
	sendReq := false
	cc := d.myImpl.ChainPlugin.Chain()

	if msg.KnownTrx.Mode == normal {
		req.ReqTrx.Mode = normal
		req.ReqTrx.Pending = 0
		for _, t := range msg.KnownTrx.IDs {
			tx := d.myImpl.localTxns.GetById().Find(*t)

			if tx.IsEnd() {
				fcLog.Debug("did not find %s", t.String())
				//At this point the details of the txn are not known, just its id. This
				//effectively gives 120 seconds to learn of the details of the txn which
				//will update the expiry in bcast_transaction
				c.trxState.Insert(TransactionState{*t, true, true, 0,
					common.TimePointSec(common.Now() + 120), common.TimePoint(0)})

				req.ReqTrx.IDs = append(req.ReqTrx.IDs, t)
				d.reqTrx = append(d.reqTrx, *t)
			} else {
				fcLog.Debug("big msg manager found txn id in table,%s", t.String())
			}
		}
		sendReq = !(len(req.ReqTrx.IDs) == 0)
		fcLog.Debug("big msg manager send_req ids list has %d entries", len(req.ReqTrx.IDs))

	} else if msg.KnownTrx.Mode != none {
		log.Error("passed a notice_message with something other than a normal on none known_trx")
		return
	}

	if msg.KnownBlocks.Mode == normal {
		req.ReqBlocks.Mode = normal

		for _, blkID := range msg.KnownBlocks.IDs {
			b := &types.SignedBlock{}
			entry := PeerBlockState{*blkID, 0, true, true, common.TimePoint(0)}
			Try(func() {
				b = cc.FetchBlockById(*blkID)
				if b != nil {
					entry.BlockNum = b.BlockNumber()
				}
			}).Catch(func(ex *exception.AssertException) {
				netLog.Info("caught assert on fetch_block_by_id, %s", ex.What())
				//keep going, client can ask another peer
			}).Catch(func(e interface{}) {
				netLog.Error("failed to retrieve block for id")
			}).End()

			if common.Empty(b) {
				sendReq = true
				req.ReqBlocks.IDs = append(req.ReqBlocks.IDs, blkID)
				entry.RequestedTime = common.Now()
			}
			c.addPeerBlock(&entry)
		}
	} else if msg.KnownBlocks.Mode != none {
		netLog.Error("passed a notice_message with something other than a normal on none known_blocks")
		return
	}

	netLog.Debug("send req = %s", sendReq)
	if sendReq {
		c.write(&req)
		c.fetchWait()
		c.lastReq = &req
	}
}

func (d *dispatchManager) retryFetch(c *Connection) {
	if common.Empty(c.lastReq) {
		return
	}
	fcLog.Debug("failed to fetch from %s", c.peerAddr)
	var (
		tid   common.TransactionIdType
		bid   common.BlockIdType
		isTxn bool
	)

	reqTrxCount := len(c.lastReq.ReqTrx.IDs)
	reqBlockCount := len(c.lastReq.ReqBlocks.IDs)
	if c.lastReq.ReqTrx.Mode == normal && reqTrxCount > 0 {
		isTxn = true
		tid = *c.lastReq.ReqTrx.IDs[reqTrxCount-1]
	} else if c.lastReq.ReqBlocks.Mode == normal && reqBlockCount > 0 {
		bid = *c.lastReq.ReqBlocks.IDs[reqBlockCount-1]
	} else {
		fcLog.Debug("no retry,block mode = %s trx mode = %s\n", modeTostring[c.lastReq.ReqBlocks.Mode], modeTostring[c.lastReq.ReqTrx.Mode])
		return
	}

	for _, conn := range d.myImpl.connections {
		if conn == c || conn.lastReq != nil {
			continue
		}

		sendIt := false
		if isTxn {
			trx := conn.trxState.GetById().Find(tid)
			sendIt = !trx.IsEnd() && trx.Value().IsKnownByPeer
		} else {
			blk := conn.blkState.GetById().Find(bid)
			sendIt = !blk.IsEnd() && blk.Value().IsKnown
		}
		if sendIt {
			conn.write(c.lastReq)
			conn.fetchWait()
			conn.lastReq = c.lastReq
			return
		}

	}

	// at this point no other peer has it, re-request or do nothing?
	if c.connected() {
		c.write(c.lastReq)
		c.fetchWait()
	}
}
