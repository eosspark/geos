package rlp

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"hash"
)

type Sha512 struct {
	Hash_ [8]uint64 `eos:"array"`
}

func NewSha512() hash.Hash {
	return sha512.New()
}

func NewSha512Nil() *Sha512{
	data := [8]uint64{0,0,0,0,0,0,0,0}
	return &Sha512{
		Hash_:data,
	}
}

func NewSha512String(s string) *Sha512 {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	result := new(Sha512)
	for i := range result.Hash_ {
		result.Hash_[i] = binary.LittleEndian.Uint64(bytes[i*8 : (i+1)*8])
	}

	return result
}

func NewSha512Byte(bytes []byte) *Sha512 {
	result := new(Sha512)
	for i := range result.Hash_ {
		result.Hash_[i] = binary.LittleEndian.Uint64(bytes[i*8 : (i+1)*8])
	}

	return result
}

func Hash512(t interface{}) Sha512 {
	cereal, err := EncodeToBytes(t)
	if err != nil {
		panic(err)
	}
	h := sha512.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	result := Sha512{}
	for i := range result.Hash_ {
		result.Hash_[i] = binary.LittleEndian.Uint64(hashed[i*8 : (i+1)*8])
	}

	return result
}

func (h Sha512) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h.Bytes()))
}

func (h Sha512) String() string {
	return hex.EncodeToString(h.Bytes())
}

func (h Sha512) Bytes() []byte {
	result := make([]byte, 64)
	for i := range h.Hash_ {
		binary.LittleEndian.PutUint64(result[i*8:(i+1)*8], h.Hash_[i])
	}
	return result
}

func (h Sha512) Or(h1 Sha512) Sha512 {
	result := Sha512{}
	for i := range result.Hash_ {
		result.Hash_[i] = h.Hash_[i] ^ h1.Hash_[i]
	}
	return result
}
