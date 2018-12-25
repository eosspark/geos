package unittests

import (
	"fmt"
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initializeBaseTester() (*AuthorizationManager, *BaseTester) {
	bt := newBaseTester(true, SPECULATIVE)
	am := bt.Control.Authorization
	return am, bt
}

func initializeValidatingTester() (*AuthorizationManager, *ValidatingTester) {
	vt := newValidatingTester(true, SPECULATIVE)
	am := vt.ValidatingControl.Authorization
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
	_, vt := initializeValidatingTester()
	vt.CreateAccount(common.N("alice"), common.DefaultConfig.SystemAccountName, false, true)
	vt.CreateAccount(common.N("bob"), common.DefaultConfig.SystemAccountName, false, true)

	// Deleting active or owner should fail
	Try(func() {
		vt.DeleteAuthority(
			common.N("alice"),
			common.N("active"),
			&[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.OwnerName}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("alice"), "owner")},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	}).End()
	Try(func() {
		vt.DeleteAuthority(
			common.N("alice"),
			common.N("owner"),
			&[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.OwnerName}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("alice"), "owner")},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	}).End()

	// Change owner permission
	newOwnerPrivKey := vt.getPrivateKey(common.N("alice"), "new_owner")
	newOwnerPubKey := newOwnerPrivKey.PublicKey()
	vt.SetAuthority2(
		common.N("alice"),
		common.N("owner"),
		types.NewAuthority(newOwnerPubKey, 0),
		common.N(""),
	)
	vt.ProduceBlocks(1, false)

	// Ensure the permission is updated
	po := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("owner")}
	err := vt.Control.DB.Find("byOwner", po, &po)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, po.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, po.Name, common.PermissionName(common.N("owner")))
	assert.Equal(t, po.Parent, int64(0))
	ownerId := po.ID
	auth := po.Auth.ToAuthority()
	assert.Equal(t, auth.Threshold, uint32(1))
	assert.Equal(t, len(auth.Keys), 1)
	assert.Equal(t, len(auth.Accounts), 0)
	assert.Equal(t, auth.Keys[0].Key, newOwnerPubKey)
	assert.Equal(t, auth.Keys[0].Weight, types.WeightType(1))

	// Change active permission, remember that the owner key has been changed
	newActivePrivKey := vt.getPrivateKey(common.N("alice"), "new_active")
	newActivePubKey := newActivePrivKey.PublicKey()
	vt.SetAuthority(
		common.N("alice"),
		common.N("active"),
		types.NewAuthority(newActivePubKey, 0),
		common.N("owner"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{vt.getPrivateKey(common.N("alice"), "active")},
	)
	vt.ProduceBlocks(1, false)

	obj := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("active")}
	err = vt.Control.DB.Find("byOwner", obj, &obj)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, obj.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, obj.Name, common.PermissionName(common.N("active")))
	assert.Equal(t, obj.Parent, ownerId)
	auth = obj.Auth.ToAuthority()
	assert.Equal(t, auth.Threshold, uint32(1))
	assert.Equal(t, len(auth.Keys), 1)
	assert.Equal(t, len(auth.Accounts), 0)
	assert.Equal(t, auth.Keys[0].Key, newActivePubKey)
	assert.Equal(t, auth.Keys[0].Weight, types.WeightType(1))

	spendingPrivKey := vt.getPrivateKey(common.N("alice"), "spending")
	spendingPubKey := spendingPrivKey.PublicKey()
	tradingPrivKey := vt.getPrivateKey(common.N("alice"), "trading")
	tradingPubKey := tradingPrivKey.PublicKey()

	// Bob attempts to create new spending auth for Alice
	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey, 0),
			common.N("active"),
			&[]types.PermissionLevel{{common.N("bob"), common.N("active")}},
			&[]ecc.PrivateKey{vt.getPrivateKey(common.N("bob"), "active")},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	})

	// Create new spending auth
	vt.SetAuthority(
		common.N("alice"),
		common.N("spending"),
		types.NewAuthority(spendingPubKey, 0),
		common.N("active"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)

	vt.ProduceBlocks(1, false)

	pObj := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("spending")}
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

	// Update spending auth parent to be its own, should fail
	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey, 0),
			common.N("spending"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	// Update spending auth parent to be owner, should fail
	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("spending"),
			types.NewAuthority(spendingPubKey, 0),
			common.N("owner"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	// Remove spending auth
	vt.DeleteAuthority(
		common.N("alice"),
		common.N("spending"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)
	delete := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("spending")}
	err = vt.Control.DB.Find("byOwner", delete, &delete)
	assert.NotNil(t, err)
	vt.ProduceBlocks(1, false)

	// Create new trading auth
	vt.SetAuthority(
		common.N("alice"),
		common.N("trading"),
		types.NewAuthority(tradingPubKey, 0),
		common.N("active"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)
	// Recreate spending auth again, however this time, it's under trading instead of owner
	vt.SetAuthority(
		common.N("alice"),
		common.N("spending"),
		types.NewAuthority(spendingPubKey, 0),
		common.N("trading"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("trading")}},
		&[]ecc.PrivateKey{tradingPrivKey},
	)

	// Verify correctness of trading and spending
	trading := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("trading")}
	spending := entity.PermissionObject{Owner: common.N("alice"), Name: common.N("spending")}
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
	err = vt.Control.DB.Find("id", tradingParent, &tradingParent)
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, tradingParent.Owner, common.AccountName(common.N("alice")))
	assert.Equal(t, tradingParent.Name, common.AccountName(common.N("active")))

	// Delete trading, should fail since it has children (spending)
	Try(func() {
		vt.DeleteAuthority(
			common.N("alice"),
			common.N("trading"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
			&[]ecc.PrivateKey{newActivePrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	// Update trading parent to be spending, should fail since changing parent authority is not supported
	Try(func() {
		vt.SetAuthority(
			common.N("alice"),
			common.N("trading"),
			types.NewAuthority(tradingPubKey, 0),
			common.N("spending"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("trading")}},
			&[]ecc.PrivateKey{tradingPrivKey},
		)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	// Delete spending auth
	vt.DeleteAuthority(
		common.N("alice"),
		common.N("spending"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)

	// Delete trading auth, now it should succeed since it doesn't have any children anymore
	vt.DeleteAuthority(
		common.N("alice"),
		common.N("trading"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("active")}},
		&[]ecc.PrivateKey{newActivePrivKey},
	)

	trading = entity.PermissionObject{Owner: common.N("alice"), Name: common.N("trading")}
	spending = entity.PermissionObject{Owner: common.N("alice"), Name: common.N("spending")}
	err = vt.Control.DB.Find("byOwner", trading, &trading)
	assert.NotNil(t, err)
	err = vt.Control.DB.Find("byOwner", spending, &spending)
	assert.NotNil(t, err)

	vt.close()
}

func TestLinkAuths(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)

	spendingPrivKey := vt.getPrivateKey(common.N("alice"), "spending")
	spendingPubKey := spendingPrivKey.PublicKey()
	scudPrivKey := vt.getPrivateKey(common.N("alice"), "scud")
	scudPubKey := scudPrivKey.PublicKey()

	vt.SetAuthority2(
		common.N("alice"),
		common.N("spending"),
		types.NewAuthority(spendingPubKey, 0),
		common.N("active"),
	)
	vt.SetAuthority2(
		common.N("alice"),
		common.N("scud"),
		types.NewAuthority(scudPubKey, 0),
		common.N("spending"),
	)

	// Send req auth action with alice's spending key, it should fail
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	}).End()
	// Link authority for eosio reqauth action with alice's spending key
	vt.LinkAuthority(common.N("alice"), common.N("eosio"), common.N("spending"), common.N("reqauth"))
	// Now, req auth action with alice's spending key should succeed
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
		&[]ecc.PrivateKey{spendingPrivKey},
	)

	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// Relink the same auth should fail
	Try(func() {
		vt.LinkAuthority(common.N("alice"), common.N("eosio"), common.N("spending"), common.N("reqauth"))
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	}).End()
	// Unlink alice with eosio reqauth
	vt.UnlinkAuthority(common.N("alice"), common.N("eosio"), common.N("reqauth"))
	// Now, req auth action with alice's spending key should fail
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
			&[]ecc.PrivateKey{spendingPrivKey},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	}).End()
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// Send req auth action with scud key, it should fail
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("scud")}},
			&[]ecc.PrivateKey{scudPrivKey},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	}).End()
	// Link authority for any eosio action with alice's scud key
	vt.LinkAuthority(common.N("alice"), common.N("eosio"), common.N("scud"), common.N(""))
	// Now, req auth action with alice's scud key should succeed
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("scud")}},
		&[]ecc.PrivateKey{scudPrivKey},
	)
	// req auth action with alice's spending key should also be fine, since it is the parent of alice's scud key
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
		&[]ecc.PrivateKey{spendingPrivKey},
	)
	vt.close()
}

