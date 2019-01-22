package unittests

import (
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
	. "github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/unittests/test_contracts"
	"github.com/stretchr/testify/assert"
	"testing"
)

type IdentityTester struct {
	ValidatingTester
	test         *testing.T
	AbiSer       abi_serializer.AbiSerializer
	AbiSerTest   abi_serializer.AbiSerializer
	ProducerName string
}

func NewIdentityTester(test *testing.T) *IdentityTester {
	t := &IdentityTester{test: test}
	t.ValidatingTester = *newValidatingTester(true, chain.SPECULATIVE)

	t.ProduceBlocks(2, false)

	t.CreateAccounts([]common.AccountName{common.N("identity"), common.N("identitytest"), common.N("alice"), common.N("bob"), common.N("carol")}, false, true)
	t.ProduceBlocks(100, false)

	t.SetCode(common.N("identity"), wast2wasm(IdentityWast), nil)
	t.SetAbi(common.N("identity"), IdentityAbi, nil)
	t.SetCode(common.N("identitytest"), wast2wasm(IdentityTestWast), nil)
	t.SetAbi(common.N("identitytest"), IdentityTestAbi, nil)
	t.ProduceBlocks(1, false)

	accnt := t.Control.GetAccount(common.N("identity"))
	var abi abi_serializer.AbiDef
	assert.True(test, abi_serializer.ToABI(accnt.Abi, &abi))
	t.AbiSer.SetAbi(&abi, t.AbiSerializerMaxTime)

	acntTest := t.Control.GetAccount(common.N("identitytest"))
	var abiTest abi_serializer.AbiDef
	assert.True(test, abi_serializer.ToABI(acntTest.Abi, &abiTest))
	t.AbiSerTest.SetAbi(&abiTest, t.AbiSerializerMaxTime)

	ap := t.Control.ActiveProducers()
	FcAssert(0 < len(ap.Producers), "No producers")
	t.ProducerName = ap.Producers[0].ProducerName.String()

	return t
}

func (t *IdentityTester) GetResultUint64() uint64 {
	db := t.Control.DataBase()
	tid := TableIdObject{Code: common.N("identitytest"), Scope: 0, Table: common.N("result")}
	FcAssert(db.Find("byCodeScopeTable", tid, &tid) == nil, "Table id not found")

	idx, _ := db.GetIndex("byScopePrimary", KeyValueObject{})

	obj := KeyValueObject{TId: tid.ID}
	itr, _ := idx.LowerBound(obj)
	if !itr.End() {
		_ = itr.Data(&obj)
	}
	FcAssert(!itr.End() && obj.TId == tid.ID, "lower_bound failed")

	FcAssert(uint64(common.N("result")) == obj.PrimaryKey, "row with result not found")
	FcAssert(obj.Value.Size() == 8, "unexpected result size")
	return binary.LittleEndian.Uint64(obj.Value)
}

func (t *IdentityTester) GetOwnerForIdentity(identity uint64) uint64 {
	getOwnerAct := types.Action{
		Account: common.N("identitytest"),
		Name:    common.N("getowner"),
		Data: t.AbiSerTest.VariantToBinary("getowner", &common.Variants{
			"identity": identity,
		}, t.AbiSerializerMaxTime),
	}
	assert.Equal(t.test, t.Success(), t.PushAction(&getOwnerAct, common.N("alice")))
	return t.GetResultUint64()
}

func (t *IdentityTester) GetIdentityForAccount(account string) uint64 {
	getIdentityAct := types.Action{
		Account: common.N("identitytest"),
		Name:    common.N("getidentity"),
		Data: t.AbiSerTest.VariantToBinary("getidentity", &common.Variants{
			"account": account,
		}, t.AbiSerializerMaxTime),
	}
	assert.Equal(t.test, t.Success(), t.PushAction(&getIdentityAct, common.N("alice")))
	return t.GetResultUint64()
}

