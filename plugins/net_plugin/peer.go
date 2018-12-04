package net_plugin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io"
	"net"
	"reflect"
	"runtime"
	"time"
)

const (
	tsBufferSize int = 32
)

type Peer struct {
	blkState      *multiIndexNet
	trxState      *multiIndexNet
	peerRequested *syncState //this peer is requesting info from us
	//socket socketPtr

	connection           net.Conn
	reader               io.Reader
	pendingMessageBuffer [1024 * 1024]byte
	//outstandingReadBytes optional<int>
	blkBuffer []byte

	nodeID             common.NodeIdType
	lastHandshakeRecv  *HandshakeMessage
	lastHandshakeSent  *HandshakeMessage
	sentHandshakeCount uint16 //int16
	connecting         bool
	syncing            bool
	protocolVersion    uint16
	peerAddr           string
	responseExpected   time.Timer
	//pendingFetch optional<request_message>
	pendingFetch *RequestMessage

	noRetry     GoAwayReason
	forkHead    common.BlockIdType
	forkHeadNum uint32
	//lastReq optional<request_message>
	lastReq *RequestMessage

	// Members set from network data
	org common.TimePoint //originate timestamp
	rec common.TimePoint //receive timestamp
	dst common.TimePoint //destination timestamp
	xmt common.TimePoint // transmit timestamp

	// Computed data
	offset float64 //peer offset double TODO

	ts [tsBufferSize]byte //working buffer for making human readable timestamps

	impl *netPluginIMpl
}

type PeerStatus struct {
	Peer          string
	Connecting    bool
	Syncing       bool
	LastHandshake HandshakeMessage
}

func NewPeer(impl *netPluginIMpl, conn net.Conn, reader io.Reader) *Peer {
	return &Peer{
		blkState:      newPeerBlockStatueIndex(),
		trxState:      newTransactionStateIndex(),
		peerRequested: new(syncState),

		connection:         conn,
		reader:             reader,
		peerAddr:           conn.RemoteAddr().String(),
		lastHandshakeSent:  &HandshakeMessage{},
		lastHandshakeRecv:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		impl:               impl,
	}
}

func (p *Peer) getStatus() *PeerStatus {
	return &PeerStatus{
		Peer:          p.peerAddr,
		Connecting:    p.connecting,
		Syncing:       p.syncing,
		LastHandshake: *p.lastHandshakeRecv,
	}
}

func (p *Peer) connected() bool {

	return false
}

func (p *Peer) current() bool {

	return false
}

func (p *Peer) reset() {
	p.peerRequested = nil //TODO
	p.blkState.clear()
	p.trxState.clear()
}

func (p *Peer) close() {

}

func (p *Peer) sendHandshake(impl *netPluginIMpl) {
	handshakePopulate(impl, p.lastHandshakeSent)
	p.sentHandshakeCount += 1
	p.lastHandshakeSent.Generation = p.sentHandshakeCount
	//fc_dlog(logger, "Sending handshake generation ${g} to ${ep}",
	//	("g",last_handshake_sent.generation)("ep", peer_name()));

	fmt.Printf("Sending handshake generation %d to %s\n", p.lastHandshakeSent.Generation, p.peerAddr)
	p.write(p.lastHandshakeSent)
}

func handshakePopulate(impl *netPluginIMpl, hello *HandshakeMessage) {
	hello.NetworkVersion = netVersionBase + netVersion
	hello.ChainID = impl.chainID
	hello.NodeID = impl.nodeID
	hello.Key = *impl.getAuthenticationKey()
	hello.Time = common.Now()
	hello.Token = crypto.Hash256(hello.Time)
	hello.Signature = *impl.signCompact(&hello.Key, &hello.Token)

	// If we couldn't sign, don't send a token.
	if common.Empty(hello.Signature) {
		hello.Token = *crypto.NewSha256Nil()
	}

	hello.P2PAddress = impl.p2PAddress + "-" + hello.NodeID.String()[:7]

	switch runtime.GOOS {
	case "darwin":
		hello.OS = "osx"
	case "linux":
		hello.OS = "linux"
	case "windows":
		hello.OS = "win32"
	default:
		hello.OS = "other"
	}
	hello.Agent = impl.userAgentName

	//controller& cc = my_impl->chain_plug->chain();
	hello.HeadID = common.BlockIdNil()
	//hello.head_num = cc.fork_db_head_block_num();
	//hello.last_irreversible_block_num = cc.last_irreversible_block_num();

	//hello.HeadNum =  cc.fork_db_head_block_num();
	//hello.LastIrreversibleBlockID = hello.LastIrreversibleBlockID

	if hello.LastIrreversibleBlockNum > 0 {
		Try(func() {
			//hello.LastIrreversibleBlockID = cc.get_block_id_for_num(hello.last_irreversible_block_num)
			//hello.LastIrreversibleBlockID =
		}).Catch(func(ex *exception.UnknownBlockException) {
			//ilog("caught unkown_block");
			fmt.Println("caught unkown_block")
			//hello.LastIrreversibleBlockNum =0
		}).End()
	}
	if hello.HeadNum > 0 {
		Try(func() {
			//hello.id = cc.get_block_id_for_num( hello.head_num )
			//hello.id =
		}).Catch(func(ex *exception.UnknownBlockException) {
			hello.HeadNum = 0
		}).End()
	}

}

