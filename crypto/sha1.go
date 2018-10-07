package crypto

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/eosspark/eos-go/crypto/rlp"
	"hash"
)

type Sha1 struct {
	Hash [5]uint32 `eos:"array"`
}

func NewSha1() hash.Hash {
	return sha1.New()
}
func NewSha1Nil() *Sha1 {
	data := [5]uint32{0, 0, 0, 0, 0}
	return &Sha1{
		Hash: data,
	}
}

func NewSha1String(s string) *Sha1 {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	result := new(Sha1)
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
	}

	return result
}

func NewSha1Byte(s []byte) *Sha1 {
	result := new(Sha1)
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(s[i*4 : (i+1)*4])
	}
	return result
}

func (h Sha1) Bytes() []byte {
	result := make([]byte, 20)
	for i := range h.Hash {
		binary.LittleEndian.PutUint32(result[i*4:(i+1)*4], h.Hash[i])
	}
	return result
}

func (h Sha1) String() string {
	return hex.EncodeToString(h.Bytes())
}

func Hash1(t interface{}) Sha1 {
	cereal, err := rlp.EncodeToBytes(t)
	if err != nil {
		panic(err)
	}
	h := sha1.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	result := Sha1{}
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint32(hashed[i*4 : (i+1)*4])
	}

	return result
}

func (h Sha1) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h.Bytes()))
}
