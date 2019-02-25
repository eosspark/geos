package net_plugin

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
)

type stages byte

const (
	libCatchup = stages(iota)
	headCatchup
	inSync
)

type syncManager struct {
	syncKnownLibNum      uint32
	syncLastRequestedNum uint32
	syncNextExpectedNum  uint32
	syncReqSpan          uint32
	source               *Connection
	state                stages
	chainPlugin          *chain_plugin.ChainPlugin
	myImpl               *netPluginIMpl
}

func NewSyncManager(impl *netPluginIMpl, span uint32) *syncManager {
	s := &syncManager{
		syncKnownLibNum:      0,
		syncLastRequestedNum: 0,
		syncNextExpectedNum:  1,
		syncReqSpan:          span,
		state:                inSync,
		source:               &Connection{},
		myImpl:               impl,
	}
	s.chainPlugin = app.App().FindPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin)
	EosAssert(s.chainPlugin != nil, &exception.MissingChainPluginException{}, "")
	return s
}

func stageStr(s stages) string {
	switch s {
	case libCatchup:
		return "lib catchup"
	case headCatchup:
		return "head catchup"
	case inSync:
		return "in sync"
	default:
		return "unkown"
	}
}

func (s *syncManager) setStage(newState stages) {
	if s.state == newState {
		return
	}
	FcLog.Info("old state %s becoming %s", stageStr(s.state), stageStr(newState))
	s.state = newState
}

func (s *syncManager) isActive(c *Connection) bool {
	if s.state == headCatchup && c != nil {
		fhSet := c.forkHead != common.BlockIdNil()
		FcLog.Info("fork_head_num = %d fork_head set = %s", c.forkHeadNum, fhSet)

		return c.forkHead != common.BlockIdNil() && c.forkHeadNum < s.chainPlugin.Chain().ForkDbHeadBlockNum()
	}
	return s.state != inSync
}

func (s *syncManager) resetLibNum(c *Connection) {
	if s.state == inSync {
		s.source.reset()
	}
	if c.current() {
		if c.lastHandshakeRecv.LastIrreversibleBlockNum > s.syncKnownLibNum {
			s.syncKnownLibNum = c.lastHandshakeRecv.LastIrreversibleBlockNum
		}
	} else if c == s.source {
		s.syncLastRequestedNum = 0
		s.requestNextChunk(c)
	}
}

func (s *syncManager) syncRequired() bool {
	FcLog.Info("last req = %d, last recv = %d, known = %d, our head %d",
		s.syncLastRequestedNum, s.syncNextExpectedNum, s.syncKnownLibNum, s.chainPlugin.Chain().ForkDbHeadBlockNum())
	return s.syncLastRequestedNum < s.syncKnownLibNum || s.chainPlugin.Chain().ForkDbHeadBlockNum() < s.syncLastRequestedNum
}

func findConnection(peer string) {

}

func (s *syncManager) requestNextChunk(conn *Connection) {
	headBlock := s.chainPlugin.Chain().ForkDbHeadBlockNum()

	if headBlock < s.syncLastRequestedNum && s.source != nil && s.source.current() {
		FcLog.Info("ignoring request, head is %d last req = %d source is %s", headBlock, s.syncLastRequestedNum, s.source.peerAddr)
		return
	}

	/* ----------
	* next chunk provider selection criteria
	* a provider is supplied and able to be used, use it.
	* otherwise select the next available from the list, round-robin style.
	 */

	if conn != nil && conn.current() {
		s.source = conn
	} else {
		if len(s.myImpl.connections) == 1 {
			if s.source == nil {
				s.source = s.myImpl.connections[0]
			}
		} else {
			if len(s.myImpl.connections) == 1 {
				if s.source == nil {
					s.source = s.myImpl.connections[0]
				}
			} else {
				cptr := 0
				end := len(s.myImpl.connections) - 1
				cend := end
				if s.source != nil {
					cptr, _ = s.myImpl.findConnection(s.source.peerAddr)
					cend = cptr
					if cptr == -1 {
						//not there - must have been closed! cend is now connections.end, so just flatten the ring.
						s.source.reset()
						cptr = 0
					} else {
						//was found - advance the start to the next. cend is the old source.
						if cptr+1 == end && cend != end {
							cptr = 0
						}
					}
				}
				cstartIt := cptr
				for {
					//select the first one which is current and break out.
					if s.myImpl.connections[cptr].current() {
						s.source = s.myImpl.connections[cptr]
						break
					}

					if cptr == end {
						cptr = 0
					} else {
						cptr++
					}
					if cstartIt == cptr {
						break
					}
				}
				// no need to check the result, either source advanced or the whole list was checked and the old source is reused.
			}
		}
	}

	// verify there is an available source
	if s.source == nil || !s.source.current() {
		netLog.Error("Unable to continue syncing at this time")
		s.syncLastRequestedNum = s.chainPlugin.Chain().LastIrreversibleBlockNum()
		s.setStage(inSync) // probably not, but we can't do anything else
		return
	}

	if s.syncLastRequestedNum != s.syncKnownLibNum {
		start := s.syncNextExpectedNum
		end := start + s.syncReqSpan - 1
		if end > s.syncKnownLibNum {
			end = s.syncKnownLibNum
		}
		if end > 0 && end >= start {
			FcLog.Info("requesting range %s to %d, from %d\n", s.source.peerAddr, end, start)
			s.source.requestSyncBlocks(start, end)
			s.syncLastRequestedNum = end
		}
	}
}

