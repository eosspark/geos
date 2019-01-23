package unittests

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func nameSuffix(n uint64) uint64 {
	remainingBitsAfterLastActualDot := uint32(0)
	tmp := uint32(0)
	for remainingBits := 59; remainingBits >= 4; remainingBits -= 5 {
		// Get characters one-by-one in name in order from left to right (not including the 13th character)
		c := (n >> uint(remainingBits)) & 0x1F
		if c == 0 {
			tmp = uint32(remainingBits)
		} else {
			remainingBitsAfterLastActualDot = tmp
		}
	}
	thirteenthCharacter := n & 0x0F
	if thirteenthCharacter != 0 {
		remainingBitsAfterLastActualDot = tmp
	}

	if remainingBitsAfterLastActualDot == 0 {
		return n
	}

	mask := uint64(1<<remainingBitsAfterLastActualDot - 16)
	shift := 64 - remainingBitsAfterLastActualDot
	return (n&mask)<<shift + thirteenthCharacter<<(shift-1)
}

func fromString(s string) common.Asset {
	return common.Asset{}.FromString(&s)
}

func TestNameSuffixTests(t *testing.T) {
	assert.Equal(t, common.Name(0), common.Name(nameSuffix(uint64(common.N(common.S(0))))))
	assert.Equal(t, common.Name(common.N("abcdehijklmn")), common.Name(nameSuffix(uint64(common.N("abcdehijklmn")))))
	assert.Equal(t, common.Name(common.N("abcdehijklmn1")), common.Name(nameSuffix(uint64(common.N("abcdehijklmn1")))))
	assert.Equal(t, common.Name(common.N("def")), common.Name(nameSuffix(uint64(common.N("abc.def")))))
	assert.Equal(t, common.Name(common.N("def")), common.Name(nameSuffix(uint64(common.N(".abc.def")))))
	assert.Equal(t, common.Name(common.N("def")), common.Name(nameSuffix(uint64(common.N("..abc.def")))))
	assert.Equal(t, common.Name(common.N("def")), common.Name(nameSuffix(uint64(common.N("abc..def")))))
	assert.Equal(t, common.Name(common.N("ghi")), common.Name(nameSuffix(uint64(common.N("abc.def.ghi")))))
	assert.Equal(t, common.Name(common.N("abcdefghij")), common.Name(nameSuffix(uint64(common.N(".abcdefghij")))))
	assert.Equal(t, common.Name(common.N("1")), common.Name(nameSuffix(uint64(common.N(".abcdefghij.1")))))
	assert.Equal(t, common.Name(common.N("bcdefghij")), common.Name(nameSuffix(uint64(common.N("a.bcdefghij")))))
	assert.Equal(t, common.Name(common.N("1")), common.Name(nameSuffix(uint64(common.N("a.bcdefghij.1")))))
	assert.Equal(t, common.Name(common.N("c")), common.Name(nameSuffix(uint64(common.N("......a.b.c")))))
	assert.Equal(t, common.Name(common.N("123")), common.Name(nameSuffix(uint64(common.N("abcdefhi.123")))))
	assert.Equal(t, common.Name(common.N("123")), common.Name(nameSuffix(uint64(common.N("abcdefhij.123")))))
}

