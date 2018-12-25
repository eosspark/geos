package eos_math

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type Uint128 struct {
	Low  uint64
	High uint64
}

type Uint128Bytes struct {
	LowBytes  []byte
	HighBytes []byte
}

func (u Uint128) ToUint128Bytes() Uint128Bytes {
	lowBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lowBytes, u.Low)
	highBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(highBytes, u.High)
	return Uint128Bytes{lowBytes, highBytes}
}

func (u Uint128) IsZero() bool {
	if u.Low == 0 && u.High == 0 {
		return true
	}
	return false
}

func CreateUint128(i int) Uint128 {
	if i >= 0 {
		return Uint128{Low: uint64(i), High: 0}
	} else {
		fmt.Println("error")
		return Uint128{}
	}
}

func (u Uint128) GetAt(i uint) bool {
	if i < 64 {
		return u.Low&(0x01<<i) != 0
	} else {
		return u.High&(0x01<<(i-64)) != 0
	}
}

func (u *Uint128) Set(i uint, b uint) {
	if i < 64 {
		if b == 1 {
			u.Low |= 0x01 << i
		}
		if b == 0 {
			u.Low &= math.MaxUint64 - 0x01<<i
		}
	}
	if i >= 64 {
		if b == 1 {
			u.High |= 0x01 << (i - 64)
		}
		if b == 0 {
			u.High &= math.MaxUint64 - 0x01<<(i-64)
		}
	}
}

func MaxUint128() Uint128 {
	return Uint128{math.MaxUint64, math.MaxUint64}
}

func MinUint128() Uint128 {
	return Uint128{0, 0}
}

func (u *Uint128) LeftShift() {
	if u.GetAt(63) {
		u.Low = u.Low << 1
		u.High = u.High << 1
		u.Set(64, 1)
	} else {
		u.Low = u.Low << 1
		u.High = u.High << 1
	}
}

func (u *Uint128) LeftShifts(shift int) {
	for i := 0; i < shift; i++ {
		u.LeftShift()
	}
}

func (u *Uint128) RightShift() {
	if u.GetAt(64) {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
		u.Set(63, 1)
	} else {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
	}
}

func (u *Uint128) RightShifts(shift int) {
	for i := 0; i < shift; i++ {
		u.RightShift()
	}
}

//if u > v , return 1; u < v , return -1; u = v , return 0 .
func (u Uint128) Compare(v Uint128) int {
	if u.High > v.High {
		return 1
	} else if u.High < v.High {
		return -1
	}
	if u.Low > v.Low {
		return 1
	} else if u.Low < v.Low {
		return -1
	}
	return 0
}

func (u Uint128) Add(v Uint128) Uint128 {
	if u.Low+v.Low < u.Low {
		u.High += v.High + 1
	} else {
		u.High += v.High
	}
	u.Low += v.Low
	return u
}

func (u Uint128) Sub(v Uint128) Uint128 {
	if u.Compare(v) < 0 {
		fmt.Println("Uint128 cannot less than 0")
	}
	if u.Low >= v.Low {
		u.Low -= v.Low
		u.High -= v.High
	} else {
		u.Low += math.MaxUint64 - v.Low + 1
		u.High -= v.High + 1
	}
	return u
}

func (u Uint128) Mul(v Uint128) Uint128 {
	Product := MulUint64(u.Low, v.Low)
	tmp1 := MulUint64(u.High, v.Low).Low
	tmp2 := MulUint64(u.Low, v.High).Low
	tmp := Uint128{0, tmp1}.Add(Uint128{0, tmp2})
	Product = Product.Add(tmp)
	return Product
}

func (u Uint128) Div(divisor Uint128) (Uint128, Uint128) {
	if divisor.IsZero() {
		fmt.Println("divisor cannot be zero")
	}
	Quotient := Uint128{}
	Remainder := Uint128{}
	for i := 0; i < 128; i++ {
		Remainder.LeftShift()
		Quotient.LeftShift()
		if u.GetAt(127 - uint(i)) {
			Remainder.Low++
		}
		if Remainder.Compare(divisor) >= 0 {
			Quotient.Low++
			Remainder = Remainder.Sub(divisor)
		}
	}
	return Quotient, Remainder

	//new(big.Int).SetBytes()
	//Quotient := new(big.Int).Div(bigIntU, bigIntV)
	//if Quotient.BitLen() <= 64{
	//	return Uint128{0,Quotient.Uint64()}
	//} else {
	//	return Uint128{new(big.Int).Div(Quotient, bigIntMul).Uint64(), Quotient.Uint64()}
	//}
}

func (u Uint128) Sqrt() uint64 {
	if u.IsZero() {
		return 0
	}
	max := uint64(math.MaxUint64)
	min := uint64(1)
	m := max>>1 + min
	highUint128 := Uint128{}
	lowUint128 := Uint128{}
	product := Uint128{}
	for {
		product = MulUint64(m, m)
		if u.Compare(product) > 0 {
			min = m
		} else if u.Compare(product) < 0 {
			max = m
		} else {
			return m
		}
		highUint128.Low = max
		lowUint128.Low = min
		product = highUint128.Add(lowUint128)
		product.RightShift()
		m = product.Low
		if max == m || min == m {
			return m
		}
	}
}

func (u Uint128) String() string {
	uHigh := new(big.Int).SetUint64(u.High)
	uLow := new(big.Int).SetUint64(u.Low)

	uBigInt := new(big.Int).SetUint64(math.MaxUint64)
	one := new(big.Int).SetUint64(1)
	uBigInt = new(big.Int).Add(uBigInt, one)

	uBigInt = new(big.Int).Mul(uBigInt, uHigh)
	uBigInt = new(big.Int).Add(uBigInt, uLow)
	return uBigInt.String()
}

func MulUint64(u, v uint64) Uint128 {
	uH := u >> 32
	vH := v >> 32
	uL := u << 32 >> 32
	vL := v << 32 >> 32
	mulH := uH * vH
	mulL := uL * vL
	mulHL := (uH+uL)*(vH+vL) - mulH - mulL
	mixH := mulHL >> 32
	mixL := mulHL << 32 >> 32 << 32

	//(uH+uL)*(vH+vL) may more than maxUint64
	specialH := (uH + uL) >> 1 * (vH + vL) >> 1
	specialH = specialH >> 62 << 32

	if mulL+mixL < mulL {
		return Uint128{mulL + mixL, mulH + mixH + 1 + specialH}
	}
	return Uint128{mulL + mixL, mulH + mixH + specialH}
}