func (s *syncManager) sendHandshakes(impl *netPluginIMpl) {
	for _, ci := range impl.connections {
		if ci.current() {
			ci.sendHandshake()
		}
	}
}

func (s *syncManager) recvHandshake(c *Connection, msg *HandshakeMessage) {
	cc := s.chainPlugin.Chain()
	libNum := cc.LastIrreversibleBlockNum()
	peerLib := msg.LastIrreversibleBlockNum
	s.resetLibNum(c)
	c.syncing = false

	//--------------------------------
	// sync need checks; (lib == last irreversible block)
	//
	// 0. my head block id == peer head id means we are all caugnt up block wise
	// 1. my head block num < peer lib - start sync locally
	// 2. my lib > peer head num - send an last_irr_catch_up notice if not the first generation
	//
	// 3  my head block num <= peer head block num - update sync state and send a catchup request
	// 4  my head block num > peer block num ssend a notice catchup if this is not the first generation
	//
	//-----------------------------

	head := cc.ForkDbHeadBlockNum()
	headID := cc.ForkDbHeadBlockId()

	if headID == msg.HeadID {
		FcLog.Info("sync check statue 0")
		// notify peer of our pending transactions

		note := NoticeMessage{}
		note.KnownBlocks.Mode = none
		note.KnownTrx.Mode = catchUp
		note.KnownBlocks.Pending = uint32(s.myImpl.localTxns.Size())
		c.enqueue(&note, true)
		return
	}

	if head < peerLib {
		FcLog.Info("sync check state 1")
		//wait for receipt of a notice message before initiating sync
		if c.protocolVersion < protoExplicitSync {
			s.startSync(c, peerLib)
		}
		return
	}

	if libNum > msg.HeadNum {
		FcLog.Info("sync check state 2")
		if msg.Generation > 1 || c.protocolVersion > protoBase {
			note := NoticeMessage{}
			note.KnownBlocks.Mode = lastIrrCatchUp
			note.KnownBlocks.Pending = head
			note.KnownTrx.Mode = lastIrrCatchUp
			note.KnownTrx.Pending = libNum
			c.enqueue(&note, true)
		}
		c.syncing = true
		return
	}

	if head <= msg.HeadNum {
		FcLog.Info("sync check state 3")
		s.verifyCatchup(c, msg.HeadNum, msg.HeadID)
		return
	} else {
		FcLog.Info("sync check state 4")
		if msg.Generation > 1 || c.protocolVersion > protoBase {
			note := NoticeMessage{}
			note.KnownBlocks.Mode = catchUp
			note.KnownBlocks.Pending = head
			note.KnownBlocks.IDs = append(note.KnownBlocks.IDs, headID)
			note.KnownTrx.Mode = none
			c.enqueue(&note, true)
		}
		c.syncing = true
		return
	}

	//netLog.Error("sync check failed to resolve status")
}

func (s *syncManager) startSync(c *Connection, target uint32) {
	if target > s.syncKnownLibNum {
		s.syncKnownLibNum = target
	}
	if !s.syncRequired() {
		bNum := s.myImpl.ChainPlugin.Chain().LastIrreversibleBlockNum()
		hNum := s.myImpl.ChainPlugin.Chain().ForkDbHeadBlockNum()
		FcLog.Info("we are already caught up, my irr = %d,head =%d,target = %d", bNum, hNum, target)
		return
	}
	if s.state == inSync {
		s.setStage(libCatchup)
		s.syncNextExpectedNum = s.myImpl.ChainPlugin.Chain().LastIrreversibleBlockNum() + 1
	}

	FcLog.Info("Catching up with chain, our last req is %d, theirs is %d peer %s", s.syncLastRequestedNum, target, c.peerAddr)

	s.requestNextChunk(c)
}

