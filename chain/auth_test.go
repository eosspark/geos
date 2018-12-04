package chain

import (
	"fmt"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"testing"
)

func initializeBaseTester() (*AuthorizationManager, *BaseTester) {
	bt := newBaseTester(true, SPECULATIVE)
	am := bt.Control.Authorization
	return am, bt
}

func initializeValidatingTester() (*AuthorizationManager, *ValidatingTester) {
	vt := newValidatingTester(true, SPECULATIVE)
	am := vt.ValidatingController.Authorization
	return am, vt
}

func TestMissingSigs(t *testing.T) {
	_, b := initializeBaseTester()
	b.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	Try(func() {
		b.PushReqAuth(common.N("alice"), &[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.ActiveName}}, &[]ecc.PrivateKey{})
	}).Catch(func(e UnsatisfiedAuthorization) {
		fmt.Println(e)
	}).End()
	/*trace := */ b.PushReqAuth2(common.N("alice"), "owner", false)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	//TODO: wait for controller::signal
	//assert.Equal(t,b.ChainHasTransaction(&trace.ID),true)
	b.close()
}

func TestMissingMultiSigs(t *testing.T) {
	_, b := initializeBaseTester()
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	b.CreateAccount(common.N("alice"), common.DefaultConfig.SystemAccountName, true, true)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	Try(func() {
		b.PushReqAuth2(common.N("alice"), "owner", false)
	}).Catch(func(e UnsatisfiedAuthorization) {
		fmt.Println(e)
	}).End()
	/*trace := */ b.PushReqAuth2(common.N("alice"), "owner", true)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	//TODO: wait for controller::signal
	//assert.Equal(t,b.ChainHasTransaction(&trace.ID),true)
	b.close()
}

func TestMissingAuths(t *testing.T) {
	_, b := initializeBaseTester()
	b.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	Try(func() {
		b.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("bob"), common.DefaultConfig.ActiveName}},
			&[]ecc.PrivateKey{b.getPrivateKey(common.N("bob"), "active")},
		)
	}).Catch(func(e MissingAuthException) {
		fmt.Println(e)
	}).End()
	b.close()
}

func TestDelegateAuth(t *testing.T) {
	a, b := initializeBaseTester()
	b.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
	delegatedAuth := types.SharedAuthority{
		Threshold: 1,
		Keys:      []types.KeyWeight{},
		Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: common.N("bob"), Permission: common.DefaultConfig.ActiveName}, Weight: 1}},
	}
	pk := b.getPrivateKey(common.N("alice"), "active")
	realAuth := types.SharedAuthority{
		Threshold: 1,
		Keys:      []types.KeyWeight{{pk.PublicKey(), 1}},
		Accounts: []types.PermissionLevelWeight{
			{Permission: types.PermissionLevel{Actor: common.N("alice"), Permission: common.DefaultConfig.EosioCodeName}, Weight: 1},
		},
	}
	originalAuth := a.GetPermission(&types.PermissionLevel{Actor: common.N("alice"), Permission: common.DefaultConfig.ActiveName}).Auth
	assert.Equal(t, originalAuth.Equals(realAuth), true)
	b.SetAuthority2(common.N("alice"), common.DefaultConfig.ActiveName, delegatedAuth.ToAuthority(), common.DefaultConfig.OwnerName)

	newAuth := a.GetPermission(&types.PermissionLevel{Actor: common.N("alice"), Permission: common.DefaultConfig.ActiveName}).Auth
	assert.Equal(t, newAuth.Equals(delegatedAuth), true)

	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs*2), 0)
	auth := a.GetPermission(&types.PermissionLevel{Actor: common.N("alice"), Permission: common.DefaultConfig.ActiveName}).Auth
	assert.Equal(t, newAuth.Equals(auth), true)

	b.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.ActiveName}},
		&[]ecc.PrivateKey{b.getPrivateKey(common.N("bob"), "active")},
	)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	b.close()
}