func TestLinkThenUpdateAuth(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.CreateAccount(common.N("alice"), common.DefaultConfig.SystemAccountName, false, true)
	firstPrivKey := vt.getPrivateKey(common.N("alice"), "first")
	firstPubKey := firstPrivKey.PublicKey()
	secondPrivKey := vt.getPrivateKey(common.N("alice"), "second")
	secondPubKey := secondPrivKey.PublicKey()

	vt.SetAuthority2(common.N("alice"), common.N("first"), types.NewAuthority(firstPubKey, 0), common.N("active"))

	vt.LinkAuthority(common.N("alice"), common.N("eosio"), common.N("first"), common.N("reqauth"))
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("first")}},
		&[]ecc.PrivateKey{firstPrivKey},
	)

	vt.ProduceBlocks(13, false)

	// Update "first" auth public key
	vt.SetAuthority2(common.N("alice"), common.N("first"), types.NewAuthority(secondPubKey, 0), common.N("active"))
	// Authority updated, using previous "first" auth should fail on linked auth
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("first")}},
			&[]ecc.PrivateKey{firstPrivKey},
		)
	}).Catch(func(e UnsatisfiedAuthorization) {
		fmt.Println(e)
	}).End()
	// Using updated authority, should succeed
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("first")}},
		&[]ecc.PrivateKey{secondPrivKey},
	)
	vt.close()
}