func (t *IdentityTester) CreateIdentity(accountName string, identity uint64, auth bool /*= true*/) ActionResult {
	createAct := types.Action{
		Account: common.N("identity"),
		Name:    common.N("create"),
		Data: t.AbiSer.VariantToBinary("create", &common.Variants{
			"creator":  accountName,
			"identity": identity,
		}, t.AbiSerializerMaxTime),
	}
	var authorizer common.AccountName
	if auth {
		authorizer = common.N(accountName)
	} else {
		if accountName == "bob" {
			authorizer = common.N("alice")
		} else {
			authorizer = common.N("bob")
		}
	}
	return t.PushAction(&createAct, authorizer)
}

func (t *IdentityTester) GetIdentity(idnt uint64) common.Variants {
	db := t.Control.DataBase()
	tid := TableIdObject{Code: common.N("identity"), Scope: common.N("identity"), Table: common.N("ident")}
	FcAssert(db.Find("byCodeScopeTable", tid, &tid) == nil, "object id not found")

	obj := KeyValueObject{TId: tid.ID, PrimaryKey: idnt}
	idx, _ := db.GetIndex("byScopePrimary", obj)

	itr, _ := idx.LowerBound(obj)
	if !itr.End() {
		itr.Data(&obj)
	}
	FcAssert(!itr.End() && obj.TId == tid.ID, "lower_bound failed")
	assert.Equal(t.test, idnt, obj.PrimaryKey)

	data := make([]byte, obj.Value.Size())
	copy(data, obj.Value)

	return t.AbiSer.BinaryToVariant("identrow", data, t.AbiSerializerMaxTime, false)
}

func (t *IdentityTester) Certify(certifier string, identity uint64, fields []common.Variants, auth bool /*= true*/) ActionResult {
	certAct := types.Action{
		Account: common.N("identity"),
		Name:    common.N("certprop"),
		Data: t.AbiSer.VariantToBinary("certprop", &common.Variants{
			"bill_storage_to": certifier,
			"certifier":       certifier,
			"identity":        identity,
			"value":           fields,
		}, t.AbiSerializerMaxTime),
	}
	var authorizer common.AccountName
	if auth {
		authorizer = common.N(certifier)
	} else {
		if certifier == "bob" {
			authorizer = common.N("alice")
		} else {
			authorizer = common.N("bob")
		}
	}
	return t.PushAction(&certAct, authorizer)
}

func (t *IdentityTester) GetCertrow(identity uint64, property string, trusted uint64, certifier string) common.Variants {
	db := t.Control.DataBase()
	tid := TableIdObject{Code: common.N("identity"), Scope: common.AccountName(identity), Table: common.N("certs")}
	if db.Find("byCodeScopeTable", tid, &tid) != nil {
		return nil
	}

	key := (&types.FixedKey{Size: 32}).MakeFromWordSequence(uint64(common.N(property)), trusted, uint64(common.N(certifier))).GetArray()
	obj := Idx256Object{TId: tid.ID, SecondaryKey: eos_math.Uint256{Low: key[0], High: key[1]}}
	idx, _ := db.GetIndex("bySecondary", obj)
	itr, _ := idx.LowerBound(obj)
	if !itr.End() {
		itr.Data(&obj)
	}
	if !itr.End() && obj.TId == tid.ID && obj.SecondaryKey == (eos_math.Uint256{Low: key[0], High: key[1]}) {
		primaryKey := obj.PrimaryKey
		obj := KeyValueObject{TId: tid.ID, PrimaryKey: primaryKey}
		idx, _ := db.GetIndex("byScopePrimary", obj)

		itr, _ := idx.LowerBound(obj)
		if !itr.End() {
			itr.Data(&obj)
		}
		FcAssert(!itr.End() && obj.TId == tid.ID && primaryKey == obj.PrimaryKey, "Record found in secondary index, but not found in primary index.")

		data := make([]byte, obj.Value.Size())
		copy(data, obj.Value)
		return t.AbiSer.BinaryToVariant("certrow", data, t.AbiSerializerMaxTime, false)

	} else {
		return nil
	}
}

