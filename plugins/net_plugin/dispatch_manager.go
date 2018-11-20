package net_plugin

import (
	"encoding/binary"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"unsafe"
)

type BlockRequest struct {
	id         common.BlockIdType
	localRetry bool
}
type blockOrigin struct {
	id   common.BlockIdType
	peer *Peer
}
type transactionOrigin struct {
	id   common.TransactionIdType
	Peer *Peer
}

type dispatchManager struct {
	justSendItMax uint32
	regBlks       []BlockRequest
	reqTrx        []common.TransactionIdType
	//receivedBlocks       map[common.BlockIdType]*Peer
	//receivedTransactions map[common.TransactionIdType]*Peer

	receivedBlocks       map[common.BlockIdType][]*Peer
	receivedTransactions map[common.TransactionIdType][]*Peer
}

func NewDispatchManager() *dispatchManager {
	return &dispatchManager{}
}

func (d *dispatchManager) bcastBlock(myImpl *netPluginIMpl, bsum *types.SignedBlock) {
	skips := map[*Peer]int{}

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

	pbstate := peerBlockState{
		id:          bid,
		blockNum:    bnum,
		isKnown:     false,
		isNoticed:   true,
		requestTime: common.TimePoint(0),
	}
	// skip will be empty if our producer emitted this block so just send it
	if (largeMsgNotify && msgsiz > d.justSendItMax) && len(skips) > 0 {
		//fc_ilog(logger, "block size is ${ms}, sending notify",("ms", msgsiz))
		netlog.Info("block_size is %d ,sending notify", msgsiz)
		myImpl.sendAll(&pendingNotify, func(p *Peer) bool {
			_, ok := skips[p]
			if ok || !p.current() {
				return false
			}
			unknown := p.addPeerBlock(&pbstate)
			if !unknown {
				//elog("${p} already has knowledge of block ${b}", ("p",c->peer_name())("b",pbstate.block_num))
				netlog.Error("%s already has knowledge of block %d", p.peerAddr, pbstate.blockNum)
			}
			return unknown
		})
	} else {
		pbstate.isKnown = true
		//for _,_ =range my.peers{
		p := Peer{}
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
	//fc_dlog(logger,"not sending rejected transaction ${tid}",("tid",id));
	netlog.Debug("not sending rejected block %s", id)
	_, ok := d.receivedBlocks[*id]
	if ok {
		delete(d.receivedBlocks, *id)
	}
}

func (d *dispatchManager) recvBlock(p *Peer, id *common.BlockIdType, bnum uint32) {
	d.receivedBlocks[*id] = append(d.receivedBlocks[*id], p)
	IdsCount := len(p.lastReq.ReqBlocks.IDs)
	if p != nil && p.lastReq != nil && p.lastReq.ReqBlocks.Mode != none && IdsCount > 0 && p.lastReq.ReqBlocks.IDs[IdsCount-1] == id {
		p.lastReq = &RequestMessage{} //TODO
	}

	pbs := peerBlockState{
		id:          *id,
		blockNum:    bnum,
		isKnown:     false,
		isNoticed:   true,
		requestTime: common.TimePoint(0),
	}
	p.addPeerBlock(&pbs)
	//fc_dlog(logger, "canceling wait on ${p}", ("p",c->peer_name()));

	netlog.Debug("canceling wait on %s", p.peerAddr)
	p.cancelWait()
}

func (d *dispatchManager) bcastTransaction(myImpl *netPluginIMpl, trx *types.PackedTransaction) { // TODO impl
	skips := map[*Peer]int{}
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

	if myImpl.localTxns.getIndex("by_id").findLocalTrxById(id) != nil { //found
		netlog.Info("found trxid in local_trxs")
		return
	}

	trxExpiration := trx.Expiration()
	msg := PackedTransactionMessage{*trx}
	packedTrxBuf, _ := rlp.EncodeToBytes(msg)

	packSize := uint32(len(packedTrxBuf))
	bufsize := packSize + uint32(unsafe.Sizeof(packSize))
	buffer := make([]byte, bufsize)

	binary.LittleEndian.PutUint32(buffer, uint32(unsafe.Sizeof(packSize))) //TODO binary.LittleEndian.PutUint32(buffer,4)
	buffer = append(buffer, packedTrxBuf...)
	nts := nodeTransactionState{
		id:            id,
		expires:       trxExpiration,
		packedTxn:     *trx,
		serializedTxn: buffer,
		blockNum:      0,
		trueBlock:     0,
		requests:      0,
	}
	myImpl.localTxns.insertNodeTrx(&nts)

	if !largeMsgNotify || bufsize <= d.justSendItMax {
		packedTrx := PackedTransactionMessage{*trx}
		myImpl.sendAll(&packedTrx, func(p *Peer) bool {
			_, ok := skips[p]
			if ok || p.syncing {
				return false
			}

			//bs := c->trx_state.find(id);
			//	bool unknown = bs == c->trx_state.end()
			bs := p.trxState.getIndex("by_id").findTrxById(id)
			unknown := bs == nil

			if unknown {
				//fc_dlog(logger, "sending notice to ${n}", ("n",c->peer_name() ) );
				//c->trx_state.insert(transaction_state({id,false,true,0,trx_expiration,time_point() }))
				netlog.Debug("sending notice to  %s", p.peerAddr)
				p.trxState.insertTrx(&transactionState{
					id, false, true, 0, trxExpiration, common.TimePoint(0),
				})
			} else {
				//c->trx_state.modify(bs, ute)
				p.trxState.modify(bs, true, func(in common.ElementObject) {
					re := in.(*transactionState)
					re.expires = trxExpiration
				})
			}
			return unknown
		})
	} else {
		pendingNotify := NoticeMessage{}
		pendingNotify.KnownTrx.Mode = normal
		pendingNotify.KnownTrx.IDs = append(pendingNotify.KnownTrx.IDs, &id)
		pendingNotify.KnownBlocks.Mode = none
		myImpl.sendAll(&pendingNotify, func(p *Peer) bool {
			_, ok := skips[p]
			if ok || p.syncing {
				return false
			}
			//bs := c->trx_state.find(id);
			//	bool unknown = bs == c->trx_state.end()
			bs := p.trxState.getIndex("by_id").findTrxById(id)
			unknown := bs == nil
			if unknown {
				//fc_dlog(logger, "sending notice to ${n}", ("n",c->peer_name() ) );
				//c->trx_state.insert(transaction_state({id,false,true,0,trx_expiration,time_point() }))
				netlog.Debug("sending notice to  %s", p.peerAddr)
				p.trxState.insertTrx(&transactionState{
					id, false, true, 0, trxExpiration, common.TimePoint(0),
				})
			} else {
				//ute :=updateTxnExpiry{trxExpiration}
				//c->trx_state.modify(bs, ute)
				p.trxState.modify(bs, true, func(in common.ElementObject) {
					re := in.(*transactionState)
					re.expires = trxExpiration
				})
			}
			return unknown
		})
	}

}

func (d *dispatchManager) rejectedTransaction(id *common.TransactionIdType) {
	//fc_dlog(logger,"not sending rejected transaction ${tid}",("tid",id));
	netlog.Debug("not sending rejected transaction %s \n", id)
	_, ok := d.receivedTransactions[*id]
	if ok {
		delete(d.receivedTransactions, *id)
	}
}

func (d *dispatchManager) recvTransaction(p *Peer, id *common.TransactionIdType) {
	d.receivedTransactions[*id] = append(d.receivedTransactions[*id], p)
	idsCount := len(p.lastReq.ReqTrx.IDs)
	if p != nil && p.lastReq != nil && p.lastReq.ReqTrx.Mode != none && idsCount > 0 && p.lastReq.ReqTrx.IDs[idsCount-1] == id { //TODO c && c->last_req
		//p.lastReq.reset()
		p.lastReq = &RequestMessage{}
	}
	//fc_dlog(logger, "canceling wait on ${p}", ("p",c->peer_name()));
	netlog.Debug("canceling wait on %s \n", p.peerAddr)
	p.cancelWait()

}

func (d *dispatchManager) recvNotice(myImpl *netPluginIMpl, p *Peer, msg *NoticeMessage, generated bool) {
	req := RequestMessage{}
	req.ReqTrx.Mode = none
	req.ReqBlocks.Mode = none
	sendReq := false
	//controller &cc = my_impl->chain_plug->chain();

	if msg.KnownTrx.Mode == normal {
		req.ReqTrx.Mode = normal
		req.ReqTrx.Pending = 0
		for _, t := range msg.KnownTrx.IDs {
			tx := myImpl.localTxns.getIndex("by_id").findLocalTrxById(*t)
			if tx == nil {
				//fc_dlog(logger,"did not find ${id}",("id",t));
				netlog.Debug("did not find %s", t.String())

				//At this point the details of the txn are not known, just its id. This
				//effectively gives 120 seconds to learn of the details of the txn which
				//will update the expiry in bcast_transaction

				trxState := transactionState{
					id:              *t,
					isKnownByPeer:   true,
					isNoticedToPeer: true,
					blockNum:        0,
					expires:         common.TimePointSec(common.Now()) + 120,
					requestedTime:   common.TimePoint(0),
				}
				p.trxState.insertTrx(&trxState)
				req.ReqTrx.IDs = append(req.ReqTrx.IDs, t)
				d.reqTrx = append(d.reqTrx, *t)
			} else {
				//fc_dlog(logger,"big msg manager found txn id in table, ${id}",("id", t));
				netlog.Debug("big msg manager found txn id in table,%s", t.String())
			}
		}
		sendReq = !(len(req.ReqTrx.IDs) == 0)
		//fc_dlog(logger,"big msg manager send_req ids list has ${ids} entries", ("ids", req.req_trx.ids.size()));
		netlog.Debug("big msg manager send_req ids list has %d entries\n", len(req.ReqTrx.IDs))

	} else if msg.KnownTrx.Mode == none {
		netlog.Error("passed a notice_message with something other than a normal on none known_trx")
		return
	}

	if msg.KnownBlocks.Mode == normal {
		req.ReqBlocks.Mode = normal
		b := types.SignedBlock{}
		for _, blkID := range msg.KnownBlocks.IDs {
			entry := peerBlockState{*blkID, 0, true, true, common.TimePoint(0)}
			Try(func() {
				//b = cc.fetchBlockByID(blkID)
				if &b != nil { //TODO
					entry.blockNum = b.BlockNumber()
				}
			}).Catch(func(ex exception.AssertException) {
				netlog.Info("caught assert on fetch_block_by_id, %s", ex.What())
				//keep going, client can ask another peer
			}).Catch(func(interface{}) {
				netlog.Error("failed to retrieve block for id")
			}).End()

			if common.Empty(b) {
				sendReq = true
				req.ReqBlocks.IDs = append(req.ReqBlocks.IDs, blkID)
				entry.requestTime = common.Now()
			}
			p.addPeerBlock(&entry)
		}
	} else if msg.KnownBlocks.Mode != none {
		netlog.Error("passed a notice_message with something other than a normal on none known_blocks")
		return
	}

	netlog.Debug("send req = %s\n", sendReq)
	if sendReq {
		p.write(&req)
		p.fetchWait()
		p.lastReq = &req
	}
}

func (d *dispatchManager) retryFetch(p *Peer) {
	if common.Empty(p.lastReq) {
		return
	}
	//fc_wlog( logger, "failed to fetch from ${p}",("p",c->peer_name()))
	netlog.Debug("failed to fetch from %s \n", p.peerAddr)
	var tid common.TransactionIdType
	var bid common.BlockIdType
	isTxn := false

	reqTrxCount := len(p.lastReq.ReqTrx.IDs)
	reqBlockCount := len(p.lastReq.ReqBlocks.IDs)
	if p.lastReq.ReqTrx.Mode == normal && reqTrxCount > 0 {
		isTxn = true
		tid = *p.lastReq.ReqTrx.IDs[reqTrxCount-1]
	} else if p.lastReq.ReqBlocks.Mode == normal && reqBlockCount > 0 {
		bid = *p.lastReq.ReqBlocks.IDs[reqBlockCount-1]
	} else {
		//fc_wlog( logger,"no retry, block mpde = ${b} trx mode = ${t}",
		//	("b",modes_str(c->last_req->req_blocks.mode))("t",modes_str(c->last_req->req_trx.mode)));
		netlog.Debug("no retry,block mode = %s trx mode = %s\n", modeTostring[p.lastReq.ReqBlocks.Mode], modeTostring[p.lastReq.ReqTrx.Mode])
		return
	}

	//for
	peer := &Peer{}
	if peer == p || !common.Empty(peer.lastReq) {
		//continue
	}

	sendit := false

	if isTxn {
		trx := p.trxState.getIndex("by_id").findTrxById(tid)
		if trx != nil && trx.isKnownByPeer {
			sendit = true
		}
	} else {
		blk := p.blkState.getIndex("by_id").findPeerById(bid)
		if blk != nil && blk.isKnown {
			sendit = true
		}
	}
	if sendit {
		peer.write(p.lastReq)
		peer.fetchWait()
		peer.lastReq = p.lastReq
		return
	}
	//}

	// at this point no other peer has it, re-request or do nothing?
	if p.connected() {
		p.write(p.lastReq)
		p.fetchWait()
	}
}
