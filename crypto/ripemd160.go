package crypto

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception/try"
	"golang.org/x/crypto/ripemd160"
	"hash"
)

type Ripemd160 struct {
	Hash [5]uint32 `eos:"array"`
}

func NewRipemd160() hash.Hash {
	return ripemd160.New()
}
func NewRipemd160Nil() *Ripemd160 {
	data := [5]uint32{0, 0, 0, 0, 0}
	return &Ripemd160{
		Hash: data,
	}
}
func NewRipemd160String(s string) *Ripemd160 {
	bytes, err := hex.DecodeString(s)
	Throw(err)

	result := new(Ripemd160)
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
	}

	return result
}

func NewRipemd160Byte(s []byte) *Ripemd160 {
	result := new(Ripemd160)
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(s[i*4 : (i+1)*4])
	}
	return result
}

func (h Ripemd160) Bytes() []byte {
	result := make([]byte, 20)
	for i := range h.Hash {
		binary.LittleEndian.PutUint32(result[i*4:(i+1)*4], h.Hash[i])
	}
	return result
}

func (h Ripemd160) String() string {
	return hex.EncodeToString(h.Bytes())
}

func HashRipemd160(t interface{}) Ripemd160 {
	cereal, err := rlp.EncodeToBytes(t)
	Throw(err)
	h := ripemd160.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	result := Ripemd160{}
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(hashed[i*4 : (i+1)*4])
	}

	return result
}

func (h Ripemd160) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h.Bytes()))
}