func (t *IdentityTester) GetAccountrow(account string) common.Variants {
	db := t.Control.DataBase()
	acnt := common.N(account)
	tid := TableIdObject{Code: common.N("identity"), Scope: acnt, Table: common.N("account")}
	if db.Find("byCodeScopeTable", tid, &tid) != nil {
		return nil
	}

	obj := KeyValueObject{TId: tid.ID, PrimaryKey: uint64(common.N("account"))}
	idx, _ := db.GetIndex("byScopePrimary", obj)
	itr, _ := idx.LowerBound(obj)
	if !itr.End() {
		itr.Data(&obj)
	}
	if !itr.End() && obj.TId == tid.ID && uint64(common.N("account")) == obj.PrimaryKey {
		data := make([]byte, obj.Value.Size())
		copy(data, obj.Value)
		return t.AbiSer.BinaryToVariant("accountrow", data, t.AbiSerializerMaxTime, false)

	} else {
		return nil
	}
}

func (t *IdentityTester) Settrust(trustor string, trusting string, trust uint64, auth bool /*= true*/) ActionResult {
	settrustAct := types.Action{
		Account: common.N("identity"),
		Name:    common.N("settrust"),
		Data: t.AbiSer.VariantToBinary("settrust", &common.Variants{
			"trustor":  trustor,
			"trusting": trusting,
			"trust":    trust,
		}, t.AbiSerializerMaxTime),
	}
	tr := common.N(trustor)
	if auth {
		return t.PushAction(&settrustAct, tr)
	}
	return t.PushAction(&settrustAct, 0)
}

func (t *IdentityTester) Gettrust(trustor string, trusting string) bool {
	db := t.Control.DataBase()
	tid := TableIdObject{Code: common.N("identity"), Scope: common.N(trustor), Table: common.N("trust")}
	if db.Find("byCodeScopeTable", tid, &tid) != nil {
		return false
	}
	tng := uint64(common.N(trusting))
	obj := KeyValueObject{TId: tid.ID, PrimaryKey: tng}
	idx, _ := db.GetIndex("byScopePrimary", obj)
	itr, _ := idx.LowerBound(obj)
	return !itr.End() && obj.TId == tid.ID && tng == obj.PrimaryKey //true if found
}

const identityVal uint64 = 0xffffffffffffffff //64-bit value

func asUint64(t interface{}) uint64 {
	switch tp := t.(type) {
	case string:
		return eos_math.MustParseUint64(tp)
	case float64:
		return uint64(tp)
	default:
		return 0
	}
}

func asString(t interface{}) string {
	switch tp := t.(type) {
	case string:
		return tp
	case float64:
		return fmt.Sprintf("%f", tp)
	default:
		return ""
	}
}

func TestIdentityCreate(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", identityVal, true))
	idnt := t.GetIdentity(identityVal)
	assert.Equal(test, identityVal, asUint64(idnt["identity"]))
	assert.Equal(test, "alice", asString(idnt["creator"]))

	//attempts to create already existing identity should fail
	assert.Equal(test, t.WasmAssertMsg("identity already exists"), t.CreateIdentity("alice", identityVal, true))
	assert.Equal(test, t.WasmAssertMsg("identity already exists"), t.CreateIdentity("bob", identityVal, true))

	//alice can create more identities
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", 2, true))
	idnt = t.GetIdentity(2)
	assert.Equal(test, uint64(2), asUint64(idnt["identity"]))
	assert.Equal(test, "alice", asString(idnt["creator"]))

	//bob can create an identity as well
	assert.Equal(test, t.Success(), t.CreateIdentity("bob", 1, true))

	//identity == 0 has special meaning, should be impossible to create
	assert.Equal(test, t.WasmAssertMsg("identity=0 is not allowed"), t.CreateIdentity("alice", 0, true))

	//creating adentity without authentication is not allowed
	assert.Equal(test, t.Error("missing authority of alice"), t.CreateIdentity("alice", 3, false))

	idntBob := t.GetIdentity(1)
	assert.Equal(test, uint64(1), asUint64(idntBob["identity"]))
	assert.Equal(test, "bob", asString(idntBob["creator"]))

	//previously created identity should still exist
	idnt = t.GetIdentity(identityVal)
	assert.Equal(test, identityVal, asUint64(idnt["identity"]))
}