func TestCreateAccount(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.CreateAccount(common.N("yc"), common.DefaultConfig.SystemAccountName, false, true)
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// Verify account create properly
	ycOwnerAuthority := entity.PermissionObject{Owner: common.N("yc"), Name: common.N("owner")}
	vt.Control.DB.Find("byOwner", ycOwnerAuthority, &ycOwnerAuthority)
	assert.Equal(t, ycOwnerAuthority.Auth.Threshold, uint32(1))
	assert.Equal(t, len(ycOwnerAuthority.Auth.Accounts), 1)
	assert.Equal(t, len(ycOwnerAuthority.Auth.Keys), 1)
	assert.Equal(t, ycOwnerAuthority.Auth.Keys[0].Key, vt.getPublicKey(common.N("yc"), "owner"))
	assert.Equal(t, ycOwnerAuthority.Auth.Keys[0].Weight, types.WeightType(1))

	ycActiveAuthority := entity.PermissionObject{Owner: common.N("yc"), Name: common.N("active")}
	vt.Control.DB.Find("byOwner", ycActiveAuthority, &ycActiveAuthority)
	assert.Equal(t, ycActiveAuthority.Auth.Threshold, uint32(1))
	assert.Equal(t, len(ycActiveAuthority.Auth.Accounts), 1)
	assert.Equal(t, len(ycActiveAuthority.Auth.Keys), 1)
	assert.Equal(t, ycActiveAuthority.Auth.Keys[0].Key, vt.getPublicKey(common.N("yc"), "active"))
	assert.Equal(t, ycActiveAuthority.Auth.Keys[0].Weight, types.WeightType(1))

	// Create duplicate name
	Try(func() {
		vt.CreateAccount(common.N("yc"), common.DefaultConfig.SystemAccountName, false, true)
	}).Catch(func(e AccountNameExistsException) {
		fmt.Println(e)
	})

	// Creating account with name more than 12 chars
	Try(func() {
		vt.CreateAccount(common.N("ychahahahahah"), common.DefaultConfig.SystemAccountName, false, true)
	}).Catch(func(e ActionValidateException) {
		fmt.Println(e)
	})

	// Create account with eosio. prefix with privileged account
	vt.CreateAccount(common.N("eosio.yc"), common.DefaultConfig.SystemAccountName, false, true)

	//Create account with eosio. prefix with non-privileged account, should fail
	Try(func() {
		vt.CreateAccount(common.N("eosio.hn"), common.N("yc"), false, true)
	}).Catch(func(e Exception) {
		fmt.Println(e)
	})
	vt.close()
}

