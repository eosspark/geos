package net_plugin

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"net"
)

type BlockRequest struct {
	id         common.BlockIdType
	localRetry bool
}
type blockOrigin struct {
	id common.BlockIdType
	//origin connectionPtr
}
type transactionOrigin struct {
	id common.TransactionIdType
	//origin connectionPtr
}

type dispatchManager struct {
	justSendItMax        uint32
	regBlks              []BlockRequest
	reqTrx               []common.TransactionIdType
	receivedBlocks       []blockOrigin
	receivedTransactions []transactionOrigin
}

func (d *dispatchManager) bcastTransaction(msg *types.PackedTransaction) {

}

func (d *dispatchManager) rejectedTransaction(msg *common.TransactionIdType) {

}

func (d *dispatchManager) bcastBlock(msg *types.SignedBlock) {

}

func (d *dispatchManager) rejectedBlock(id *common.BlockIdType) {

}

func (d *dispatchManager) recvBlock(conn net.Conn, msg *common.BlockIdType, bnum uint32) {

}

func (d *dispatchManager) recvTransaction(conn net.Conn, id *common.TransactionIdType) {

}

func (d *dispatchManager) recvNotice(conn net.Conn, msg *NoticeMessage, generated bool) {

}

func (d *dispatchManager) retryFetch(conn net.Conn) {

}
