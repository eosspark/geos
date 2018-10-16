package common

import "math"

type Int128 struct{
	High uint64
	Low  uint64
}

func (u Int128) IsZero() bool {
	if u.Low == 0 && u.High == 0 {
		return true
	}
	return false
}

func (u Int128) GetAt(i uint) bool {
	if i < 64 {
		return u.Low & ( 0x01 << i ) != 0
	} else {
		return u.High & ( 0x01 << (i - 64) ) != 0
	}
}

func (u *Int128) Set(i uint, b uint) {
	if i < 64 {
		if b == 1 {
			u.Low |= 0x01 << i
		}
		if b == 0 {
			u.Low &= math.MaxUint64 - 0x01 << i
		}
	}
	if i >= 64 {
		if b == 1 {
			u.High |= 0x01 << (i - 64)
		}
		if b == 0 {
			u.High &= math.MaxUint64 - 0x01 << (i - 64)
		}
	}
}

func (u *Int128) LeftShift() {
	if u.GetAt(63) {
		u.Low = u.Low << 1
		u.Set(64, 1)
	} else {
		u.Low = u.Low << 1
		u.High = u.High << 1
	}
}

func (u *Int128) LeftShifts(shift int){
	for i := 0; i < shift; i++ {
		u.LeftShift()
	}
}

func (u *Int128) RightShift() {
	signSymbol := u.GetAt(127)
	if u.GetAt(64) {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
		u.Set(63, 1)
	} else {
		u.High = u.High >> 1
		u.Low = u.Low >> 1
	}
	if signSymbol {
		u.Set(127, 1)
	}
}

func (u *Int128) RightShifts(shift int){
	for i := 0; i < shift; i++ {
		u.RightShift()
	}
}

func (u Int128) ToTrueForm() Uint128{
	if u.GetAt(127){
		for i := uint(0); i < 127; i++{
			if u.GetAt(i){
				u.Set(i, 0)
			} else {
				u.Set(i, 1)
			}
		}
		One := Int128{0,1}
		u = u.Add(One)
		u.Set(127,1)
	}
	return Uint128{u.High,u.Low}
}

func (Uint128) ToComplement() Int128{
	return Int128{}
}

func (u Int128) Add(v Int128) Int128{
	if u.Low+v.Low < u.Low {
		u.High += v.High + 1
	} else {
		u.High += v.High
	}
	u.Low += v.Low
	return u
}

func (u Int128) Sub(v Int128) Int128{
	if u.Low >= v.Low {
		u.Low -= v.Low
		u.High -= v.High
	} else {
		u.Low += math.MaxUint64 - v.Low + 1
		u.High -= v.High + 1
	}
	return u
}

func (u Int128) Mul(v Int128) Int128{
	signBit := false
	if u.GetAt(127) != v.GetAt(127){
		signBit = true
	}
	uTrueForm := u.ToTrueForm()
	vTrueForm := v.ToTrueForm()
	productTrueForm := uTrueForm.Mul(vTrueForm)
	if signBit == true{
		productTrueForm.Set(127,1)
	} else {
		productTrueForm.Set(127,0)
	}
	Product := productTrueForm.ToComplement()
	return Product
}

func (u Int128) Div(v Int128) (Int128, Int128){
	Quotient := Int128{}
	Remainder := Int128{}
	return Quotient, Remainder
}