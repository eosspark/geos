package types

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/log"
)

type TransactionHeader struct {
	Expiration     common.TimePointSec `json:"expiration"`
	RefBlockNum    uint16              `json:"ref_block_num"`
	RefBlockPrefix uint32              `json:"ref_block_prefix"`

	MaxNetUsageWords uint32 `json:"max_net_usage_words"`
	MaxCpuUsageMS    uint8  `json:"max_cpu_usage_ms"`
	DelaySec         uint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration
}

func (th TransactionHeader) GetRefBlocknum(headBlocknum uint32) uint32 {
	return headBlocknum/0xffff*0xffff + headBlocknum%0xffff
}

func (th TransactionHeader) VerifyReferenceBlock(referenceBlock common.BlockIdType) bool {
	return th.RefBlockNum == uint16(common.EndianReverseU32(uint32(referenceBlock.Hash[0]))) &&
		th.RefBlockPrefix == uint32(referenceBlock.Hash[1])
}

func (th TransactionHeader) Validate() {
	if th.MaxNetUsageWords >= uint32(0xffffffff)/8 {
		panic("declared max_net_usage_words overflows when expanded to max net usage")
	}
}

type Transaction struct { // WARN: is a `variant` in C++, can be a SignedTransaction or a Transaction.
	TransactionHeader

	ContextFreeActions []*Action    `json:"context_free_actions"`
	Actions            []*Action    `json:"actions"`
	Extensions         []*Extension `json:"transaction_extensions"`
}

// NewTransaction creates a transaction. Unless you plan on adding HeadBlockID later, to be complete, opts should contain it.  Sign
func NewTransaction(actions []*Action, opts *TxOptions) *Transaction {
	if opts == nil {
		opts = &TxOptions{}
	}

	tx := &Transaction{Actions: actions}
	tx.Fill(opts.HeadBlockID, opts.DelaySecs, opts.MaxNetUsageWords, opts.MaxCpuUsageMS)
	return tx
}

func (tx *Transaction) TotalActions() uint32 {
	return uint32(len(tx.ContextFreeActions) + len(tx.Actions))
}

func (tx *Transaction) FirstAuthorizor() common.AccountName {
	for _, a := range tx.Actions {
		for _, auth := range a.Authorization {
			return auth.Actor
		}
	}
	return common.AccountName(0)
}
func (tx *Transaction) SetExpiration(in uint32) {
	tx.Expiration = common.TimePointSec(in)
}

func (tx *Transaction) GetSignatureKeys(chainId common.ChainIdType, allowDeplicateKeys bool, useCache bool) []ecc.PublicKey {
	//TODO
	return nil
}

type Extension struct {
	Type uint16          `json:"type"`
	Data common.HexBytes `json:"data"`
}

// Fill sets the fields on a transaction.  If you pass `headBlockID`, then `api` can be nil. If you don't pass `headBlockID`, then the `api` is going to be called to fetch

//canada eos code
func (tx *Transaction) Fill(headBlockID common.BlockIdType, delaySecs, maxNetUsageWords uint32, maxCPUUsageMS uint8) {
	tx.setRefBlock(headBlockID)

	if tx.ContextFreeActions == nil {
		tx.ContextFreeActions = make([]*Action, 0, 0)
	}
	if tx.Extensions == nil {
		tx.Extensions = make([]*Extension, 0, 0)
	}

	tx.MaxNetUsageWords = uint32(maxNetUsageWords)
	tx.MaxCpuUsageMS = maxCPUUsageMS
	tx.DelaySec = uint32(delaySecs)

	//tx.SetExpiration(30 * time.Second)
	tx.SetExpiration(30)
}

func (tx *Transaction) setRefBlock(blockID common.BlockIdType) {
	tx.RefBlockNum = uint16(blockID.Hash[0])
	tx.RefBlockPrefix = uint32(blockID.Hash[1])
}

type SignedTransaction struct {
	Transaction

	Signatures      []ecc.Signature   `json:"signatures"`
	ContextFreeData []common.HexBytes `json:"context_free_data"`

	packed *PackedTransaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     *tx,
		Signatures:      make([]ecc.Signature, 0),
		ContextFreeData: make([]common.HexBytes, 0),
	}
}