func TestCertifyDecertify(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", identityVal, true))

	//alice (creator of the identity) certifies 1 property
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "name",
		"type":       "string",
		"data":       common.HexBytes("Alice Smith"),
		"memo":       "",
		"confidence": 100,
	}}, true))

	obj := t.GetCertrow(identityVal, "name", 0, "alice")
	// check action
	assert.True(test, obj != nil)
	assert.Equal(test, "name", asString(obj["property"]))
	assert.Equal(test, uint64(0), asUint64(obj["trusted"]))
	assert.Equal(test, "alice", asString(obj["certifier"]))
	assert.Equal(test, uint64(100), asUint64(obj["confidence"]))
	assert.Equal(test, "string", asString(obj["type"]))
	assert.Equal(test, "Alice Smith", t.ToString(obj["data"]))

	//check that there is no trusted row with the same data
	assert.True(test, t.GetCertrow(identityVal, "name", 1, "alice") == nil)

	//bob certifies 2 properties
	fields := []common.Variants{{
		"property":   "email",
		"type":       "string",
		"data":       common.HexBytes("alice@alice.name"),
		"memo":       "official email",
		"confidence": 95,
	}, {
		"property":   "address",
		"type":       "string",
		"data":       common.HexBytes("1750 Kraft Drive SW, Blacksburg, VA 24060"),
		"memo":       "official address",
		"confidence": 80,
	}}

	//shouldn't be able to certify without authorization
	assert.Equal(test, t.Error("missing authority of bob"), t.Certify("bob", identityVal, fields, false))

	//certifying non-existent identity is not allowed
	nonExistent := uint64(21)
	assert.Equal(test, t.WasmAssertMsg("identity does not exist"), t.Certify("alice", nonExistent, []common.Variants{{
		"property":   "name",
		"type":       "string",
		"data":       common.HexBytes("Alice Smith"),
		"memo":       "",
		"confidence": 100,
	}}, true))

	//parameter "type" should be not longer than 32 bytes
	assert.Equal(test, t.WasmAssertMsg("certrow::type should be not longer than 32 bytes"), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "height",
		"type":       "super_long_type_name_which_is_not_allowed",
		"data":       common.HexBytes("Alice Smith"),
		"memo":       "",
		"confidence": 100,
	}}, true))

	//bob also should be able to certify
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, fields, true))

	obj = t.GetCertrow(identityVal, "email", 0, "bob")
	assert.True(test, obj != nil)
	assert.Equal(test, "email", asString(obj["property"]))
	assert.Equal(test, uint64(0), asUint64(obj["trusted"]))
	assert.Equal(test, "bob", asString(obj["certifier"]))
	assert.Equal(test, uint64(95), asUint64(obj["confidence"]))
	assert.Equal(test, "string", asString(obj["type"]))
	assert.Equal(test, "alice@alice.name", t.ToString(obj["data"]))

	obj = t.GetCertrow(identityVal, "address", 0, "bob")
	assert.True(test, obj != nil)
	assert.Equal(test, "address", asString(obj["property"]))
	assert.Equal(test, uint64(0), asUint64(obj["trusted"]))
	assert.Equal(test, "bob", asString(obj["certifier"]))
	assert.Equal(test, uint64(80), asUint64(obj["confidence"]))
	assert.Equal(test, "string", asString(obj["type"]))
	assert.Equal(test, "1750 Kraft Drive SW, Blacksburg, VA 24060", t.ToString(obj["data"]))

	//now alice certifies another email
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "email",
		"type":       "string",
		"data":       common.HexBytes("alice.smith@gmail.com"),
		"memo":       "",
		"confidence": 100,
	}}, true))

	obj = t.GetCertrow(identityVal, "email", 0, "alice")
	assert.True(test, obj != nil)
	assert.Equal(test, "email", asString(obj["property"]))
	assert.Equal(test, uint64(0), asUint64(obj["trusted"]))
	assert.Equal(test, "alice", asString(obj["certifier"]))
	assert.Equal(test, uint64(100), asUint64(obj["confidence"]))
	assert.Equal(test, "string", asString(obj["type"]))
	assert.Equal(test, "alice.smith@gmail.com", t.ToString(obj["data"]))

	//email certification made by bob should be still in place
	obj = t.GetCertrow(identityVal, "email", 0, "bob")
	assert.Equal(test, "bob", asString(obj["certifier"]))
	assert.Equal(test, "alice@alice.name", t.ToString(obj["data"]))

	//remove email certification made by alice
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "email",
		"type":       "string",
		"data":       common.HexBytes(""),
		"memo":       "",
		"confidence": 0,
	}}, true))
	assert.True(test, t.GetCertrow(identityVal, "email", 0, "alice") == nil)

	//email certification made by bob should still be in place
	obj = t.GetCertrow(identityVal, "email", 0, "bob")
	assert.Equal(test, "bob", asString(obj["certifier"]))
	assert.Equal(test, "alice@alice.name", t.ToString(obj["data"]))

	//name certification made by alice should still be in place
	obj = t.GetCertrow(identityVal, "name", 0, "alice")
	assert.True(test, obj != nil)
	assert.Equal(test, "Alice Smith", t.ToString(obj["data"]))
}

