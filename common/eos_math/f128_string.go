package eos_math

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Binary to decimal floating point conversion.
// Algorithm:
//   1) store mantissa in multiprecision decimal
//   2) shift decimal by exponent
//   3) read digits out & format
type decimalSlice struct {
	d      []byte
	nd, dp int
	neg    bool
}

type floatInfo struct {
	mantbits uint
	expbits  uint
	bias     int
}

//var float32info = floatInfo{23, 8, -127}
//var float64info = floatInfo{52, 11, -1023}
//var float80info = floatInfo{64, 15, -16383}
var float128info = floatInfo{112, 15, -16383}

type Mant struct {
	High uint64
	Low  uint64
}

func (f Float128) String() string {
	dst := make([]byte, 0)
	neg := f.High>>63 != 0
	exp := int(f.High>>48) & 0x7FFF
	mant := Mant{
		High: f.High & uint64(0x0000FFFFFFFFFFFF),
		Low:  f.Low,
	}
	flt := &float128info

	switch exp {
	case 1<<flt.expbits - 1:
		// Inf, NaN
		var s string
		switch {
		case mant.High != 0 && mant.Low != 0:
			s = "NaN"
		case neg:
			s = "-Inf"
		default:
			s = "+Inf"
		}
		return s

	case 0:
		// denormalized
		exp++

	default:
		// add implicit top bit
		mant.High |= uint64(1) << (flt.mantbits - 64) //f128 64
	}

	exp += flt.bias
	return string(bigFtoa(dst, 18, 'e', neg, mant, exp, flt))
}

// bigFtoa uses multiprecision computations to format a float.
func bigFtoa(dst []byte, prec int, fmt byte, neg bool, mant Mant, exp int, flt *floatInfo) []byte {
	d := new(decimal)
	d.Assign(mant)

	d.Shift(exp-int(flt.mantbits), mant)

	// Round appropriately.
	d.Round(prec + 1)

	digs := decimalSlice{d: d.d[:], nd: d.nd, dp: d.dp}

	return fmtE(dst, neg, digs, prec, fmt)
}

// %e: -d.ddddde±dd
func fmtE(dst []byte, neg bool, d decimalSlice, prec int, fmt byte) []byte {
	// sign
	if neg {
		dst = append(dst, '-')
	}

	// first digit
	ch := byte('0')
	if d.nd != 0 {
		ch = d.d[0]
	}
	dst = append(dst, ch)

	// .moredigits
	if prec > 0 {
		dst = append(dst, '.')
		i := 1
		m := min(d.nd, prec+1)
		if i < m {
			dst = append(dst, d.d[i:m]...)
			i = m
		}
		for ; i <= prec; i++ {
			dst = append(dst, '0')
		}
	}

	// e±
	dst = append(dst, fmt)
	exp := d.dp - 1
	if d.nd == 0 { // special case: 0 has exponent 0
		exp = 0
	}
	if exp < 0 {
		ch = '-'
		exp = -exp
	} else {
		ch = '+'
	}
	dst = append(dst, ch)

	// dd or ddd
	switch {
	case exp < 10:
		dst = append(dst, '0', byte(exp)+'0')
	case exp < 100:
		dst = append(dst, byte(exp/10)+'0', byte(exp%10)+'0')
	default:
		dst = append(dst, byte(exp/100)+'0', byte(exp/10)%10+'0', byte(exp%10)+'0')
	}

	return dst
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