func TestUpdateAuths(t *testing.T) {
	_, vt := initializeBaseTester()
	vt.CreateAccount(common.N("alice"), common.DefaultConfig.SystemAccountName, false, true)
	vt.CreateAccount(common.N("bob"), common.DefaultConfig.SystemAccountName, false, true)
	Try(func() {
		vt.DeleteAuthority(
			common.N("alice"),
			common.N("active"),
			&[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.OwnerName}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("alice"),"owner")},
		)
	}).Catch(func(e ActionValidateException){
		fmt.Println(e)
	}).End()

	Try(func() {
		vt.DeleteAuthority(
			common.N("bob"),
			common.N("owner"),
			&[]types.PermissionLevel{{common.N("bob"), common.DefaultConfig.OwnerName}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("bob"),"owner")},
		)
	}).Catch(func(e ActionValidateException){
		fmt.Println(e)
	}).End()

	newOwnerPrivKey := vt.getPrivateKey(common.N("alice"),"new_owner")
	newOwnerPubKey := newOwnerPrivKey.PublicKey()
	vt.SetAuthority2(common.N("alice"),common.N("owner"),types.NewAuthority(newOwnerPubKey,0),common.N(""))
	vt.ProduceBlocks(1,false)

	po := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("new_owner")}
	err := vt.Control.DB.Find("byOwner", po, &po)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, po.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, po.Name, common.PermissionName(common.N("owner")))
	assert.Equal(t, po.Parent, 0)
	ownerId := po.ID
	auth:= po.Auth.ToAuthority()
	assert.Equal(t, auth.Threshold, 1)
	assert.Equal(t, len(auth.Keys), 1)
	assert.Equal(t, len(auth.Accounts), 0)
	assert.Equal(t, auth.Keys[0].Key, newOwnerPubKey)
	assert.Equal(t, auth.Keys[0].Weight, 1)

	newActivePrivKey := vt.getPrivateKey(common.N("alice"),"new_owner")
	newActivePubKey := newActivePrivKey.PublicKey()
	vt.SetAuthority2(common.N("alice"),common.N("owner"),types.NewAuthority(newActivePubKey,0),common.N(""))
	vt.ProduceBlocks(1,false)

	obj := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("new_active")}
	err = vt.Control.DB.Find("byOwner", obj, &obj)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, obj.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, obj.Name, common.PermissionName(common.N("active")))
	assert.Equal(t, obj.Parent, ownerId)
	auth = obj.Auth.ToAuthority()
	assert.Equal(t, auth.Threshold, 1)
	assert.Equal(t, len(auth.Keys), 1)
	assert.Equal(t, len(auth.Accounts), 0)
	assert.Equal(t, auth.Keys[0].Key, newActivePubKey)
	assert.Equal(t, auth.Keys[0].Weight, 1)

	spendingPrivKey := vt.getPrivateKey(common.N("alice"), "spending")
	spendingPubKey := spendingPrivKey.PublicKey()
	tradingPrivKey := vt.getPrivateKey(common.N("alice"),"trading")
	tradingPubKey := tradingPrivKey.PublicKey()

	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey,0),
			common.N("active"),
			&[]types.PermissionLevel{{common.N("bob"),common.N("active")}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("bob"),"active")},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	})

	vt.SetAuthority(
		common.N("alice"),
		common.N("spending"),
		types.NewAuthority(spendingPubKey,0),
		common.N("active"),
		&[]types.PermissionLevel{{common.N("alice"),common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)

	vt.ProduceBlocks(1,false)

	pObj := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("spending")}
	err = vt.Control.DB.Find("byOwner", pObj, &pObj)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, pObj.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, pObj.Name, common.PermissionName(common.N("spending")))

	parent := entity.PermissionObject{ID: pObj.Parent}
	err = vt.Control.DB.Find("id", parent, &parent)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, parent.Owner, common.N("alice"))
	assert.Equal(t, parent.Name, common.N("active"))

	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey,0),
			common.N("spending"),
			&[]types.PermissionLevel{{common.N("alice"),common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey,0),
			common.N("owner"),
			&[]types.PermissionLevel{{common.N("alice"),common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	vt.DeleteAuthority(
		common.N("alice"),
		common.N("spending"),
		&[]types.PermissionLevel{{common.N("alice"),common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)
	delete := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("spending")}
	err = vt.Control.DB.Find("byOwner", delete, &delete)
	if err != nil {
		fmt.Println(err)
	}
	vt.ProduceBlocks(1,false)

	trading := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("trading")}
	spending := entity.PermissionObject{Owner:common.N("alice"), Name:common.N("spending")}
	err = vt.Control.DB.Find("byOwner", trading, &trading)
	if err != nil {
		fmt.Println(err)
	}
	err = vt.Control.DB.Find("byOwner", spending, &spending)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, trading.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, spending.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, trading.Name, common.PermissionName(common.N("trading")))
	assert.Equal(t, spending.Name, common.PermissionName(common.N("spending")))
	assert.Equal(t, spending.Parent, trading.ID)

	tradingParent := entity.PermissionObject{ID: trading.Parent}
	err = vt.Control.DB.Find("byOwner", tradingParent, &tradingParent)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, tradingParent.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, tradingParent.Name, common.AccountName(common.N("active")))

	Try(func() {
		vt.DeleteAuthority(
			common.N("alice"),
			common.N("trading"),
			&[]types.PermissionLevel{{common.N("alice"),common.N("active")}},
			&[]ecc.PrivateKey{newActivePrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("trading"),
			types.NewAuthority(tradingPubKey,0),
			common.N("spending"),
			&[]types.PermissionLevel{{common.N("alice"),common.N("trading")}},
			&[]ecc.PrivateKey{tradingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	vt.DeleteAuthority(
		common.N("alice"),
		common.N("spending"),
		&[]types.PermissionLevel{{common.N("alice"),common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)
	vt.DeleteAuthority(
		common.N("alice"),
		common.N("trading"),
		&[]types.PermissionLevel{{common.N("alice"),common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)

	trading = entity.PermissionObject{Owner:common.N("alice"), Name:common.N("trading")}
	spending = entity.PermissionObject{Owner:common.N("alice"), Name:common.N("spending")}
	err = vt.Control.DB.Find("byOwner", trading, &trading)
	assert.Equal(t, err != nil, 0)
	err = vt.Control.DB.Find("byOwner", spending, &spending)
	assert.Equal(t, err != nil, 0)

	vt.close()
}
