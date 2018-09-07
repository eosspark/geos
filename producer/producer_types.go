package producer_plugin

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"time"
)

//enums
type EnumPendingBlockMode int

const (
	producing = EnumPendingBlockMode(iota)
	speculating
)

type EnumStartBlockRusult int

const (
	succeeded = EnumStartBlockRusult(iota)
	failed
	waiting
	exhausted
)

//timer
//type scheduleTimer struct {
//	internal *time.Timer
//	duration time.Duration
//}
//
//func (pt *scheduleTimer) expiresFromNow(m common.Microseconds) {
//	pt.duration = time.Microsecond * time.Duration(m)
//}
//
//func (pt *scheduleTimer) expiresUntil(t common.TimePoint) {
//	pt.expiresFromNow(t.Sub(common.Now()))
//}
//
//func (pt *scheduleTimer) expiresAt(epoch common.Microseconds) {
//	pt.expiresUntil(common.TimePoint(epoch))
//}
//
//func (pt *scheduleTimer) asyncWait(valid func() bool, call func()) {
//	pt.internal = time.NewTimer(pt.duration)
//	<-pt.internal.C
//	if valid() {
//		go call()
//	}
//}
//
//func (pt *scheduleTimer) cancel() {
//	if pt.internal != nil {
//		pt.internal.Stop()
//		pt.internal = nil
//	}
//}

type signatureProviderType func([]byte) ecc.Signature

type transactionIdWithExpireIndex map[common.TransactionIDType]common.TimePoint

type respVariant struct {
	err error
	trx int //TODO:transaction_trace_ptr

}

type RuntimeOptions struct {
	MaxTransactionTime      int32
	MaxIrreversibleBlockAge int32
	ProduceTimeOffsetUs     int32
	LastBlockTimeOffsetUs   int32
	SubjectiveCpuLeewayUs   int32
	IncomingDeferRadio      float64
}

type WhitelistAndBlacklist struct {
	ActorWhitelist    map[common.AccountName]struct{}
	ActorBlacklist    map[common.AccountName]struct{}
	ContractWhitelist map[common.AccountName]struct{}
	ContractBlacklist map[common.AccountName]struct{}
	ActionBlacklist   map[[2]common.Name]struct{}
	KeyBlacklist      map[ecc.PublicKey]struct{}
}

type GreylistParams struct {
	Accounts []common.AccountName
}

type tuple struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(respVariant)
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key ecc.PrivateKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest []byte) (sign ecc.Signature) {
		sign, err = key.Sign(digest)
		return
	}
	return
}

func makeKeosdSignatureProvider(produce *ProducerPlugin, url string, pubKey ecc.PublicKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest []byte) ecc.Signature {
		if produce != nil {
			//TODO
			return ecc.Signature{}
		} else {
			return ecc.Signature{}
		}
	}
	return
}

//errors
var ErrProducerFail = errors.New("called produce_block while not actually producing")
var ErrMissingPendingBlockState = errors.New("pending_block_state does not exist but it should, another plugin may have corrupted it")
var ErrProducerPriKeyNotFound = errors.New("attempting to produce a block for which we don't have the private key")
var ErrBlockFromTheFuture = errors.New("received a block from the future, ignoring it")
