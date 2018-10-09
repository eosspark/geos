package common

type Uint256 struct {
	High  Uint128
	Low   Uint128
}

func (u *Uint256) IsZero() bool {
	if u.High.IsZero() && u.Low.IsZero() {
		return true
	}
	return false
}

func (u *Uint256) LeftShift() Uint256 {
	if u.Low.GetAt(127) {
		u.Low.LeftShift()
		u.High.LeftShift()
		u.High.Low += 1
	} else {
		u.Low.LeftShift()
		u.High.LeftShift()
	}
	return *u
}

func (u *Uint256) RightShift() Uint256 {
	if u.High.GetAt(0) {
		u.High.RightShift()
		u.Low.RightShift()
		u.Low.High += 0x01 << 63
	}
	return *u
}

func (u *Uint256) GetAt(i uint) bool {
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
