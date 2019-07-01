package types

import (
	"encoding/json"
	"testing"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	data := []byte{0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0}
	packedTrx := PackedTransaction{}
	err := rlp.DecodeBytes(data, &packedTrx)
	assert.NoError(t, err)
	id := packedTrx.ID()
	assert.Equal(t, "e97f9f1e4aaafe1b92feded9bdd140247465de773154bcccab86986e1806fa33", id.String())
	//trx := packedTrx.GetTransaction()
	//re, _ := json.Marshal(trx)
	//fmt.Println("Trx:  ", string(re))
	result, err := rlp.EncodeToBytes(packedTrx)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestSignedBlock(t *testing.T) {
	data := []byte{0x66, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0xec, 0x3b, 0x3f, 0x3a, 0xe, 0x67, 0x96, 0xc6, 0xc5, 0x0, 0x43, 0xd1, 0x47, 0xac, 0xe2, 0x31, 0x93, 0xa6, 0x6e, 0x4b, 0x88, 0x55, 0x7b, 0x81, 0x50, 0x93, 0xa5, 0xf6, 0x3a, 0x1c, 0x77, 0x50, 0xf6, 0x33, 0x7f, 0x9e, 0x46, 0x91, 0xa4, 0xca, 0x2e, 0x32, 0x55, 0x48, 0x7, 0x2c, 0xc2, 0x82, 0x7a, 0xae, 0x7f, 0xba, 0x5f, 0xaa, 0x17, 0xb0, 0x38, 0xd5, 0xf9, 0xb, 0x44, 0x48, 0xb9, 0x47, 0x64, 0x11, 0x34, 0x75, 0x7b, 0xc0, 0x15, 0x27, 0xeb, 0x52, 0x14, 0x3b, 0x4d, 0x61, 0xf9, 0xd6, 0x49, 0x24, 0xf2, 0x4b, 0x7f, 0x19, 0x20, 0x7c, 0x2f, 0x46, 0xc2, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x42, 0x22, 0x7c, 0xd4, 0x29, 0x7, 0x34, 0x3d, 0x90, 0x3d, 0xd2, 0x6a, 0x10, 0xbb, 0x41, 0x93, 0x57, 0x4c, 0xf, 0xad, 0xec, 0x90, 0x27, 0x66, 0xa9, 0xe5, 0x4f, 0x4c, 0xde, 0xfd, 0x3, 0xc3, 0x25, 0x3f, 0x7d, 0x66, 0x77, 0x1e, 0x14, 0xf3, 0x5f, 0x9c, 0xd9, 0xc, 0xcf, 0xe9, 0x6a, 0x5b, 0x3d, 0xfa, 0x80, 0x8f, 0xf, 0x6c, 0xea, 0xf7, 0x9b, 0xdf, 0x2f, 0x74, 0xab, 0x6f, 0x47, 0x9e, 0x1, 0x0, 0x2d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	signedBlock := SignedBlock{}
	err := rlp.DecodeBytes(data, &signedBlock)
	assert.NoError(t, err)
	//data2, err := json.Marshal(signedBlock)
	//assert.NoError(t,nil,err)
	//fmt.Println("Receive P2PMessag ", string(data2))

	result, err := rlp.EncodeToBytes(signedBlock)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestTransaction_GetSignatureKeys(t *testing.T) {
	data := []byte{0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0}
	packedTrx := PackedTransaction{}
	err := rlp.DecodeBytes(data, &packedTrx)
	assert.NoError(t, err)
	trx := packedTrx.GetTransaction()

	chainID := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	set := trx.GetSignatureKeys(packedTrx.Signatures, &chainID, []common.HexBytes{}, false, true)
	digist := trx.SigDigest(&chainID, nil)
	bytes, err := rlp.EncodeToBytes(digist)
	assert.Equal(t, digist.Bytes(), bytes)
	pub, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	keys := set.Find(func(value ecc.PublicKey) bool {
		return set.Comparator(value, pub) == 0
	})
	assert.Equal(t, pub, keys)
}

func TestTransactionID(t *testing.T) {
	data := []byte{0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	signedBlock := SignedBlock{}
	err := rlp.DecodeBytes(data, &signedBlock)
	assert.NoError(t, err)
	//data, err := json.Marshal(signedBlock)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Receive P2PMessag ", string(data))

	transactions := signedBlock.Transactions
	for _, TrxReceipt := range transactions {

		if TrxReceipt.Trx.TransactionID == common.TransactionIdType(crypto.NewSha256Nil()) {
			packedTrx := TrxReceipt.Trx.PackedTransaction

			//enc, _ := rlp.EncodeToBytes(packedTrx)
			//fmt.Printf("%#v\n", enc)
			data, err := json.Marshal(packedTrx.PackedTrx)
			assert.NoError(t, err)
			assert.Equal(t, "\"f035c95b3700380ad3d100000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000\"",
				string(data))
			trx := packedTrx.GetTransaction()
			assert.Equal(t, "9efe81d63d3ba19d6f4a7da3f89ac3315532ec8a6ebfaeb1c52dda63a620f322", trx.ID().String())

			signedTrx := packedTrx.GetSignedTransaction()
			//fmt.Println(signedTrx.Transaction)
			//data, err = json.Marshal(signedTrx.Transaction)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//fmt.Println("trx signed  ", string(data))

			newPackedTrx := NewPackedTransactionBySignedTrx(signedTrx, CompressionNone)
			assert.Equal(t, packedTrx.PackedTrx, newPackedTrx.PackedTrx)
		}
	}
}

//func TestReceiveSignedBlock(t *testing.T) {
//	data1 := []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0x66, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0xec, 0x3b, 0x3f, 0x3a, 0xe, 0x67, 0x96, 0xc6, 0xc5, 0x0, 0x43, 0xd1, 0x47, 0xac, 0xe2, 0x31, 0x93, 0xa6, 0x6e, 0x4b, 0x88, 0x55, 0x7b, 0x81, 0x50, 0x93, 0xa5, 0xf6, 0x3a, 0x1c, 0x77, 0x50, 0xf6, 0x33, 0x7f, 0x9e, 0x46, 0x91, 0xa4, 0xca, 0x2e, 0x32, 0x55, 0x48, 0x7, 0x2c, 0xc2, 0x82, 0x7a, 0xae, 0x7f, 0xba, 0x5f, 0xaa, 0x17, 0xb0, 0x38, 0xd5, 0xf9, 0xb, 0x44, 0x48, 0xb9, 0x47, 0x64, 0x11, 0x34, 0x75, 0x7b, 0xc0, 0x15, 0x27, 0xeb, 0x52, 0x14, 0x3b, 0x4d, 0x61, 0xf9, 0xd6, 0x49, 0x24, 0xf2, 0x4b, 0x7f, 0x19, 0x20, 0x7c, 0x2f, 0x46, 0xc2, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x42, 0x22, 0x7c, 0xd4, 0x29, 0x7, 0x34, 0x3d, 0x90, 0x3d, 0xd2, 0x6a, 0x10, 0xbb, 0x41, 0x93, 0x57, 0x4c, 0xf, 0xad, 0xec, 0x90, 0x27, 0x66, 0xa9, 0xe5, 0x4f, 0x4c, 0xde, 0xfd, 0x3, 0xc3, 0x25, 0x3f, 0x7d, 0x66, 0x77, 0x1e, 0x14, 0xf3, 0x5f, 0x9c, 0xd9, 0xc, 0xcf, 0xe9, 0x6a, 0x5b, 0x3d, 0xfa, 0x80, 0x8f, 0xf, 0x6c, 0xea, 0xf7, 0x9b, 0xdf, 0x2f, 0x74, 0xab, 0x6f, 0x47, 0x9e, 0x1, 0x0, 0x2d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//	data2 := []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0x76, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x15, 0xfc, 0xc0, 0xb8, 0x41, 0x5, 0x50, 0x37, 0xe7, 0xa7, 0xf2, 0xe3, 0x99, 0xed, 0xad, 0x14, 0x9a, 0x58, 0xa1, 0xb, 0x38, 0xfc, 0x75, 0xae, 0x7c, 0x19, 0xc7, 0xcc, 0xd7, 0x9, 0x5a, 0x6c, 0x4e, 0x90, 0x6d, 0xe1, 0x2f, 0x21, 0x87, 0x3, 0x7, 0x45, 0x69, 0xfc, 0xb6, 0xb3, 0xe4, 0xc1, 0x89, 0x64, 0x3d, 0x5e, 0xce, 0x6c, 0x13, 0x41, 0xb5, 0x5b, 0xe8, 0x81, 0xdd, 0x9f, 0xf3, 0x7e, 0xdc, 0xf6, 0x77, 0x79, 0x16, 0x1b, 0xe4, 0xc, 0x78, 0x93, 0x72, 0xb1, 0x63, 0xac, 0xe6, 0x4, 0xb7, 0x9d, 0x30, 0xa4, 0x27, 0xff, 0x91, 0x1c, 0x46, 0x63, 0xa, 0xe9, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0xd, 0x90, 0x27, 0xae, 0xb, 0x4f, 0x8, 0x2e, 0x72, 0x26, 0x49, 0xe2, 0xb0, 0xcd, 0xc4, 0x64, 0x1e, 0x3d, 0xb0, 0xd9, 0x63, 0x7, 0xce, 0x36, 0x49, 0x71, 0xf6, 0xcb, 0x35, 0x3b, 0xe0, 0xe6, 0x31, 0x48, 0x9f, 0x46, 0x7c, 0xaa, 0x12, 0x7b, 0xf2, 0x22, 0x72, 0xf9, 0x4f, 0x9a, 0x2b, 0x8c, 0x14, 0xeb, 0x99, 0xb0, 0x35, 0x26, 0xb0, 0x23, 0x4e, 0x99, 0xf5, 0xa8, 0xc3, 0x67, 0x3e, 0xa9, 0x1, 0x0, 0xd2, 0x0, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x37, 0x3, 0x2f, 0x8f, 0x54, 0x1c, 0xd3, 0x36, 0x82, 0x9a, 0x88, 0xa7, 0x63, 0x20, 0xe7, 0xb3, 0xea, 0xee, 0xa7, 0x50, 0x89, 0x26, 0x19, 0x9f, 0x91, 0x30, 0x3f, 0x1d, 0x5b, 0xac, 0xc1, 0xf6, 0x11, 0x7e, 0x7d, 0x9e, 0x48, 0x95, 0xba, 0xdf, 0xe3, 0x8, 0x99, 0x1e, 0xcc, 0xf1, 0xb3, 0xfa, 0xb9, 0x32, 0xd7, 0x78, 0x52, 0x73, 0x11, 0x6d, 0xea, 0x79, 0x7c, 0x8, 0x68, 0x2e, 0x79, 0xbe, 0x0, 0x0, 0x98, 0x1, 0x58, 0xeb, 0xc3, 0x5b, 0x14, 0x0, 0x4c, 0xc1, 0x4b, 0x27, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1a, 0xa3, 0x6a, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//	fmt.Println(data1)
//	fmt.Println(data2)
//}

//receive data:  []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//signed Block Num: 57
//---------------*************----------------
//Receive signedBlock:    {"timestamp":"2018-10-19T01:39:31.000","producer":"eosio","confirmed":0,"previous":"000000381eef635b5c02e6dfb1aba578c32360b3517caed7d047bb864c110ecd","transaction_mroot":"c98fe67f358897605862d4e9a41363282f5ce436de079a8bcd9001a15ac3864d","action_mroot":"8547a5327f04cafc37437c2c1fde6c3b7a6e1f4614d22165a949e416cd65b8d9","schedule_version":0,"new_producers":null,"header_extensions":[],"producer_signature":"SIG_K1_K25NdAjgDn4niYaJZFTD4dESHqLMVCEdqqazQhAyCkmre1LtW6VfZ2tT5vKFutjbfXC8hYba2iTZcoz1M5k4agrSDZSUYh","transactions":[{"status":"executed","cpu_usage_us":285,"net_usage_words":25,"trx":[{"signatures":["SIG_K1_KdivyH4YyTLkPcqW7GDCovLkAixmMuUgWHnDjySkuA3iRPFXS3hr3BMUtcCVamRgrgYVSDF9RXjGa6Dfd6V9SeRa7tM4Pi"],"compression":"none","packed_context_free_data":"","packed_trx":"f035c95b3700380ad3d100000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000","UnpackedTrx":null},"0000000000000000000000000000000000000000000000000000000000000000"]}],"block_extensions":[]}
//encode result: []byte{0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//[]byte{0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0}
//encode result : packedTrx:    "f035c95b3700380ad3d100000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000"
//trx receive   {"expiration":1539913200,"ref_block_num":55,"ref_block_prefix":3520268856,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000000"}],"transaction_extensions":[]}
//cbb52b1177b73e47f0282d8fb706158acadf78f872f0bb93e94acad159dec50a
//---------------*************----------------
//receive data:  []byte{0xb9, 0x0, 0x0, 0x0, 0x7, 0xa7, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x39, 0x98, 0x9c, 0x66, 0x39, 0x71, 0x71, 0x73, 0xfb, 0x68, 0xf7, 0xec, 0xca, 0x7a, 0x24, 0xb3, 0xd, 0xc4, 0xca, 0x81, 0x17, 0x2a, 0x37, 0x39, 0x51, 0x2c, 0x26, 0xa2, 0xa1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xad, 0x5e, 0x40, 0x20, 0x86, 0x89, 0xc2, 0xae, 0x4e, 0x96, 0x9d, 0x33, 0x7d, 0xdd, 0x49, 0xac, 0xdb, 0xe6, 0xe7, 0xaa, 0xa0, 0xf0, 0x8, 0x2d, 0xb4, 0xfd, 0x80, 0x1c, 0x4a, 0xaf, 0x45, 0xa2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x7f, 0x48, 0x9d, 0xf0, 0x2a, 0x33, 0xf5, 0x2b, 0xdc, 0xe1, 0x35, 0x66, 0xee, 0x4, 0xb7, 0xf8, 0xb1, 0x3d, 0x4e, 0x39, 0x53, 0xd2, 0xa2, 0xed, 0xd5, 0x29, 0x2e, 0x5, 0x3a, 0xb3, 0x13, 0xfb, 0xf, 0x83, 0x8, 0xa4, 0x28, 0xce, 0xe3, 0x29, 0xcd, 0x7c, 0x91, 0xf7, 0x4a, 0x35, 0x50, 0x6e, 0x89, 0x6, 0x5c, 0x5a, 0x84, 0xe4, 0x9b, 0x90, 0xf7, 0x99, 0xd9, 0xc, 0xb8, 0xf2, 0x5f, 0xd6, 0x0, 0x0}
//signed Block Num: 58

// signed Block Num: 127
// Receive P2PMessag
// {
// 	"timestamp": "2018-06-09T11:57:33.000",
// 	"producer": "eosio",
// 	"confirmed": 0,
// 	"Previous": "0000007fc12a85c48a04d10f51243ff3e4fc91cbb00a8107fcbf5312af22ec04",
// 	"transaction_mroot": "645c4333b805c768f2ed701402754bf85ce5b7b82ab6e642bc54109809298a01",
// 	"action_mroot": "3dc8be647fbaf84a3a6eb5abd1cfa4f3e99b6704c09c2dd165afe5e702c8f8ca",
// 	"schedule_version": 0,
// 	"new_producers": null,
// 	"header_extensions": [],
// 	"ProducerSignature": "SIG_K1_KgrhinBzxUnY4sqvs3pFsxPouido11Ky3WAyFbJzjDDcA9R54bD17eyagg9M7gtA2FQ2MC4ZCmUzJRbfm6VxUNYKmmCXZP",
// 	"transactions": [{
// 		"status": "executed",
// 		"cpu_usage_us": 1460,
// 		"net_usage_words": 465,
// 		"trx": [{
// 			"signatures": ["SIG_K1_KeuZmWqgvzfBQFxE3cDAebvxRmrCCaarDi1rH6oi4aGUFXP9sKu3Abm8owbWK7BjnCcGL8K5pJYErQeth8GZEZgdYteafu"],
// 			"compression": "zlib",
// 			"packed_context_free_data": "",
// 			"packed_trx": "78dacd595d6c1cd775bef7ceefeeec90434b7664d2a9ef4e1c87096a97521d8a705a444397b214a5962db74e10c158ad76c7d2fe93bb4bca8c05ad02abb1933ce6a18003a445e0c44e60c3a80bb88253a08c60a0468304491019895f6af4a14002b430fae0164580f43be7ce2e871425a079ca48dc9939f7de73bf7bce77ceb933f32f5b73a72f89179e7d20113814fd885f2ffc257e8f7cf4abaf5f93398178f93f0e1dfaaffbb725ff7ebfa80e3ad4459ebded8c35ba3412679cd1251ce28c3de2931a8de8075702ffe5a5338aa4428ea805bfd6884fd27445bb6d862bd225f167f1ad4427f59fb264a5dd0da77ab6d71f0a87aea36a6dd8e8752bf5eab05a1934be980a9fc4a5da7abf9f76879561a3930a9745696fd0e855aa83418ac11689dc4edaa9ad6e8ae2f866900ecd4dd44fabf54a4eb728b08e7ebab6dee8a795eafaf0bc084814e645878447b229e8a9acf61b1b8d767a2ead8b2912de6184bdd5de20add3457dbd96f6078297b48f1afbe9a0b7deafa59576a3d3180e44687df8ae42a1e04f47a2200baa6015ec42c1b266a6a78b33859969c776e4aa745dc792427aafc922ada0d7df54e28eca17d2b553270e1eac9d4f6badc17ae7d027171f3f5811a51b1a0e2e2e504388866eba47c394d5ed5d10d30b952f3cf249b6dfc143f9f5ae9c3af1c8e395838bab69bfd3180cc858ed74236daf88c8a9aeaeb637c50cdbb8b32a62b7536db77b35f111fba97e9a8a7b0b25291211b81fbb6dff1d07668bff347f5720b4d032d151bc12e4afc3c975221e0e8a2252cb5b57fef4f24bc1141a3e6ef18f1ff9c107aea546f625394aa0785ed9c97ba215171695bd2ccaeef2a7ca4e1295e9dab38ea8c93fed2ebb4f6006fb0f8588ade46f46cde4b7722db9f34428b5957c5b36cb56490518e5a367f0f95095ac000d3f91cd44e07cd274778ecf9549feeaf2d29797979e2bfb81f697ef7e563bcbbfc531fdec97496227d066438dbc5876b5afbd2b650f3d9eb9183bcb5f3c158a00683cad703501acf7046cffbe00065cf958b210590123935a3c9603ff7b6e6d636b19e8821687adcf4f202a68f9870f5efcebd77fb9fc394055cb972f5f35b78f0125dd5edefafaffbe43b7cadc5e7ee3bf717b2ab47521118bea2f70928bea51282e7cdc7af4b0f221f88568e2d76f46fb664b76908d7bf7956fbcf3e668ac87667981f45846cf6771528bea98d173ecb08a267aa26674fb2cd9863bfe194ed622a073c72387d58149c703cde843b330a9e9f82738d98b6ac9745c3aacf4a4a36e4677ce9664d6f1019c9c45b5603a2e1c56f3938ef3cd686e36803b0a88af26a22bb8cd846104326885108d0ac16fa4858c7dc91e99408c5ab1577631c6437f2d212a2b3a09b044946d752492b15c09856151f225fb380c2fa38f946d20a261104e3793d5b5566cd38c1a4a227736d0d0f9a020aec0839423bc13c9fb22b2b40bf7fb913d0b91ff702682f1638fe4cd44a22937931ddd1b20c5a844ae37cb527bf70941ccc19aca161122916b2b040880b5d364d801d66a4119e8a89d198184830922324722837d798350b728085e539c9caccc225b484dc5326cad8bbb4de28f9908a35834afc530311f0c22e14b98cc9c74116d995d24db853a915d0a6410df9c843945e6046a68f4bdc79c349db405dcf3e46ab020733a93245922872f34a33f8015ac895b76188b284686f2989574e532cfe9ca21f2b07f77dad0671bfa04d7d60e9ce5697f86aa01466e89cc8885b111833714c8645d7232d3e956ec978be8ecefb65c216739c9b36596232af9c672e6847c306194cf96f33346158dad8ac6804563b26d6215b134ed8fb9f5f45aec4e08464d86633e5f51177022df1a3563a4229cb835c7407fa7515d633d8fada78b286f3712b1c0462c644484190b338a4c58642792051f97ca982c56652b3366661d91cd2b601865624c6931b688628b40008bf0a9384be997a4d76c39f6c3422b76297eb5bbdb0f363174c25f0eea317b0c877942f766e4752803f9c8d54c5ec7d0d5317475267445716c6a27ba6b763737795c44108ca603da0136fa41cc59c942a7ec66999ea3cb4ddef9def54f9c0ca9f7812620dc1dc301c9421bb364b347e6b47f5b51ec922a798ce6a48590af1d6ad90f2c0b5d1a4d8bdabfedd4bcff2cf29fcdfeb333bcda609d873a6841be94045411456884ca402a749b6f629d0452ed00a927969a37482302a9b6412a0392c775cdd07b50448d21ed190b8ea261da844ea64c463a30f22827df9f938b9cdc66397ceb206b3009a7c99da8566511f04474e192a3bce0baa911d6768db0b946d837d4885b85b56dc2da9ca82edca4509878ce55080ae4ac48b8e300dbbb20dcb2043833f6cee4ff33655958961a695a56d94b683b44d12c306399a2cc061b2cdeea41cff1e4b2041ca2449c151e4b7bd8bbc41eb9161546ac413bea4cec0e5fd5ce15687213af59a612e76df0ae91c9ecbc141bb3c5c66cb1d7028562f7cfd9fd2867ade8c35903d6c9ec8bd5c32155128c46fb49b6acddd4542489c04803669774947649843bb669e52ec0d36262e6b1db9ad81403213136b5b36b643c5ac7b86b79dc95cd4f3de25cef38ebed813e80097822c0ccc13f2b649ddc2e1296f5b14cce6d2e2d534b42ef9244c40e2fdcd2b0a13d8e9be43ace14f8c867a65d25dfbdfebddb4f90b7c9bd43646698f4f85cecae842ae00c10dd5de60d28bb4c60d5a02776bbce2cb84133114ac62b9311e23d14d8d91119d18db32a3227f15e997402cf2aca0eb0acc336ccb4410b9c3956813031a6668bf17a2485f1a31dd6e118321e331167629a3a956936edb0a93563126c6146c7cc8df47410fcab948698d8d9183280833e11d1de76a3c5f7c68d2abb26729b71e2a6e378c69b0cbdf998e8d663248d91bbc6204fcbc918995df35e7409eb2382247816998b91206211c27370083f135e91b15aa15dfc8cc31bcc58ac64149ba2f4477770dabd41e0884888e0b8a4ec648cce958c987f1f3dba684927fb61ca42f4a021694a9004ce811f15a8a0024aae768b128583dd4cf23cfa44f7043fb2656994e57543e5ad240e4318a21c42c9d60f04f6091c7992ec491199786bb16a91bbe7a8248c43e0ad242eeaf0243f1e179352bb995cc6ee89f3559184c915dcc621acc8918d4662032661ed8a485db478c74de114ea225b559c3079ced621066b1a9ed0b3cad36b5af50dd911d4d2a4c262d3a8dceec55d4222b00c2974a28f32ab83e43778344348b7100622791bab2c07d972cb255acddb496c954332308697a87e26cf03ffbccab0236a707d3479ef1a7140f2e2b812d33a689ccbd5a8899d68985c6cc63e3516127a22f4d6d828d86d1f639bc390278915008e1c51d41e55ff22f52f5217ac0df634039bbc12ad8e334d7cded0d3ed1cb419abd0828d4ecab4a1564d2a3acd6cbc7124cd876244d391be507ba48fb72854b502b29e646beb122d722e0e8d5d2d9301a98cfadadf65e4698c0d9e73a43f326ca0452a50fd3ef1c6432017f1fdea43650f89e8d382ea50221f440b3d4e25d190cacbd58790dea9f0659cf200da02fb3662f56911fbb41ea5fdd611314b1b670d21f89ea9b578336db685571fc250b22dbfa781ef50e63c20b5f1e71bad6bb19f7ce97face3b13bd7d446e6a1aecdb5b4d72a7b63d4a100522b87d422a49641aa0cbd95990bb597d3b9875a48fb7007dbc2097c1ff06d82ef2300015ffb045d65d05db3fb749a8cdc4a46a86f63e4718e585c9f63c4a07534a49d186d1f9e1f13319647b9e248dad2bbdb41925cd4b285279cbed998910abacafae4482192b7c0fd98331e7b5f71d6c6320986095e322e0532d3dfa7a7670cf4a7a9601b6c13d8964166eafb517e68ccd6a13db472fddd46e84d10ca09c2ac4f1ee118a39721049aad04559e5025629a2a1b40710a50c133c895263365dc7b7b39db7d26cfff80ca0d2d8c2e91a12cdab2675b10a463ae69d87870ea518c836ff8dd4f89e04b8ebf69b2dac5c972b6a15ac8d5c12fa587c219d85bc7059ec7a3c0ef75d37ebfd7c78d0e5c7e7b89cb85e04826d7e6f5eec706bada4f75afdbded41bd576a3ae9feaf73a7a783ed5310f8af56073304c3be85febad778750f2be081c7a438ccbcb32b0cfa52444620a9e302f392ba4a242ef39ebfa4275a0ebe9a071ae8b9b61cf4c7436d535f484a4fad430edebcaf9b4baca5d6bbdce6a3b1da6e8c4aaaac3b42ef2efe0afd6f67829fffd702cf9bb708a513ff860f56ce38f0edebfe09432dc956eb593daf4339d7b5f3b1106d9dbee8960dfb05fed0e3269a35e196eaea641ee8d7270216d9c3b3f64b9bbdee80e0f2eee8f76bf0916ca81865e7f0788e276afdd508aad74b362140b65e1a6b8ba7eb6dda85570e91a797ede0fed9e6f323637c70da0f650145ca8368693c13edf0dd21aafeb8f0fed31a040efbf7bfdc67053d885e1f97e3a38df6bd7b3fe36d00e4adb4b39fda49fad7f70e74d109f7ed2a14907610ec8e9278bddf4c29876b65703e776db925db543e2f42e80e0dbf85c72e1460eb0374887b55e3d85c66cdc8ef1ee468716e8d052960a1b9d8db44f68cdbd4d039db39bc374e0420d3826d49e5a2c34997ec5f5d53a58cc9f4cf69ef1167c7057abf41de706c692b6ed1515eb29c50c4fa1febf53f8ed46b7750b74bce49d12664c2e5e82eca348670fa8c5f5ee6402eb779ea056edd6d2369659dd146acadc35bae7f82bcc8dfc1ef69f46c4ee15c0de3829aac220edd6d33e1a3c0ede434bfe80bf98f59f367ebb8dbe498d51659fa4f6b650a15fed54788c034d8b0f80b5630a6782daeafa0e0151903e90ddc45b5e63c0dfcf6ca2dc0ce138d7ee9dadb6c730e43e52b88e00ad60edf4d7e8d58de6d2f8cb1a650ca1c2c9ed0d71327316e9b555a1cc4c96dc996bc2ec6b5df6914efa0364befa7a3b9dcaeb3ffde4ce2f82d2a6c4bf63962971e49b2fc4f5ab2fe44379fc71751289e3c43e8e2971e487efbedc3e753d1f3b2c7be55b9fc9931dc77d2deba5af6d73987a6dfdeadf7e9ea79df8c7af7c62e5afbe93ec60118ef77ef6ccf517278c10e2c78fbc79fef56b7bf99ddbde7efdda1ebec071e6fba7fffedac4a538967e0ac12e13e2f8dbb4f3ca9b3b2dc6c7ff01fba693c9",
// 			"UnpackedTrx": null
// 		}, "0000000000000000000000000000000000000000000000000000000000000000"]
// 	}],
// 	"block_extensions": []
// }
