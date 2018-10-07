package common

import (
	"fmt"
	"strings"
)

/**
*  S Converts a base32 string to a uint64_t. This is a constexpr so that
*  this method can be used in template arguments as well.
*
*  @brief Converts a base32 string to a uint64_t.
*  @param str - String representation of the name
*  @return constexpr uint64_t - 64-bit unsigned integer representation of the name
*  @ingroup types
 */
func N(s string) (val uint64) {
	// ported from the eosio codebase, libraries/chain/include/eosio/chain/name.hpp
	var i uint32
	sLen := uint32(len(s))
	if sLen > 13 {
		//panic(fmt.Sprintf("Name is loger than 13 chacacters %s",s))
		fmt.Printf("Name is loger than 13 chacacters %s\n", s)
	}
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

	return
}

func S(in uint64) string {
	// ported from libraries/chain/name.cpp in eosio
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

/**
*  charToSymbol Converts a base32 symbol into its binary representation, used by string_to_name()
*
*  @brief Converts a base32 symbol into its binary representation, used by string_to_name()
*  @param c - Character to be converted
*  @return constexpr char - Converted character
*  @ingroup types
 */
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