func (s *SignedTransaction) String() string {

	data, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (st *SignedTransaction) GetSignatureKeys(chainId common.ChainIdType, allowDeplicateKeys bool, useCache bool) []ecc.PublicKey {
	//TODO

	return st.Transaction.GetSignatureKeys(chainId, allowDeplicateKeys, useCache)
}

func (ptrx *PackedTransaction) GetSignedTransaction() *SignedTransaction {

	switch ptrx.Compression {
	case common.CompressionNone:
		return &SignedTransaction{GetTransaction(), ptrx.Signatures, UnpackContextFreeData(ptrx.PackedContextFreeData), nil}
	case common.CompressionZlib:
		return &SignedTransaction{}
	default:
		return nil
	}
	return nil //TODO
}

func UnpackContextFreeData(data []byte) []common.HexBytes {
	t := []common.HexBytes{}
	if len(data) == 0 {
		return t
	}
	err := rlp.DecodeBytes(data, t)
	if err != nil {
		fmt.Println("UnpackContextFreeData is error :", err.Error())
	}
	return t
}

func ZlibDecompressContextFreeData(data []byte) []common.HexBytes {
	t := []common.HexBytes{}
	if len(data) == 0 {
		return t
	}

	//out := ZlibDecompress()	//TODO
	return nil
}

func ZlibDecompress() []common.HexBytes { return nil }

func (head *TransactionHeader) SetReferenceBlock(referenceBlock common.BlockIdType) {
	first := common.EndianReverseU32(uint32(referenceBlock.Hash[0]))
	head.RefBlockNum = uint16(first)
	head.RefBlockPrefix = uint32(referenceBlock.Hash[1])
	log.Info("SetReferenceBlock:", head)
}

// func (s *SignedTransaction) SignedByKeys(chainID SHA256Bytes) (out []ecc.PublicKey, err error) {
// 	trx, cfd, err := s.PackedTransactionAndCFD()
// 	if err != nil {
// 		return
// 	}

// 	for _, sig := range s.Signatures {
// 		pubKey, err := sig.PublicKey(SigDigest(chainID, trx, cfd))
// 		if err != nil {
// 			return nil, err
// 		}

// 		out = append(out, pubKey)
// 	}

// 	return
// }

// func (s *SignedTransaction) PackedTransactionAndCFD() ([]byte, []byte, error) {
// 	rawtrx, err := MarshalBinary(s.Transaction)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	rawcfd := []byte{}
// 	if len(s.ContextFreeData) > 0 {
// 		rawcfd, err = MarshalBinary(s.ContextFreeData)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 	}

// 	return rawtrx, rawcfd, nil
// }

/*func (tx *Transaction) ID() string {
	return "ID here" //todo
}*/
func (tx *Transaction) ID() common.TransactionIdType {

	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		fmt.Println("Transaction ID() is error :", err.Error()) //TODO
	}

	return common.TransactionIdType(crypto.Hash256(b))
}

// func (s *SignedTransaction) Pack(compression CompressionType) (*PackedTransaction, error) {
// 	rawtrx, rawcfd, err := s.PackedTransactionAndCFD()
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch compression {
// 	case CompressionZlib:
// 		var trx bytes.Buffer
// 		var cfd bytes.Buffer

// 		// Compress Trx
// 		writer, _ := zlib.NewWriterLevel(&trx, flate.BestCompression) // can only fail if invalid `level`..
// 		writer.Write(rawtrx)                                          // ignore error, could only bust memory
// 		err = writer.Close()
// 		if err != nil {
// 			return nil, fmt.Errorf("tx writer close %s", err)
// 		}
// 		rawtrx = trx.Bytes()

// 		// Compress ContextFreeData
// 		writer, _ = zlib.NewWriterLevel(&cfd, flate.BestCompression) // can only fail if invalid `level`..
// 		writer.Write(rawcfd)                                         // ignore errors, memory errors only
// 		err = writer.Close()
// 		if err != nil {
// 			return nil, fmt.Errorf("cfd writer close %s", err)
// 		}
// 		rawcfd = cfd.Bytes()

// 	}

// 	packed := &PackedTransaction{
// 		Signatures:            s.Signatures,
// 		Compression:           compression,
// 		PackedContextFreeData: rawcfd,
// 		PackedTransaction:     rawtrx,
// 	}

// 	return packed, nil
// }

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature        `json:"signatures"`
	Compression           common.CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData common.HexBytes        `json:"packed_context_free_data"`
	PackedTransaction     common.HexBytes        `json:"packed_trx"`
}

func (p *PackedTransaction) ID() (id common.TransactionIdType) {
	return //TODO
}

func (p *PackedTransaction) Expiration() common.TimePointSec {
	return common.TimePointSec(0) //TODO
}

