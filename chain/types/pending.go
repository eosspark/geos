/*
 *  @Time : 2018/8/26 下午3:02
 *  @Author : xueyahui
 *  @File : pending
 *  @Software: GoLand
 */
package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
)

type ActionReceipt struct {
	receiver       common.AccountName
	ActDigest      common.SHA256Bytes
	GlobalSequence uint64
	RecvSequence   uint64
	AuthSequence   map[common.AccountName]uint64
	CodeSequence   uint32 //TODO
	ABISequence    uint32
}
type PendingState struct {
	DBSeesion   eosiodb.Session
	BlockState  BlockState
	Actions     []ActionReceipt
	BlockStatus BlockStatus
}