func TestAnyAuth(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	aliceSpendingPrivKey := vt.getPrivateKey(common.N("alice"), "spending")
	aliceSpendingPubKey := aliceSpendingPrivKey.PublicKey()
	bobSpendingPrivKey := vt.getPrivateKey(common.N("bob"), "spending")
	bobSpendingPubKey := bobSpendingPrivKey.PublicKey()

	vt.SetAuthority2(common.N("alice"), common.N("spending"), types.NewAuthority(aliceSpendingPubKey, 0), common.N("active"))
	vt.SetAuthority2(common.N("bob"), common.N("spending"), types.NewAuthority(bobSpendingPubKey, 0), common.N("active"))

	// this should fail because spending is not active which is default for reqauth
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
			&[]ecc.PrivateKey{aliceSpendingPrivKey},
		)
	}).Catch(func(e IrrelevantAuthException) {
		fmt.Println(e)
	}).End()
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// link to eosio.any permission
	vt.LinkAuthority(common.N("alice"), common.N("eosio"), common.N("eosio.any"), common.N("reqauth"))
	vt.LinkAuthority(common.N("bob"), common.N("eosio"), common.N("eosio.any"), common.N("reqauth"))

	// this should succeed because eosio::reqauth is linked to any permission
	vt.PushReqAuth(
		common.N("alice"),
		&[]types.PermissionLevel{{common.N("alice"), common.N("spending")}},
		&[]ecc.PrivateKey{aliceSpendingPrivKey},
	)

	// this should fail because bob cannot authorize for alice, the permission given must be one-of alices
	Try(func() {
		vt.PushReqAuth(
			common.N("alice"),
			&[]types.PermissionLevel{{common.N("bob"), common.N("spending")}},
			&[]ecc.PrivateKey{bobSpendingPrivKey},
		)
	}).Catch(func(e MissingAuthException) {
		fmt.Println(e)
	}).End()
	vt.close()
}

func TestNoDoubleBilling(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	acc1 := common.AccountName(common.N("bill"))
	acc2 := common.AccountName(common.N("bill2"))
	acc1a := common.AccountName(common.N("bill1a"))

	vt.CreateAccount(acc1, common.DefaultConfig.SystemAccountName, false, true)
	vt.CreateAccount(acc1a, common.DefaultConfig.SystemAccountName, false, true)
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	createAcc := func(a common.AccountName) *types.TransactionTrace {
		trx := types.SignedTransaction{}
		vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)

		pls := []types.PermissionLevel{
			{acc1, common.N("active")},
			{acc1, common.N("owner")},
			{acc1a, common.N("owner")},
		}

		new := NewAccount{
			Creator: acc1,
			Name:    a,
			Owner:   types.NewAuthority(vt.getPublicKey(a, "owner"), 0),
			Active:  types.NewAuthority(vt.getPublicKey(a, "active"), 0),
		}
		data, _ := rlp.EncodeToBytes(new)
		act := &types.Action{
			Account:       new.GetAccount(),
			Name:          new.GetName(),
			Authorization: pls,
			Data:          data,
		}
		trx.Actions = append(trx.Actions, act)
		vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)
		chainId := vt.Control.GetChainId()
		pk := vt.getPrivateKey(acc1, "active")
		trx.Sign(&pk, &chainId)
		pk = vt.getPrivateKey(acc1, "owner")
		trx.Sign(&pk, &chainId)
		pk = vt.getPrivateKey(acc1a, "owner")
		trx.Sign(&pk, &chainId)
		return vt.PushTransaction(&trx, common.MaxTimePoint(), vt.DefaultBilledCpuTimeUs)
	}

	createAcc(acc2)
	usage := entity.ResourceUsageObject{Owner: acc1}
	vt.Control.DB.Find("byOwner", usage, &usage)
	usage2 := entity.ResourceUsageObject{Owner: acc1a}
	vt.Control.DB.Find("byOwner", usage2, &usage2)

	assert.True(t, usage.CpuUsage.Average() > 0)
	assert.True(t, usage.NetUsage.Average() > 0)
	assert.Equal(t, usage.CpuUsage.Average(), usage2.CpuUsage.Average())
	assert.Equal(t, usage.NetUsage.Average(), usage2.NetUsage.Average())
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	vt.close()
}

