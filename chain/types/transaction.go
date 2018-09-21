package types

import (
	"encoding/json"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/log"
	"time"
	"github.com/eosspark/eos-go/rlp"
)

type TransactionHeader struct {
	Expiration     common.JSONTime `json:"expiration"`
	RefBlockNum    uint16              `json:"ref_block_num"`
	RefBlockPrefix uint32              `json:"ref_block_prefix"`

	MaxNetUsageWords uint32 `json:"max_net_usage_words"`
	MaxCPUUsageMS    uint8  `json:"max_cpu_usage_ms"`
	DelaySec         uint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration
}

func (th TransactionHeader) GetRefBlocknum(headBlocknum uint32) uint32 {
	return headBlocknum/0xffff*0xffff + headBlocknum%0xffff
}

func (th TransactionHeader) VerifyReferenceBlock(referenceBlock common.BlockIDType) bool {
	return th.RefBlockNum == uint16(common.EndianReverseU32(uint32(referenceBlock.Hash_[0]))) &&
		th.RefBlockPrefix == uint32(referenceBlock.Hash_[1])
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
	tx.Fill(opts.HeadBlockID, opts.DelaySecs, opts.MaxNetUsageWords, opts.MaxCPUUsageMS)
	return tx
}

func (tx *Transaction) SetExpiration(in time.Duration) {
	tx.Expiration = common.JSONTime{time.Now().UTC().Add(in)}
}

func (tx *Transaction) GetSignatureKeys(chainId common.ChainIDType,allowDeplicateKeys bool,useCache bool) []common.PublicKeyType{
	//TODO
	return nil
}

type Extension struct {
	Type uint16          `json:"type"`
	Data common.HexBytes `json:"data"`
}

// Fill sets the fields on a transaction.  If you pass `headBlockID`, then `api` can be nil. If you don't pass `headBlockID`, then the `api` is going to be called to fetch
func (tx *Transaction) Fill(headBlockID common.BlockIDType, delaySecs, maxNetUsageWords uint32, maxCPUUsageMS uint8) {
	tx.setRefBlock(headBlockID)

	if tx.ContextFreeActions == nil {
		tx.ContextFreeActions = make([]*Action, 0, 0)
	}
	if tx.Extensions == nil {
		tx.Extensions = make([]*Extension, 0, 0)
	}

	tx.MaxNetUsageWords = uint32(maxNetUsageWords)
	tx.MaxCPUUsageMS = maxCPUUsageMS
	tx.DelaySec = uint32(delaySecs)

	tx.SetExpiration(30 * time.Second)
}

func (tx *Transaction) setRefBlock(blockID common.BlockIDType) {
	tx.RefBlockNum = uint16(blockID.Hash_[0])
	tx.RefBlockPrefix = uint32(blockID.Hash_[1])
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

func (st *SignedTransaction) GetSignatureKeys(chainId common.ChainIDType,allowDeplicateKeys bool,useCache bool) []common.PublicKeyType{
	//TODO

	return st.Transaction.GetSignatureKeys(chainId,allowDeplicateKeys,useCache)
}

func (head *TransactionHeader) SetReferenceBlock(referenceBlock common.BlockIDType) {
	first := common.EndianReverseU32(uint32(referenceBlock.Hash_[0]))
	head.RefBlockNum = uint16(first)
	head.RefBlockPrefix = uint32(referenceBlock.Hash_[1])
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

func (tx *Transaction) ID() string {
	return "ID here" //todo
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

func (p *PackedTransaction) ID() (id common.TransactionIDType) {
	return //TODO
}

func (p *PackedTransaction) Expiration() common.TimePointSec {
	return common.TimePointSec(0) //TODO
}

func (p *PackedTransaction) GetUnprunableSize() uint32{
	size := common.DefaultConfig.FixedNetOverheadOfPackedTrx
	size += uint32(len(p.PackedTransaction))
	max := ^uint(0)>>1 / 2
	if size >= uint32(max){
		log.Error("packed_transaction is too big")
		return 0
	}
	return size
}

func (p *PackedTransaction) GetPrunableSize() uint32{
	size,_ := rlp.EncodeSize(p.Signatures)
	size += len(p.PackedContextFreeData)
	max := ^uint(0)>>1 / 2
	if uint32(size) >= uint32(max){
		log.Error("packed_transaction is too big")
		return 0
	}
	return uint32(size)
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
	ChainID          common.ChainIDType // If specified, we won't hit the API to fetch it
	HeadBlockID      common.BlockIDType // If provided, don't hit API to fetch it.  This allows offline transaction signing.
	MaxNetUsageWords uint32
	DelaySecs        uint32
	MaxCPUUsageMS    uint8 // If you want to override the CPU usage (in counts of 1024)
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