func (p *PackedTransaction) GetUnprunableSize() uint32 {
	size := common.DefaultConfig.FixedNetOverheadOfPackedTrx
	size += uint32(len(p.PackedTransaction))
	max := ^uint(0) >> 1 / 2
	if size >= uint32(max) {
		log.Error("packed_transaction is too big")
		return 0
	}
	return size
}

func (p *PackedTransaction) GetPrunableSize() uint32 {
	size, _ := rlp.EncodeSize(p.Signatures)
	size += len(p.PackedContextFreeData)
	max := ^uint(0) >> 1 / 2
	if uint32(size) >= uint32(max) {
		log.Error("packed_transaction is too big")
		return 0
	}
	return uint32(size)
}

func GetTransaction() Transaction {
	//LocalUnpack()
	/*if (!unpacked_trx) {
		try {
			switch(compression) {
		case none:
			unpacked_trx = unpack_transaction(packed_trx);
			break;
		case zlib:
			unpacked_trx = zlib_decompress_transaction(packed_trx);
			break;
		default:
			EOS_THROW(unknown_transaction_compression, "Unknown transaction compression algorithm");
		}
		} FC_CAPTURE_AND_RETHROW((compression)(packed_trx))
	}*/
	t := Transaction{}
	return t
}

// // Unpack decodes the bytestream of the transaction, and attempts to
// // decode the registered actions.
// func (p *PackedTransaction) Unpack() (signedTx *SignedTransaction, err error) {
// 	return p.unpack(false)
// }

// // UnpackBare decodes the transcation payload, but doesn't decode the
// // nested action data structure.  See also `Unpack`.
// func (p *PackedTransaction) UnpackBare() (signedTx *SignedTransaction, err error) {
// 	return p.unpack(true)
// }

// func (p *PackedTransaction) unpack(bare bool) (signedTx *SignedTransaction, err error) {
// 	var txReader io.Reader
// 	txReader = bytes.NewBuffer(p.PackedTransaction)

// 	var freeDataReader io.Reader
// 	freeDataReader = bytes.NewBuffer(p.PackedContextFreeData)

// 	switch p.Compression {
// 	case CompressionZlib:
// 		txReader, err = zlib.NewReader(txReader)
// 		if err != nil {
// 			return nil, fmt.Errorf("new reader for tx, %s", err)
// 		}

// 		if len(p.PackedContextFreeData) > 0 {
// 			freeDataReader, err = zlib.NewReader(freeDataReader)
// 			if err != nil {
// 				return nil, fmt.Errorf("new reader for free data, %s", err)
// 			}
// 		}
// 	}

// 	data, err := ioutil.ReadAll(txReader)
// 	if err != nil {
// 		return nil, fmt.Errorf("unpack read all, %s", err)
// 	}
// 	decoder := NewDecoder(data)
// 	// decoder.DecodeActions(!bare)

// 	var tx Transaction
// 	err = decoder.Decode(&tx)
// 	if err != nil {
// 		return nil, fmt.Errorf("unpacking Transaction, %s", err)
// 	}

// 	signedTx = NewSignedTransaction(&tx)
// 	//signedTx.ContextFreeData = contextFreeData
// 	signedTx.Signatures = p.Signatures
// 	signedTx.packed = p

// 	return
// }

type DeferredTransaction struct {
	*Transaction

	SenderID   uint32             `json:"sender_id"`
	Sender     common.AccountName `json:"sender"`
	DelayUntil common.JSONTime    `json:"delay_until"`
}

// TxOptions represents options you want to pass to the transaction
// you're sending.
type TxOptions struct {
	ChainID          common.ChainIdType // If specified, we won't hit the API to fetch it
	HeadBlockID      common.BlockIdType // If provided, don't hit API to fetch it.  This allows offline transaction signing.
	MaxNetUsageWords uint32
	DelaySecs        uint32
	MaxCpuUsageMS    uint8 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress common.CompressionType
}

// FillFromChain will load ChainID (for signing transactions) and
// HeadBlockID (to fill transaction with TaPoS data).
// func (opts *TxOptions) FillFromChain(api *API) error {
// 	if opts == nil {
// 		return errors.New("TxOptions should not be nil, send an object")
// 	}

// 	if opts.HeadBlockID == nil || opts.ChainID == nil {
// 		info, err := api.cachedGetInfo()
// 		if err != nil {
// 			return err
// 		}

// 		if opts.HeadBlockID == nil {
// 			opts.HeadBlockID = info.HeadBlockID
// 		}
// 		if opts.ChainID == nil {
// 			opts.ChainID = info.ChainID
// 		}
// 	}

// 	return nil
// }