func TestAssetFromStringOverflow(t *testing.T) {
	a := common.Asset{}

	// precision = 19, magnitude < 2^61
	fromStringFunc := func() { fromString("0.1000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 19 should be <= 18", fromStringFunc)
	fromStringFunc = func() { fromString("-0.1000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 19 should be <= 18", fromStringFunc)
	fromStringFunc = func() { fromString("1.0000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 19 should be <= 18", fromStringFunc)
	fromStringFunc = func() { fromString("-1.0000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 19 should be <= 18", fromStringFunc)

	// precision = 18, magnitude < 2^58
	a = fromString("0.100000000000000000 CUR")
	assert.Equal(t, int64(100000000000000000), a.Amount)
	a = fromString("-0.100000000000000000 CUR")
	assert.Equal(t, int64(-100000000000000000), a.Amount)

	// precision = 18, magnitude = 2^62
	fromStringFunc = func() { fromString("4.611686018427387904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)
	fromStringFunc = func() { fromString("-4.611686018427387904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)
	fromStringFunc = func() { fromString("4611686018427387.904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)
	fromStringFunc = func() { fromString("-4611686018427387.904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)

	// precision = 18, magnitude = 2^62-1
	a = fromString("4.611686018427387903 CUR")
	assert.Equal(t, int64(4611686018427387903), a.Amount)
	a = fromString("-4.611686018427387903 CUR")
	assert.Equal(t, int64(-4611686018427387903), a.Amount)

	// precision = 0, magnitude = 2^62
	fromStringFunc = func() { fromString("4611686018427387904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)
	fromStringFunc = func() { fromString("-4611686018427387904 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62", fromStringFunc)

	// precision = 18, magnitude = 2^65
	fromStringFunc = func() { fromString("36.893488147419103232 CUR") }
	CheckThrowException(t, &exception.OverflowException{}, fromStringFunc)
	fromStringFunc = func() { fromString("-36.893488147419103232 CUR") }
	CheckThrowException(t, &exception.UnderflowException{}, fromStringFunc)

	// precision = 14, magnitude > 2^76
	fromStringFunc = func() { fromString("1000000000.00000000000000 CUR") }
	CheckThrowException(t, &exception.OverflowException{}, fromStringFunc)
	fromStringFunc = func() { fromString("-1000000000.00000000000000 CUR") }
	CheckThrowException(t, &exception.UnderflowException{}, fromStringFunc)

	// precision = 0, magnitude > 2^76
	fromStringFunc = func() { fromString("100000000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.ParseErrorException{}, "Couldn't parse int64", fromStringFunc)
	fromStringFunc = func() { fromString("-100000000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.ParseErrorException{}, "Couldn't parse int64", fromStringFunc)

	// precision = 20, magnitude > 2^142
	fromStringFunc = func() { fromString("100000000000000000000000.00000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 20 should be <= 18", fromStringFunc)
	fromStringFunc = func() { fromString("-100000000000000000000000.00000000000000000000 CUR") }
	CheckThrowExceptionAndMsg(t, &exception.SymbolTypeException{}, "precision 20 should be <= 18", fromStringFunc)
}

func TestAuthorityChecker(t *testing.T) {

	bt := newBaseTester(true, chain.SPECULATIVE)
	a := bt.getPublicKey(common.N("a"), "active")
	b := bt.getPublicKey(common.N("b"), "active")
	c := bt.getPublicKey(common.N("c"), "active")

	getNullAuthority := func(p *types.PermissionLevel) types.SharedAuthority { return types.SharedAuthority{} }

	makeAuthChecker := func(pta types.PermissionToAuthorityFunc, recursionDepthLimit uint16, pubkeys []ecc.PublicKey) types.AuthorityChecker {
		keySet := treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
		for _, key := range pubkeys {
			keySet.Add(key)
		}
		permissionLevelSet := treeset.NewWith(types.PermissionLevelType, types.ComparePermissionLevel)
		checkTime := func() {}
		return types.MakeAuthChecker(pta, recursionDepthLimit, keySet, permissionLevelSet, common.Microseconds(0), &checkTime)
	}

	A := types.SharedAuthority{Threshold: 2, Keys: []types.KeyWeight{{Key: a, Weight: 1}, {Key: b, Weight: 1}}}

	{
		checker := makeAuthChecker(getNullAuthority, uint16(2), []ecc.PublicKey{a, b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 2)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 0)
	}

	{
		checker := makeAuthChecker(getNullAuthority, uint16(2), []ecc.PublicKey{a, c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 0)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 2)
	}

	{
		checker := makeAuthChecker(getNullAuthority, uint16(2), []ecc.PublicKey{a, b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 2)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	{
		checker := makeAuthChecker(getNullAuthority, uint16(2), []ecc.PublicKey{b, c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 0)
	}

	A = types.SharedAuthority{Threshold: 3, Keys: []types.KeyWeight{{Key: a, Weight: 1}, {Key: b, Weight: 1}, {Key: c, Weight: 1}}}

	{
		checker := makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{c, b, a})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{a, b})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{a, c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{b, c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
	}

	A = types.SharedAuthority{Threshold: 1, Keys: []types.KeyWeight{{Key: a, Weight: 1}, {Key: b, Weight: 1}}}

	{
		checker := makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{a})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
	}

	A = types.SharedAuthority{Threshold: 1, Keys: []types.KeyWeight{{Key: a, Weight: 2}, {Key: b, Weight: 1}}}

	{
		checker := makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{a})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getNullAuthority, 2, []ecc.PublicKey{c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
	}

	getCAuthority := func(p *types.PermissionLevel) types.SharedAuthority {
		return types.SharedAuthority{Threshold: 1, Keys: []types.KeyWeight{{Key: c, Weight: 1}}}
	}

	A = types.SharedAuthority{
		Threshold: 2,
		Keys:      []types.KeyWeight{{Key: a, Weight: 2}, {Key: b, Weight: 1}},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, Weight: 1}},
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{a})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.AllKeysUsed())
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{b})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 0)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 0)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 2)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 0)
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{b, c, a})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 1)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 2)
	}

	A = types.SharedAuthority{
		Threshold: 3,
		Keys:      []types.KeyWeight{{a, 2}, {b, 1}},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, Weight: 3}},
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{a, b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.AllKeysUsed())
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{a, b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 1)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 2)
	}

	A = types.SharedAuthority{
		Threshold: 2,
		Keys:      []types.KeyWeight{{a, 1}, {b, 1}},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, Weight: 1}},
	}

	{
		checker := makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{a})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{b})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{c})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{a, b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{a, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{a, b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 2)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	A = types.SharedAuthority{
		Threshold: 2,
		Keys:      []types.KeyWeight{{a, 1}, {b, 1}},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, Weight: 2}},
	}

	{
		checker := makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{a, b})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{a})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getCAuthority, 2, []ecc.PublicKey{b})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
	}

	{
		checker := makeAuthChecker(getCAuthority, uint16(2), []ecc.PublicKey{a, b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 1)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 2)
	}

	d := bt.getPublicKey(common.N("d"), "active")
	e := bt.getPublicKey(common.N("e"), "active")

	getAuthority := func(p *types.PermissionLevel) types.SharedAuthority {
		if p.Actor == common.N("top") {
			return types.SharedAuthority{
				Threshold: 2,
				Keys:      []types.KeyWeight{{d, 1}},
				Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("bottom"), Permission: common.N("bottom")}, Weight: 1}},
			}
		}
		return types.SharedAuthority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{e, 1}},
		}
	}

	A = types.SharedAuthority{
		Threshold: 5,
		Keys:      []types.KeyWeight{{a, 2}, {b, 2}, {c, 2}},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("top"), Permission: common.N("top")}, Weight: 5}},
	}

	{
		checker := makeAuthChecker(getAuthority, uint16(2), []ecc.PublicKey{d, e})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.AllKeysUsed())
	}

	{
		checker := makeAuthChecker(getAuthority, uint16(2), []ecc.PublicKey{a, b, c, d, e})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 2)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 3)
	}

	{
		checker := makeAuthChecker(getAuthority, uint16(2), []ecc.PublicKey{a, b, c, e})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		usedKeys := checker.GetUsedKeys()
		assert.True(t, usedKeys.Size() == 3)
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	{
		checker := makeAuthChecker(getAuthority, 1, []ecc.PublicKey{a, b, c})
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		checker = makeAuthChecker(getAuthority, 1, []ecc.PublicKey{d, e})
		assert.True(t, !checker.SatisfiedAcd(&A, nil, 0))
	}

	assert.True(t, ecc.ComparePubKey(b, a) == -1)
	assert.True(t, ecc.ComparePubKey(b, c) == -1)
	assert.True(t, ecc.ComparePubKey(a, c) == -1)

	{
		// valid key order: b < a < c
		A = types.SharedAuthority{
			Threshold: 2,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
		}
		B := types.SharedAuthority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{b, 1}, {c, 1}},
		}
		C := types.SharedAuthority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{a, 1}, {c, 1}, {b, 1}},
		}
		D := types.SharedAuthority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{b, 1}, {c, 1}, {c, 1}},
		}
		E := types.SharedAuthority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{b, 1}, {b, 1}, {c, 1}},
		}
		F := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
		}
		checker := makeAuthChecker(getNullAuthority, uint16(2), []ecc.PublicKey{a, b, c})
		assert.True(t, types.Validate(A.ToAuthority()))
		assert.True(t, types.Validate(B.ToAuthority()))
		assert.True(t, !types.Validate(C.ToAuthority()))
		assert.True(t, !types.Validate(D.ToAuthority()))
		assert.True(t, !types.Validate(E.ToAuthority()))
		assert.True(t, !types.Validate(F.ToAuthority()))

		assert.True(t, !checker.AllKeysUsed())
		unusedKeys := checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 3)
		assert.True(t, checker.SatisfiedAcd(&A, nil, 0))
		assert.True(t, checker.SatisfiedAcd(&B, nil, 0))
		assert.True(t, !checker.AllKeysUsed())
		unusedKeys = checker.GetUnusedKeys()
		assert.True(t, unusedKeys.Size() == 1)
	}

	{
		A := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("a"), Permission: common.N("world")}, 1},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
				{types.PermissionLevel{Actor: common.N("hi"), Permission: common.N("world")}, 1},
			},
		}
		B := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
			},
		}
		C := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("there")}, 1},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
			},
		}
		D := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 2},
			},
		}
		E := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 2},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("there")}, 1},
			},
		}
		F := types.SharedAuthority{
			Threshold: 4,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("hi"), Permission: common.N("world")}, 2},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
			},
		}
		G := types.SharedAuthority{
			Threshold: 7,
			Keys:      []types.KeyWeight{{b, 1}, {a, 1}, {c, 1}},
			Accounts: []types.PermissionLevelWeight{
				{types.PermissionLevel{Actor: common.N("a"), Permission: common.N("world")}, 1},
				{types.PermissionLevel{Actor: common.N("hello"), Permission: common.N("world")}, 1},
				{types.PermissionLevel{Actor: common.N("hi"), Permission: common.N("world")}, 1},
			},
		}
		assert.True(t, types.Validate(A.ToAuthority()))
		assert.True(t, types.Validate(B.ToAuthority()))
		assert.True(t, types.Validate(C.ToAuthority()))
		assert.True(t, !types.Validate(D.ToAuthority()))
		assert.True(t, !types.Validate(E.ToAuthority()))
		assert.True(t, !types.Validate(F.ToAuthority()))
		assert.True(t, !types.Validate(G.ToAuthority()))
	}
	bt.close()
}