func (p *Peer) PeerName() string {
	return p.peerAddr
}

//void connection::sync_wait( ) {
//response_expected->expires_from_now( my_impl->resp_expected_period);
//connection_wptr c(shared_from_this());
//response_expected->async_wait( [c]( boost::system::error_code ec){
//connection_ptr conn = c.lock();
//if (!conn) {
//// connection was destroyed before this lambda was delivered
//return;
//}
//
//conn->sync_timeout(ec);
//} );
//}
func (p *Peer) syncWait() {

}

func (p *Peer) cancelWait() {

}

func (p *Peer) fetchWait() {

}

func (p *Peer) cancelSync(reason GoAwayReason) {
	//fc_dlog(logger,"cancel sync reason = ${m}, write queue size ${o} peer ${p}",
	//	("m",reason_str(reason)) ("o", write_queue.size())("p", peer_name()));

	p.impl.log.Debug("cancel sync reason = %s, write queue size %d peer %s", ReasonToString[reason], p.write)

	p.cancelWait()
	//p.flushQueues()
	switch reason {
	case validation, fatalOther:
		p.noRetry = reason
		p.write(&GoAwayMessage{Reason: reason})
	default:
		//fc_dlog(logger, "sending empty request but not calling sync wait on ${p}", ("p",peer_name()))
		p.impl.log.Debug("sending empty request but not calling sync wait on %s", p.peerAddr)
		p.write(&SyncRequestMessage{0, 0})
	}

}

func (p *Peer) txnSendPending(impl *netPluginIMpl, ids []*common.TransactionIdType) {
	for tx := impl.localTxns.getIndex("by_id").begin(); tx.next(); {
		nts := tx.value.(*nodeTransactionState)

		if len(nts.serializedTxn) != 0 && nts.blockNum == 0 {
			found := false
			for _, known := range ids {
				if *known == nts.id {
					found = true
					break
				}
			}
			if !found {
				impl.localTxns.modify(nts, true, func(in common.ElementObject) {
					re := in.(*nodeTransactionState)
					exp := re.expires.SecSinceEpoch()
					re.expires = common.TimePointSec(exp + 1*60)
					if re.requests == 0 {
						re.trueBlock = re.blockNum
						re.blockNum = 0
					}
					re.requests += 1
					if re.requests == 0 {
						re.blockNum = re.trueBlock
					}
				})

				//queue_write(std::make_shared<vector<char>>(tx->serialized_txn),
				//	true,
				//	[tx_id=tx->id](boost::system::error_code ec, std::size_t ) {
				//auto& local_txns = my_impl->local_txns;
				//auto tx = local_txns.get<by_id>().find(tx_id);
				//if (tx != local_txns.end()) {
				//local_txns.modify(tx, decr_in_flight);
				//} else {
				//fc_wlog(logger, "Local pending TX erased before queued_write called callback");
				//}
				//});

			}
		}
	}

}