func (s *syncManager) reassignFetch(c *Connection, reason GoAwayReason) {
	FcLog.Info("reassign_fetch, our last req is %d, next expected is %d peer %s", s.syncLastRequestedNum, s.syncNextExpectedNum, c.peerAddr)
	if c == s.source {
		c.cancelSync(reason)
		s.syncLastRequestedNum = 0
		s.requestNextChunk(c)
	}
}

func (s *syncManager) verifyCatchup(c *Connection, num uint32, id common.BlockIdType) {
	req := RequestMessage{}
	req.ReqBlocks.Mode = catchUp

	for _, peer := range s.myImpl.connections {
		if peer.forkHead == id || peer.forkHeadNum > num {
			req.ReqBlocks.Mode = none
		}
		break
	}

	if req.ReqBlocks.Mode == catchUp {
		c.forkHead = id
		c.forkHeadNum = num
		FcLog.Info("got a catch_up notice while in %s, fork head num = %d target LIB = %d next_expected = %d",
			stageStr(s.state), num, s.syncKnownLibNum, s.syncNextExpectedNum)
		if s.state == libCatchup {
			return
		}
		s.setStage(headCatchup)

	} else {
		c.forkHead = common.BlockIdNil()
		c.forkHeadNum = 0
	}

	req.ReqTrx.Mode = none
	c.enqueue(&req, true)
}

func (s *syncManager) recvNotice(c *Connection, msg *NoticeMessage) {
	FcLog.Info("sync_manager got %s block notice", modeStr[msg.KnownBlocks.Mode])
	if msg.KnownBlocks.Mode == catchUp {
		IDsCount := len(msg.KnownBlocks.IDs)
		if IDsCount == 0 {
			netLog.Error("got a catch up with ids size = 0")
		} else {
			s.verifyCatchup(c, msg.KnownBlocks.Pending, msg.KnownBlocks.IDs[IDsCount-1])
		}
	} else {
		c.lastHandshakeRecv.LastIrreversibleBlockNum = msg.KnownTrx.Pending
		s.resetLibNum(c)
		s.startSync(c, msg.KnownBlocks.Pending)
	}
}

func (s *syncManager) rejectedBlock(c *Connection, blkNum uint32) {
	if s.state != inSync {
		FcLog.Debug("block %d not accepted from %s", blkNum, c.peerAddr)
		s.syncLastRequestedNum = 0
		s.source.reset()
		s.myImpl.close(c)
		s.setStage(inSync)
		s.sendHandshakes(s.myImpl)
	}
}

func (s *syncManager) recvBlock(c *Connection, blkID common.BlockIdType, blkNum uint32) {
	FcLog.Debug("got block %d from %s,state:%s", blkNum, c.peerAddr, stageStr(s.state))
	if s.state == libCatchup {
		if blkNum != s.syncNextExpectedNum {
			FcLog.Info("expected block %d but got %d", s.syncNextExpectedNum, blkNum)
			s.myImpl.close(c)
			return
		}
		s.syncNextExpectedNum = blkNum + 1
	}

	if s.state == headCatchup {
		FcLog.Debug("sync_manager in head_catchup state")
		s.setStage(inSync)
		s.source.reset()

		nullID := crypto.NewSha256Nil()
		for _, cp := range s.myImpl.connections {
			if cp.forkHead.Equals(nullID) {
				continue
			}
			if cp.forkHead.Equals(blkID) || cp.forkHeadNum < blkNum {
				cp.forkHead = nullID
				cp.forkHeadNum = 0
			} else {
				s.setStage(headCatchup)
			}
		}
	} else if s.state == libCatchup {
		if blkNum == s.syncKnownLibNum {
			FcLog.Debug("All caught up with last known last irreversible block resending handshake")
			s.setStage(inSync)
			s.sendHandshakes(s.myImpl)
		} else if blkNum == s.syncLastRequestedNum {
			s.requestNextChunk(nil)
		} else {
			//FcLog.Debug("calling sync_wait on connecting %s", c.peerAddr)
			//c.syncWait()
		}
	}
}
