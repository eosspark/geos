package common
import (
	"fmt"
	"math"
)

type Uint128 struct {
	High uint64
	Low  uint64
}

func (u *Uint128) IsZero() bool {
	if u.Low == 0 && u.High == 0 {
		return true
	}
	return false
}

func (u *Uint128) LeftShift() Uint128 {
	if u.Low >> 63 == 1 {
		u.Low = u.Low << 1
		u.High = u.High << 1 + 1
	} else {
		u.Low = u.Low << 1
		u.High = u.High << 1
	}
	return *u
}

func (u *Uint128) RightShift() Uint128 {
	if u.High << 63 >> 63 == 1 {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
		u.Low += 0x01 << 63
	} else {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
	}

	return *u
}

func (u *Uint128) GetAt(i uint) bool {
	if i < 64 {
		return u.Low & ( 0x01 << i ) != 0
	} else {
		return u.High & ( 0x01 << (i - 64) ) != 0
	}
}

//if u > v , return 1; u < v , return -1; u = v , return 0 .
func (u *Uint128) Compare(v Uint128) int {
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

func (u *Uint128) Add(v Uint128) {
	if u.Low+v.Low < u.Low {
		u.High += v.High + 1
	} else {
		u.High += v.High
	}
	u.Low += v.Low
}

func (u *Uint128) Sub(v Uint128) {
	if u.Low >= v.Low {
		u.Low -= v.Low
		u.High -= v.High
	} else {
		u.Low += math.MaxUint64 - v.Low + 1
		u.High -= v.High + 1
	}
}

func (u *Uint128) Div(divisor Uint128) (Uint128, Uint128) {
	if divisor.IsZero() {
		fmt.Println("divisor cannot be zero")
	}
	Quotient := Uint128{}
	Remainder := Uint128{}
	for i := 0; i < 128; i++ {
		Remainder.LeftShift()
		Quotient.LeftShift()
		if u.GetAt(127 - uint(i)) {
			Remainder.Low ++
		}
		if Remainder.Compare(divisor) >= 0 {
			Quotient.Low ++
			Remainder.Sub(divisor)
		}
	}
	return Quotient, Remainder
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
	specialH := (uH+uL) >> 1 * (vH+vL) >> 1
	specialH = specialH >> 62 << 32

	if mulL + mixL < mulL {
		return Uint128{mulH + mixH + 1 + specialH, mulL + mixL}
	}
	return Uint128{mulH + mixH + specialH, mulL + mixL}
}