func (p *Peer) txnSend(impl *netPluginIMpl, ids []common.TransactionIdType) {
	for _, t := range ids {
		tx := impl.localTxns.getIndex("by_id").findLocalTrxById(t)
		if tx != nil && len(tx.serializedTxn) != 0 {
			//my_impl->local_txns.modify( tx,incr_in_flight);
			impl.localTxns.modify(tx, true, func(in common.ElementObject) {
				re := in.(*nodeTransactionState)
				exp := re.expires.SecSinceEpoch()
				re.expires = common.TimePointSec(exp + 1*60)
				if re.requests == 0 {
					re.trueBlock = re.blockNum
					re.blockNum = 0
				}
				re.requests += 1
				if re.requests == 0 {
					re.blockNum = re.trueBlock
				}
			})
			//queue_write(std::make_shared<vector<char>>(tx->serialized_txn),
			//		 true,
			//		 [t](boost::system::error_code ec, std::size_t ) {
			//			auto& local_txns = my_impl->local_txns;
			//			auto tx = local_txns.get<by_id>().find(t);
			//			if (tx != local_txns.end()) {
			//			   local_txns.modify(tx, decr_in_flight);
			//			} else {
			//			   fc_wlog(logger, "Local TX erased before queued_write called callback");
			//			}
			//		 });
		}
	}

}

// void connection::blk_send_branch() {
//       controller &cc = my_impl->chain_plug->chain();
//       uint32_t head_num = cc.fork_db_head_block_num ();
//       notice_message note;
//       note.known_blocks.mode = normal;
//       note.known_blocks.pending = 0;
//       fc_dlog(logger, "head_num = ${h}",("h",head_num));
//       if(head_num == 0) {
//          enqueue(note);
//          return;
//       }
//       block_id_type head_id;
//       block_id_type lib_id;
//       uint32_t lib_num;
//       try {
//          lib_num = cc.last_irreversible_block_num();
//          lib_id = cc.last_irreversible_block_id();
//          head_id = cc.fork_db_head_block_id();
//       }
//       catch (const assert_exception &ex) {
//          elog( "unable to retrieve block info: ${n} for ${p}",("n",ex.to_string())("p",peer_name()));
//          enqueue(note);
//          return;
//       }
//       catch (const fc::exception &ex) {
//       }
//       catch (...) {
//       }

//       vector<signed_block_ptr> bstack;
//       block_id_type null_id;
//       for (auto bid = head_id; bid != null_id && bid != lib_id; ) {
//          try {
//             signed_block_ptr b = cc.fetch_block_by_id(bid);
//             if ( b ) {
//                bid = b->previous;
//                bstack.push_back(b);
//             }
//             else {
//                break;
//             }
//          } catch (...) {
//             break;
//          }
//       }
//       size_t count = 0;
//       if (!bstack.empty()) {
//          if (bstack.back()->previous == lib_id) {
//             count = bstack.size();
//             while (bstack.size()) {
//                enqueue(*bstack.back());
//                bstack.pop_back();
//             }
//          }
//          fc_ilog(logger, "Sent ${n} blocks on my fork",("n",count));
//       } else {
//          fc_ilog(logger, "Nothing to send on fork request");
//       }

//       syncing = false;
//    }

func (p *Peer) blk_send_branch() {

	//controller &cc = my_impl->chain_plug->chain();
	//uint32_t head_num = cc.fork_db_head_block_num ();
	var headNum uint32

	var note NoticeMessage

	note.KnownBlocks.Mode = normal
	note.KnownBlocks.Pending = 0
	log.Debug("head_num = %d", headNum)
	if headNum == 0 {
		//enqueue(note)TODO
		return
	}

	var headID common.BlockIdType
	var libID common.BlockIdType
	//var libNum uint32
	Try(func() {
		//libNum = cc.last_irreversible_block_num()
		//libID = cc.last_irreversible_block_id()
		//headID = cc.fork_db_head_block_id()
	}).Catch(func(ex *exception.AssertException) {
		log.Error("unable to retrieve block info: %s for %s", ex.What(), p.peerAddr)
		p.write(&note)
		return
	}).Catch(func(ex exception.Exception) {

	}).Catch(func(interface{}) {}).End()

	var bstack []*types.SignedBlock
	var nullID common.BlockIdType
	breaking := false
	for bid := headID; bid != nullID && bid != libID; {
		Try(func() {
			var b *types.SignedBlock
			//b :=cc.fetch_block_by_id(bid)
			if b != nil {
				bid = b.Previous
				bstack = append(bstack, b)
			} else {
				breaking = true
				return
			}
		}).Catch(func(interface{}) {
			breaking = true
			return
		}).End()

		if breaking {
			break
		}
	}

	//if (!bstack.empty()) {
	//   if (bstack.back()->previous == lib_id) {
	//      count = bstack.size();
	//      while (bstack.size()) {
	//         enqueue(*bstack.back());
	//         bstack.pop_back();
	//      }
	//   }
	//   fc_ilog(logger, "Sent ${n} blocks on my fork",("n",count));
	//}
	count := len(bstack)
	if count != 0 {
		if bstack[count-1].Previous == libID {
			for count != 1 {
				p.write(&SignedBlockMessage{*bstack[count-1]})
				count -= 1
			}
		}
		log.Info("Sent %d blocks on my fork", len(bstack))
	} else {
		log.Info("Nothing to send on fork request")
	}
	p.syncing = false

}

