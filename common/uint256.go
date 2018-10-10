package common

//import "math"

type Uint256 struct {
	High  Uint128
	Low   Uint128
}

func (u Uint256) IsZero() bool {
	if u.High.IsZero() && u.Low.IsZero() {
		return true
	}
	return false
}

func (u Uint256) LeftShift() Uint256 {
	if u.GetAt(127) {
		u.Low.LeftShift()
		u.High.LeftShift()
		u.Set(128, 1)
	} else {
		u.Low.LeftShift()
		u.High.LeftShift()
	}
	return u
}

func (u Uint256) RightShift() Uint256 {
	if u.GetAt(128) {
		u.High.RightShift()
		u.Low.RightShift()
		u.Set(127, 1)
	}
	return u
}

func (u Uint256) GetAt(i uint) bool {
	if i < 128 {
		return u.Low.GetAt(i)
	} else {
		return u.High.GetAt(i - 128)
	}
}

func (u *Uint256) Set(i uint, b uint) {
	if i < 128 {
		if b == 1 {
			u.Low.Set(i, 1)
		}
		if b == 0 {
			u.Low.Set(i, 0)
		}
	}
	if i >= 128 {
		if b == 1 {
			u.High.Set(i, 1)
		}
		if b == 0 {
			u.High.Set(i, 0)
		}
	}
}

func  (u Uint256) Compare(v Uint256) int {
	if u.High.Compare(v.High) > 0 {
		return 1
	} else if u.High.Compare(v.High) < 0 {
		return -1
	}
	if u.Low.Compare(v.Low) > 0 {
		return 1
	} else if u.Low.Compare(v.Low) < 0 {
		return -1
	}
	return 0
}

func (u Uint256) Add(v Uint256) Uint256{
	if u.Low.Add(v.Low).Compare(u.Low) < 0 {
		u.High = u.High.Add(v.High).Add(Uint128{0,1})
	} else {
		u.High = u.High.Add(v.High)
	}
	u.Low = u.Low.Add(v.Low)
	return u
}
//
//func (u *Uint256) Sub(v Uint256) Uint256{
//	if u.Low.Compare(v.Low) >= 0 {
//		u.Low.Sub(v.Low)
//		v.High.Sub(v.High)
//	} else {
//		u.Low.Sub(Uint128{math.MaxUint64,math.MaxUint64}.Sub(v.Low).Add(Uint128{0,1}))
//		u.High.Sub(v.High.Add(Uint128{0,1}))
//	}
//	return *u
//}