func TestTrustUntrust(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.Settrust("bob", "alice", 1, true))
	assert.Equal(test, true, t.Gettrust("bob", "alice"))

	//relation of trust in opposite direction should not exist
	assert.Equal(test, false, t.Gettrust("alice", "bob"))

	//remove trust
	assert.Equal(test, t.Success(), t.Settrust("bob", "alice", 0, true))
	assert.Equal(test, false, t.Gettrust("bob", "alice"))
}

func TestCertifyDecertifyOwner(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", identityVal, true))

	//bob certifies ownership for alice
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "bob") != nil)

	//it should not affect "account" singleton in alice's scope since it's not self-certification
	assert.Equal(test, true, t.GetAccountrow("alice") == nil)
	//it also shouldn't affect "account" singleton in bob's scope since he certified not himself
	assert.Equal(test, true, t.GetAccountrow("bob") == nil)

	// alice certifies her ownership, should populate "account" singleton
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	obj := t.GetCertrow(identityVal, "owner", 0, "alice")
	assert.True(test, obj != nil)
	assert.Equal(test, "owner", asString(obj["property"]))
	assert.Equal(test, uint64(0), asUint64(obj["trusted"]))
	assert.Equal(test, "alice", asString(obj["certifier"]))
	assert.Equal(test, uint64(100), asUint64(obj["confidence"]))
	assert.Equal(test, "account", asString(obj["type"]))
	assert.Equal(test, uint64(common.N("alice")), t.ToUint64(obj["data"]))

	//check that singleton "account" in the alice's scope contains the identity
	obj = t.GetAccountrow("alice")
	assert.True(test, obj != nil)
	assert.Equal(test, identityVal, asUint64(obj["identity"]))

	// ownership was certified by alice, but not by a block producer or someone trusted by a block producer
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))

	//remove bob's certification
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 0,
	}}, true))

	//singleton "account" in the alice's scope still should contain the identity
	obj = t.GetAccountrow("alice")
	assert.True(test, obj != nil)
	assert.Equal(test, identityVal, asUint64(obj["identity"]))

	//remove owner certification
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 0,
	}}, true))
	obj = t.GetCertrow(identityVal, "owner", 0, "alice")
	assert.True(test, obj == nil)

	//check that singleton "account" in the alice's scope doesn't contain the identity
	obj = t.GetAccountrow("alice")
	assert.True(test, obj == nil)
}

func TestOwnerCertifiedByProducer(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", identityVal, true))

	// certify owner by a block producer, should result in trusted certification
	assert.Equal(test, t.Success(), t.Certify(t.ProducerName, identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "",
		"confidence": 100,
	}}, true))

	obj := t.GetCertrow(identityVal, "owner", 1, t.ProducerName)
	assert.True(test, obj != nil)
	assert.Equal(test, "owner", asString(obj["property"]))
	assert.Equal(test, uint64(1), asUint64(obj["trusted"]))
	assert.Equal(test, t.ProducerName, asString(obj["certifier"]))
	assert.Equal(test, uint64(common.N("alice")), t.ToUint64(obj["data"]))

	//uncertified copy of that row shouldn't exist
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, t.ProducerName) == nil)

	//alice still has not claimed the identity - she is not the official owner yet
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))

	//alice claims it
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "alice") != nil)

	//now alice should be the official owner
	assert.Equal(test, uint64(common.N("alice")), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, identityVal, t.GetIdentityForAccount("alice"))

	//block producer decertifies ownership
	assert.Equal(test, t.Success(), t.Certify(t.ProducerName, identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "",
		"confidence": 0,
	}}, true))
	//self-certification made by alice still exists
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "alice") != nil)

	//but now she is not official owner
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))
}