func (p *Peer) blkSend(impl *netPluginIMpl, ids []common.BlockIdType) {
	//    controller &cc = my_impl->chain_plug->chain();
	var count int
	breaking := false
	for _, blkid := range ids {
		count += 1
		Try(func() {
			var b *types.SignedBlock
			//signed_block_ptr b = cc.fetch_block_by_id(blkid);
			if b != nil {
				p.impl.log.Debug("found block for id ar num %d", b.BlockNumber())
				//enqueue(net_message(*b))//TODO
				p.write(&SignedBlockMessage{*b})
			} else {
				p.impl.log.Info("fetch block by id returned null, id %s on block %d of %d for %s",
					blkid, count, len(ids), p.peerAddr)
				breaking = true
			}
		}).Catch(func(ex *exception.AssertException) {
			p.impl.log.Error("caught assert on fetch_block_by_id, %s, id %s on block %d of %d for %s",
				ex.What(), blkid, count, len(ids), p.peerAddr)
			breaking = true
		}).Catch(func(interface{}) {
			p.impl.log.Error("caught others exception fetching block id %s on block %d of %d for %s",
				blkid, count, len(ids), p.peerAddr)
			breaking = true
		})
		if breaking {
			break
		}
	}
}

func (p *Peer) stopSend() {
	p.syncing = false
}

//void enqueue( const net_message &msg, bool trigger_send = true );
//void cancel_sync(go_away_reason);
//void flush_queues();
//bool enqueue_sync_block();

func (p *Peer) requestSyncBlocks(start, end uint32) {
	syncRequest := SyncRequestMessage{
		StartBlock: start,
		EndBlock:   end,
	}
	p.write(&syncRequest)
	p.syncWait()
}

//void sync_timeout(boost::system::error_code ec);
//void fetch_timeout(boost::system::error_code ec);
//
//void queue_write(std::shared_ptr<vector<char>> buff,
//bool trigger_send,
//std::function<void(boost::system::error_code, std::size_t)> callback);
//void do_queue_write();

func isValid(msg *HandshakeMessage) bool {
	// Do some basic validation of an incoming handshake_message, so things
	// that really aren't handshake messages can be quickly discarded without
	// affecting state.
	valid := true
	if msg.LastIrreversibleBlockNum > msg.HeadNum {
		fmt.Printf("Handshake message validation: last irreversible block %d is greater than head block %d\n",
			msg.LastIrreversibleBlockNum, msg.HeadNum)
		valid = false
	}
	if len(msg.P2PAddress) == 0 {
		fmt.Println("Handshake message validation: p2p_address is null string")
		valid = false
	}
	if len(msg.OS) == 0 {
		fmt.Println("Handshake message validation: os field is null string")
		valid = false
	}
	if (common.CompareString(msg.Signature, ecc.NewSigNil()) != 0 || common.CompareString(msg.Token, crypto.NewSha256Nil()) != 0) &&
		common.CompareString(msg.Token, crypto.Hash256(msg.Time)) != 0 {
		fmt.Println("Handshake message validation: token field invalid")
		valid = false
	}

	return valid
}

//sendTime populate and queue time_message immediately using incoming time_message
func (p *Peer) sendTime(msg *TimeMessage) {
	xpkt := &TimeMessage{
		Org: msg.Org,
		Rec: msg.Dst,
		Xmt: common.Now(),
	}
	p.write(xpkt)
}

func (p *Peer) sendTimeTicker() {
	xpkt := &TimeMessage{
		Org: p.rec,
		Rec: p.dst,
		Xmt: common.Now(),
	}
	p.org = xpkt.Xmt
	p.write(xpkt)
}

