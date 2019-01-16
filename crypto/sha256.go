package crypto

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception/try"
	"hash"
)

type Sha256 struct {
	Hash [4]uint64 `eos:"array"`
}

func NewSha256() hash.Hash {
	return sha256.New()
}

func NewSha256Nil() *Sha256 {
	data := [4]uint64{0, 0, 0, 0}
	return &Sha256{
		Hash: data,
	}
}

func NewSha256String(s string) *Sha256 {
	bytes, err := hex.DecodeString(s)
	Throw(err)

	result := new(Sha256)
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint64(bytes[i*8 : (i+1)*8])
	}

	return result
}

func NewSha256Byte(s []byte) *Sha256 {
	result := new(Sha256)
	//if len(s) <32{
	//	return nil,errors.New("the length of slice is less then 32")
	//}

	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint64(s[i*8 : (i+1)*8])
	}
	return result
}

func Hash256(t interface{}) *Sha256 {
	cereal, err := rlp.EncodeToBytes(t)
	Throw(err)
	h := sha256.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	result := &Sha256{}
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint64(hashed[i*8 : (i+1)*8])
	}

	return result
}

func Hash256String(s string) *Sha256 {
	h := sha256.New()
	_, _ = h.Write([]byte(s))
	hashed := h.Sum(nil)

	result := &Sha256{}
	for i := range result.Hash {
		result.Hash[i] = binary.LittleEndian.Uint64(hashed[i*8 : (i+1)*8])
	}

	return result
}

func (h Sha256) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h.Bytes()))
}

func (h *Sha256) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	for i := range h.Hash {
		h.Hash[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return nil
}

func (h Sha256) String() string {
	return hex.EncodeToString(h.Bytes())
}

func (h Sha256) Bytes() []byte {
	result := make([]byte, 32)
	for i := range h.Hash {
		binary.LittleEndian.PutUint64(result[i*8:(i+1)*8], h.Hash[i])
	}
	return result
}

func (h Sha256) BigEndianBytes() []byte {
	result := make([]byte, 32)
	for i := range h.Hash {
		binary.BigEndian.PutUint64(result[i*8:(i+1)*8], h.Hash[i])
	}
	return result
}

func (h Sha256) Or(h1 Sha256) Sha256 {
	result := Sha256{}
	for i := range result.Hash {
		result.Hash[i] = h.Hash[i] ^ h1.Hash[i]
	}
	return result
}

func (h Sha256) Equals(h1 Sha256) bool {
	// idea to not use memcmp, from:
	//   https://lemire.me/blog/2018/08/22/avoid-lexicographical-comparisons-when-testing-for-string-equality/
	return h.Hash[0] == h1.Hash[0] &&
		h.Hash[1] == h1.Hash[1] &&
		h.Hash[2] == h1.Hash[2] &&
		h.Hash[3] == h1.Hash[3]
}

func Sha256Compare(a, b Sha256) int {
	for i := 0; i < 4; i++ {
		if a.Hash[i] > b.Hash[i] {
			return 1
		} else if a.Hash[i] < b.Hash[i] {
			return -1
		}
	}
	return 0
}
