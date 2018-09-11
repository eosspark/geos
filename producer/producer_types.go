package producer_plugin

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
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

type signatureProviderType func(sha256 rlp.Sha256) ecc.Signature

type transactionIdWithExpireIndex map[common.TransactionIDType]common.TimePoint

type ErrorORTrace struct {
	error error
	trace *types.TransactionTrace
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

type pendingIncomingTransaction struct {
	packedTransaction   *types.PackedTransaction
	persistUntilExpired bool
	next                func(ErrorORTrace)
}

func failureIsSubjective(e error, deadlineIsSubjective bool) bool {
	//TODO wait for error definition
	return false
}

func makeDebugTimeLogger() func() {
	start := time.Now()
	return func() {
		fmt.Println(time.Now().Sub(start))
	}
}

func makeKeySignatureProvider(key ecc.PrivateKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest rlp.Sha256) (sign ecc.Signature) {
		sign, err = key.Sign(digest.Bytes())
		return
	}
	return
}

func makeKeosdSignatureProvider(produce *ProducerPlugin, url string, pubKey ecc.PublicKey) (signFunc signatureProviderType, err error) {
	signFunc = func(digest rlp.Sha256) ecc.Signature {
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