func (p *Peer) addPeerBlock(entry *peerBlockState) bool {
	bptr := p.blkState.getIndex("by_id").findPeerById(entry.id)
	added := bptr == nil
	if added {
		p.blkState.insertPeerBlock(entry)
	} else {
		p.blkState.modify(bptr, false, func(in common.ElementObject) {
			re := in.(*peerBlockState)
			re.isKnown = true
		})

		if entry.blockNum == 0 {
			p.blkState.modify(bptr, true, func(in common.ElementObject) {
				re := in.(*peerBlockState)
				re.blockNum = entry.blockNum
			})
		} else {
			p.blkState.modify(bptr, false, func(in common.ElementObject) {
				re := in.(*peerBlockState)
				re.requestTime = common.Now()
			})
		}
	}

	return added
}

func (p *Peer) read(impl *netPluginIMpl) {
	defer func() {
		p.connection.Close()
		impl.loopWG.Done()
	}()

	impl.loopWG.Add(1)
	p.impl.log.Info("start read message!")

	for {
		p2pMessage, err := ReadP2PMessageData(p.reader)
		if err != nil {
			fmt.Println("Error reading from p2p client:", err)
			continue
		}
		time.Sleep(100 * time.Millisecond) //TODO for testing
		//data, err := json.Marshal(p2pMessage)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println(p.peerAddr, ": Receive P2PMessag ", string(data))

		switch msg := p2pMessage.(type) {
		case *HandshakeMessage:
			impl.handleHandshakeMsg(p, msg)
		case *ChainSizeMessage:
			impl.handleChainSizeMsg(p, msg)
		case *GoAwayMessage:
			impl.handleGoawayMsg(p, msg)
			fmt.Printf("GO AWAY Reason[%d] \n", msg.Reason)
		case *TimeMessage:
			impl.handleTimeMsg(p, msg)
		case *NoticeMessage:
			impl.handleNoticeMsg(p, msg)
		case *RequestMessage:
			impl.handleRequestMsg(p, msg)
		case *SyncRequestMessage:
			impl.handleSyncRequestMsg(p, msg)
		case *SignedBlockMessage:
			impl.handleSignedBlock(p, msg)
		case *PackedTransactionMessage:
			impl.handlePackTransaction(p, msg)
		default:
			fmt.Println("unsupport p2pmessage type")
		}

	}
}

func ReadP2PMessageData(r io.Reader) (p2pMessage P2PMessage, err error) {
	//data := make([]byte, 0)
	lengthBytes := make([]byte, 4, 4)
	_, err = io.ReadFull(r, lengthBytes)
	if err != nil {
		return
	}
	//data = append(data, lengthBytes...)
	size := binary.LittleEndian.Uint32(lengthBytes)
	payloadBytes := make([]byte, size, size)
	count, err := io.ReadFull(r, payloadBytes)
	if count != int(size) {
		err = fmt.Errorf("readfull not full read [%d] expected[%d]", count, size)
		return
	}
	if err != nil {
		fmt.Println("readfull ,error:", err)
		return
	}
	//data = append(data, payloadBytes...)
	//fmt.Printf("receive data:  %#v\n", data)

	messagetype := P2PMessageType(payloadBytes[0])
	attr, ok := messagetype.reflectTypes()
	if !ok {
		fmt.Errorf("decode, unknown p2p message type [%d]", messagetype)
		return
	}
	msg := reflect.New(attr.ReflectType)
	err = rlp.DecodeBytes(payloadBytes[1:], msg.Interface())
	if err != nil {
		return nil, err
	}

	p2pMessage = msg.Interface().(P2PMessage)

	return
}

func (p *Peer) write(message P2PMessage) {
	payload, err := rlp.EncodeToBytes(message)
	if err != nil {
		err = fmt.Errorf("p2p message, %s", err)
		return
	}
	messageLen := uint32(len(payload) + 1)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, messageLen)
	sendBuf := append(buf, byte(message.GetType()))
	sendBuf = append(sendBuf, payload...)

	p.connection.Write(sendBuf)

	//fmt.Println(p.peerAddr, ": Message bytes", sendBuf)
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	p.impl.log.Info("%s: send Message json: %s", p.peerAddr, string(data))

	return
}

func toProtocolVersion(v uint16) uint16 {
	if v >= netVersionBase {
		v -= netVersionBase
		if v <= netVersionRange {
			return v
		}
	}
	return 0
}
