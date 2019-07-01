package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

// ported from libraries/chain/name.cpp in eosio
type Name uint64
type AccountName = Name
type PermissionName = Name
type ActionName = Name
type TableName = Name
type ScopeName = Name

func (n Name) String() string {
	return S(uint64(n))
}

func (n Name) IsEmpty() bool {
	return n == 0
}

func (n Name) Pack() ([]byte, error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(n))
	return buf, nil
}

func (n *Name) Unpack(in []byte) (rlp.Unpack, error) {
	if len(in) < 8 {
		return nil, fmt.Errorf("rlp: uint64 required [%d] bytes, remaining [%d]", 8, len(in))
	}

	data := in[:8]
	out := binary.LittleEndian.Uint64(data)
	fmt.Println(Name(out))
	return nil, nil
}

//for treeset
var TypeName = reflect.TypeOf(Name(0))

func CompareName(first interface{}, second interface{}) int {
	if first.(Name) == second.(Name) {
		return 0
	}
	if first.(Name) < second.(Name) {
		return -1
	}
	return 1
}

func (n Name) Empty() bool {
	return n == 0
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(S(uint64(n)))
}

func (n *Name) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*n = Name(N(s))
	return nil
}

//N converts a base32 string to a uint64. 64-bit unsigned integer representation of the name.
func N(s string) Name {
	var i uint32
	var val uint64
	sLen := uint32(len(s))
	EosAssert(sLen <= 13, &exception.NameTypeException{}, "Name is longer than 13 characters (%s) ", s)
	for ; i <= 12; i++ {
		var c uint64
		if i < sLen {
			c = uint64(charToSymbol(s[i]))
		}

		if i < 12 {
			c &= 0x1f
			c <<= 64 - 5*(i+1)
		} else {
			c &= 0x0f
		}

		val |= c
	}

	return Name(val)
}

//S converts a uint64 to a base32 string. String representation of the name.
func S(in uint64) string {
	a := []byte{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.', '.'}

	tmp := in
	i := uint32(0)
	for ; i <= 12; i++ {
		bit := 0x1f
		if i == 0 {
			bit = 0x0f
		}
		c := base32Alphabet[tmp&uint64(bit)]
		a[12-i] = c

		shift := uint(5)
		if i == 0 {
			shift = 4
		}

		tmp >>= shift
	}

	return strings.TrimRight(string(a), ".")
}

func NameSuffix(n uint64) uint64 {
	remainingBitsAfterLastActualDot := uint32(0)
	tmp := 0
	for remainingBits := 59; remainingBits >= 4; remainingBits -= 5 { // int
		//Note: remaining_bits must remain signed integer
		//Get characters one-by-one in name in order from left to right (not including the 13th character)
		c := (n >> uint(remainingBits)) & 0x1F
		if c == 0 { // if this character is a dot
			tmp = remainingBits
		} else {
			remainingBitsAfterLastActualDot = uint32(tmp)
		}
	}

	thirteenthCharacter := n & 0x0F
	if thirteenthCharacter != 0 { // if 13th character is not a dot
		remainingBitsAfterLastActualDot = uint32(tmp)
	}
	if remainingBitsAfterLastActualDot == 0 { // there is no actual dot in the name other than potentially leading dots
		return n
	}
	// At this point remaining_bits_after_last_actual_dot has to be within the range of 4 to 59 (and restricted to increments of 5).
	// Mask for remaining bits corresponding to characters after last actual dot, except for 4 least significant bits (corresponds to 13th character).

	mask := uint64((1 << remainingBitsAfterLastActualDot) - 16)
	shift := uint32(64 - remainingBitsAfterLastActualDot)

	return uint64(((n & mask) << shift) + (thirteenthCharacter << (shift - 1)))
}

//charToSymbol converts a base32 symbol into its binary representation, used by N()
func charToSymbol(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 6
	}
	if c >= '1' && c <= '5' {
		return c - '1' + 1
	}
	return 0
}

var base32Alphabet = []byte(".12345abcdefghijklmnopqrstuvwxyz")