func TestOwnerCertifiedByTrustedAccount(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("bob", identityVal, true))

	//alice claims the identity created by bob
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "alice") != nil)

	//alice claimed the identity, but it hasn't been certified yet - she is not the official owner
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))

	//block producer trusts bob
	assert.Equal(test, t.Success(), t.Settrust(t.ProducerName, "bob", 1, true))
	assert.Equal(test, true, t.Gettrust(t.ProducerName, "bob"))

	// bob (trusted account) certifies alice's ownership, it should result in trusted certification
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 1, "bob") != nil)
	//no untrusted row should exist
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "bob") == nil)

	//now alice should be the official owner
	assert.Equal(test, uint64(common.N("alice")), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, identityVal, t.GetIdentityForAccount("alice"))

	//block producer stops trusting bob
	assert.Equal(test, t.Success(), t.Settrust(t.ProducerName, "bob", 0, true))
	assert.Equal(test, false, t.Gettrust(t.ProducerName, "bob"))

	//certification made by bob is still flaged as trusted
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 1, "bob") != nil)

	//but now alice shouldn't be the official owner
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))
}

func TestOwnerCertificationBecomesTrusted(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("bob", identityVal, true))

	// bob (not trusted so far) certifies alice's ownership, it should result in untrusted certification
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "bob") != nil)
	//no trusted row should exist
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 1, "bob") == nil)

	//alice claims the identity created by bob
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "alice") != nil)
	//alice claimed the identity, but it is certified by untrusted accounts only - she is not the official owner
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))

	//block producer trusts bob
	assert.Equal(test, t.Success(), t.Settrust(t.ProducerName, "bob", 1, true))
	assert.Equal(test, true, t.Gettrust(t.ProducerName, "bob"))

	//old certification made by bob still shouldn't be flaged as trusted
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 0, "bob") != nil)

	//but effectively bob's certification should became trusted
	assert.Equal(test, uint64(common.N("alice")), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, identityVal, t.GetIdentityForAccount("alice"))
}

func TestOwnershipContradiction(test *testing.T) {
	t := NewIdentityTester(test)
	assert.Equal(test, t.Success(), t.CreateIdentity("alice", identityVal, true))

	//alice claims identity
	assert.Equal(test, t.Success(), t.Certify("alice", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))

	// block producer certifies alice's ownership
	assert.Equal(test, t.Success(), t.Certify(t.ProducerName, identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("alice"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 1, t.ProducerName) != nil)

	//now alice is the official owner of the identity
	assert.Equal(test, uint64(common.N("alice")), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, identityVal, t.GetIdentityForAccount("alice"))

	//bob claims identity
	assert.Equal(test, t.Success(), t.Certify("bob", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("bob"))),
		"memo":       "claiming onwership",
		"confidence": 100,
	}}, true))

	//block producer trusts carol
	assert.Equal(test, t.Success(), t.Settrust(t.ProducerName, "carol", 1, true))
	assert.Equal(test, true, t.Gettrust(t.ProducerName, "carol"))

	//another trusted delegate certifies bob's identity (to the identity already certified to alice)
	assert.Equal(test, t.Success(), t.Certify("carol", identityVal, []common.Variants{{
		"property":   "owner",
		"type":       "account",
		"data":       t.Uint64ToUint8Vector(uint64(common.N("bob"))),
		"memo":       "",
		"confidence": 100,
	}}, true))
	assert.Equal(test, true, t.GetCertrow(identityVal, "owner", 1, t.ProducerName) != nil)

	//now neither alice or bob are official owners, because we have 2 trusted certifications in contradiction to each other
	assert.Equal(test, uint64(0), t.GetOwnerForIdentity(identityVal))
	assert.Equal(test, uint64(0), t.GetIdentityForAccount("alice"))
}