func TestTransactionTest(t *testing.T) {
	bt := newBaseTester(true, chain.SPECULATIVE)
	trx := types.SignedTransaction{}
	type params struct {
		From common.AccountName
	}
	ps := params{From: eosio}
	data, _ := rlp.EncodeToBytes(ps)
	act := types.Action{
		Account:       eosio,
		Name:          common.ActionName(common.N("reqauth")),
		Authorization: []types.PermissionLevel{{eosio, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	nonce := "dummy"
	data, _ = rlp.EncodeToBytes(nonce)
	contextFreeAct := types.Action{
		Account:       eosio,
		Name:          common.ActionName(common.N("nonce")),
		Authorization: []types.PermissionLevel{},
		Data:          data,
	}
	trx.ContextFreeActions = append(trx.ContextFreeActions, &contextFreeAct)
	bt.SetTransactionHeaders(&trx.Transaction, bt.DefaultExpirationDelta, 0)
	trx.Expiration = common.NewTimePointSecTp(common.Now())
	trx.Validate()
	assert.Equal(t, int(0), len(trx.Signatures))
	privKey := bt.getPrivateKey(eosio, "active")
	chainId := bt.Control.GetChainId()
	trx.SignWithoutAppend(privKey, &chainId)
	assert.Equal(t, int(0), len(trx.Signatures))
	trx.Sign(&privKey, &chainId)
	assert.Equal(t, int(1), len(trx.Signatures))
	trx.Validate()

	pkt := types.PackedTransaction{}
	pkt.SetTransaction(&trx.Transaction, types.CompressionNone)

	pkt2 := types.PackedTransaction{}
	pkt2.SetTransaction(&trx.Transaction, types.CompressionZlib)

	assert.True(t, trx.Expiration == pkt.Expiration())
	assert.True(t, trx.Expiration == pkt2.Expiration())

	assert.Equal(t, trx.ID(), pkt.ID())
	assert.Equal(t, trx.ID(), pkt2.ID())

	raw := pkt.GetRawTransaction()
	raw2 := pkt2.GetRawTransaction()
	assert.Equal(t, raw.Size(), raw2.Size())
	bt.close()
}
