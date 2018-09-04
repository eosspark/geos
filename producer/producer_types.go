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
type scheduleTimer struct {
	internal *time.Timer
	duration time.Duration
}

func (pt *scheduleTimer) expiresFromNow(d time.Duration) {
	pt.duration = d
}

func (pt *scheduleTimer) expiresUntil(t time.Time) {
	pt.expiresFromNow(time.Until(t))
}

func (pt *scheduleTimer) expiresAt(epoch int64) {
	pt.expiresUntil(time.Unix(0, epoch*1e3))
}

func (pt *scheduleTimer) asyncWait(valid func() bool, call func()) {
	pt.internal = time.NewTimer(pt.duration)
	<-pt.internal.C
	if valid() {
		go call()
	}
}

func (pt *scheduleTimer) cancel() {
	if pt.internal != nil {
		pt.internal.Stop()
		pt.internal = nil
	}
}

type signatureProviderType func([]byte) ecc.Signature

type transactionIdWithExpireIndex map[common.TransactionIDType]time.Time

type respVariant struct {
	err error
	trx int //TODO:transaction_trace_ptr

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