func TestStricterAuth(t *testing.T) {
	_, vt := initializeValidatingTester()
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	acc1 := common.AccountName(common.N("acc1"))
	acc2 := common.AccountName(common.N("acc2"))
	acc3 := common.AccountName(common.N("acc3"))
	acc4 := common.AccountName(common.N("acc4"))

	vt.CreateAccount(acc1, common.DefaultConfig.SystemAccountName, false, true)
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	createAcc := func(a common.AccountName, creator common.AccountName, threshold int) *types.TransactionTrace {
		trx := types.SignedTransaction{}
		vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)

		invalidAuth := types.Authority{
			Threshold: uint32(threshold),
			Keys:      []types.KeyWeight{{vt.getPublicKey(a, "owner"), 1}},
			Accounts:  []types.PermissionLevelWeight{{types.PermissionLevel{creator, common.DefaultConfig.ActiveName}, 1}},
			Waits:     []types.WaitWeight{},
		}

		pls := []types.PermissionLevel{{creator, common.N("active")}}

		new := NewAccount{
			Creator: creator,
			Name:    a,
			Owner:   types.NewAuthority(vt.getPublicKey(a, "owner"), 0),
			Active:  invalidAuth,
		}
		data, _ := rlp.EncodeToBytes(new)
		act := &types.Action{
			Account:       new.GetAccount(),
			Name:          new.GetName(),
			Authorization: pls,
			Data:          data,
		}
		trx.Actions = append(trx.Actions, act)
		vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)
		pk := vt.getPrivateKey(creator, "active")
		chainId := vt.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		return vt.PushTransaction(&trx, common.MaxTimePoint(), vt.DefaultBilledCpuTimeUs)
	}

	// Threshold can't be zero
	Try(func() {
		createAcc(acc2, acc1, 0)
	}).Catch(func(e Exception) {
		fmt.Println(e)
	})

	// Threshold can't be more than total weight
	Try(func() {
		createAcc(acc4, acc1, 3)
	}).Catch(func(e Exception) {
		fmt.Println(e)
	})

	createAcc(acc3, acc1, 1)
	vt.close()
}

func TestLinkAuthSpecial(t *testing.T) {
	// TODO
	_, vt := initializeValidatingTester()
	tester := common.AccountName(common.N("tester"))
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	vt.CreateAccount(common.N("currency"), common.N("eosio"), false, true)

	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	vt.CreateAccount(common.N("tester"), common.N("eosio"), false, true)
	vt.CreateAccount(common.N("tester2"), common.N("eosio"), false, true)
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	data := common.Variants{
		"account":    tester,
		"permission": common.N("first"),
		"parent":     common.N("active"),
		"auth":       types.NewAuthority(vt.getPublicKey(tester, "first"), 5),
	}

	actName := UpdateAuth{}.GetName()
	vt.PushAction2(
		&common.DefaultConfig.SystemAccountName,
		&actName,
		tester,
		&data,
		vt.DefaultExpirationDelta,
		0,
	)

	validateDisallow := func(rtype string) {
		actName := LinkAuth{}.GetName()
		data := common.Variants{
			"account":     tester,
			"code":        common.N("eosio"),
			"type":        common.N(rtype),
			"requirement": common.N("first"),
		}
		Try(func() {
			vt.PushAction2(
				&common.DefaultConfig.SystemAccountName,
				&actName,
				tester,
				&data,
				vt.DefaultExpirationDelta,
				0,
			)
		}).Catch(func(e ActionValidateException) {
			fmt.Println(e)
		})
	}

	validateDisallow("linkauth")
	validateDisallow("unlinkauth")
	validateDisallow("deleteauth")
	validateDisallow("updateauth")
	validateDisallow("canceldelay")
	vt.close()
}
