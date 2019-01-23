package unittests

import (
	"testing"
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"github.com/eosspark/eos-go/exception"
	"io/ioutil"
	"github.com/eosspark/eos-go/entity"

	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/database"
	"fmt"
)

//eos-go db don't supply func db.get_index.size ,so add func--TakeCountOf
func TakeCountOf (generatedIndex *database.MultiIndex) (count int){
	it := generatedIndex.Begin()
	count = 0
	for generatedIndex.CompareIterator(it, generatedIndex.End()) == false  {
		it.Next()
		count++
	}
	return count
}


func TestValidating(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, SPECULATIVE)
		b.ProduceBlocks(2,false)

		trx := types.SignedTransaction{}

		newco := common.AccountName(common.N("newco"))
		creator := common.DefaultConfig.SystemAccountName
		ownerAuth := types.NewAuthority(b.getPublicKey(newco, "owner"), 0)

		pl := []types.PermissionLevel{{common.DefaultConfig.SystemAccountName, common.N("active")}}
		a := NewAccount{
			creator,
			newco,
			ownerAuth,
			types.NewAuthority(b.getPublicKey(newco, "active"), 0)}
		act := newAction(pl,&a)
		trx.Actions = append(trx.Actions,act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta,0)
		trx.DelaySec = 3
		privKey := b.getPrivateKey(creator,"active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey,&chainId)


		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)

		b.ProduceBlocks(8,false)
		b.close()
	})
}

func TestValidatingError(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, SPECULATIVE)
		b.ProduceBlocks(2,false)

		trx := types.SignedTransaction{}

		newco := common.N("newco")
		creator := common.DefaultConfig.SystemAccountName
		ownerAuth := types.NewAuthority(b.getPublicKey(newco, "owner"), 0)

		pl := []types.PermissionLevel{{common.DefaultConfig.SystemAccountName, common.PermissionName(common.N("active"))}}
		a := NewAccount{
			common.N("bad"), /// a does not exist, this should error when execute
			newco,
			ownerAuth,
			types.NewAuthority(b.getPublicKey(newco, "active"), 0)}
		act := newAction(pl,&a)
		trx.Actions = append(trx.Actions,act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta,0)
		trx.DelaySec = 3
		privKey := b.getPrivateKey(creator,"active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey,&chainId)

		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)


		b.ProduceBlocks(6,false)

		scheduledTrxs := b.Control.GetScheduledTransactions()
		assert.Equal(t, len(scheduledTrxs),1)
		dtrace := b.Control.PushScheduledTransaction(&scheduledTrxs[len(scheduledTrxs)-1],common.MaxTimePoint(),0)

		assert.Equal(t, dtrace.Except != nil, true)
		assert.Equal(t, dtrace.Except.Code(), exception.MissingAuthException{}.Code())
		b.close()

	})
}

//func GetCurrencyBalance(chain testing.T)

func TestLinkDelayDirect(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))
		chain.ProduceBlocks(1,false)
		chain.CreateAccount(eosioToken,eosio,false,true)
		chain.ProduceBlocks(10,false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1,false)
		chain.CreateAccount(common.N("tester"),eosio,false,true)
		chain.CreateAccount(common.N("tester2"), eosio, false,true)
		chain.ProduceBlocks(10, false)


		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)


		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "3.0000 CUR", "memo": "hi"},
			20,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(18, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "96.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "4.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.close()
	})
}


func TestDeleteAuth(t *testing.T){
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		// can't delete auth because it doesn't exist
		deleteAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
		}
		actName :=DeleteAuth{}.GetName()

		CheckThrowMsg(t, "Failed to retrieve permission", func() {
			chain.PushAction2(
				&eosio,
				&actName,
				tester,
				&deleteAuthData,
				chain.DefaultExpirationDelta,
				0,
			)
		})


		//update auth
		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		actName =UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		//link auth
		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		//permissionLink := entity.PermissionLinkObject{}
		//err := chain.Control.DB.Find("id", permissionLink, &permissionLink)
		//fmt.Println(permissionLink, " ", err)
		//create CUR token
		chain.ProduceBlocks(1, false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		// issue to account "eosio.token"
		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)


		// transfer from eosio.token to tester
		transfer := common.N("transfer")
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)


		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		//Permission := entity.PermissionObject{Owner: tester, Name: common.N("first")}
		//err = chain.Control.DB.Find("byOwner", Permission, &Permission)
		//fmt.Println(Permission, " ", err)

		// can't delete auth because it's linked
		deleteAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
		}
		actName =DeleteAuth{}.GetName()

		CheckThrowMsg(t, "Cannot delete a linked authority", func() {
			chain.PushAction2(
				&eosio,
				&actName,
				tester,
				&deleteAuthData,
				chain.DefaultExpirationDelta,
				0,
			)
		})
		//Permission = entity.PermissionObject{Owner: tester, Name: common.N("first")}
		//err = chain.Control.DB.Find("byOwner", Permission, &Permission)
		//fmt.Println(Permission, " ", err)

		//unlink auth
		unLinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
		}
		unLinkName := UnLinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&unLinkName,
			tester,
			&unLinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		// delete auth
		deleteAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
		}
		actName = DeleteAuth{}.GetName()

		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&deleteAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "3.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "96.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "4.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayDirectParentPermission(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)



		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)


		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("active"),
			"parent": common.N("owner"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "active"), 15),

		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "3.0000 CUR", "memo": "hi"},
			20,
			15)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))



		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(28, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "96.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "4.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayDirectWalkParentPermission(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))
		chain.ProduceBlocks(1,false)
		chain.CreateAccount(eosioToken,eosio,false,true)
		chain.ProduceBlocks(10,false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1,false)
		chain.CreateAccount(common.N("tester"),eosio,false,true)
		chain.CreateAccount(common.N("tester2"), eosio, false,true)
		chain.ProduceBlocks(10, false)


		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData2 := common.Variants{
			"account": tester,
			"permission": common.N("second"),
			"parent": common.N("first"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "second"), 0),

		}
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData2,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)


		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 20),

		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "3.0000 CUR", "memo": "hi"},
			30,
			20)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(38, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "96.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "4.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayPermissionChange(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			30,
			10,
		)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2 )
		assert.Equal(t,  len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// first transfer will finally be performed
		chain.ProduceBlocks(1,false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// delayed update auth removing the delay will finally execute
		chain.ProduceBlocks(1, false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)

		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(15, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)

		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "84.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "16.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayPermissionChangeWithDelayHeirarchy(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData2 := common.Variants{
			"account": tester,
			"permission": common.N("second"),
			"parent": common.N("first"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "second"), 0),

		}

		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData2,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			30,
			10,
		)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2 )
		assert.Equal(t,  len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// first transfer will finally be performed
		chain.ProduceBlocks(1,false)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)


		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(14, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "84.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "16.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayLinkChange(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData2 := common.Variants{
			"account": tester,
			"permission": common.N("second"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "second"), 0),

		}
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData2,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		LinkAuthData = common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName = LinkAuth{}.GetName()
		pl := []types.PermissionLevel{{tester, common.PermissionName(common.N("first"))}}

		CheckThrowMsg(t, "transaction declares authority", func() {
			chain.PushAction4(
				&eosio,
				&linkName,
				&pl,
				&LinkAuthData,
				30,
				3,
			)
		})

		// this transaction will be delayed 20 blocks
		LinkAuthData = common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName = LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			30,
			10,
		)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 2, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// first transfer will finally be performed
		chain.ProduceBlocks(1,false)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// delay on minimum permission of transfer is finally removed
		chain.ProduceBlocks(1, false)


		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "84.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "16.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()

	})
}


func TestLinkDelayUnlink(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)


		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		unLinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
		}
		unLinkName := UnLinkAuth{}.GetName()

		pl := []types.PermissionLevel{ {tester, common.PermissionName(common.N("first"))}}

		CheckThrowMsg(t, "transaction declares authority", func() {
			chain.PushAction4(
				&eosio,
				&unLinkName,
				&pl,
				&unLinkAuthData,
				30,
				7,
			)
		})


		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosio,
			&unLinkName,
			tester,
			&unLinkAuthData,
			30,
			10,
		)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2 )
		assert.Equal(t,  len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// first transfer will finally be performed
		chain.ProduceBlocks(1,false)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)


		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(15, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "84.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "16.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestLinkDelayLinkChangeHeirarchy(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData2 := common.Variants{
			"account": tester,
			"permission": common.N("second"),
			"parent": common.N("first"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),

		}

		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData2,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData3 := common.Variants{
			"account": tester,
			"permission": common.N("third"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "third"), 0),

		}
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData3,
			chain.DefaultExpirationDelta,
			0,
		)

		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		LinkAuthData = common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("third"),

		}
		trace = chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			30,
			10,
		)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2 )
		assert.Equal(t,  len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// first transfer will finally be performed
		chain.ProduceBlocks(1,false)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// delay on minimum permission of transfer is finally removed
		chain.ProduceBlocks(1, false)


		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "89.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "11.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "84.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "16.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}


func TestMindelay(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)


		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 0, TakeCountOf(generatedIndex))

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// send transfer with delay_sec set to 10
		accnt := entity.AccountObject{Name: eosioToken}
		chain.Control.DB.Find("byName", accnt, &accnt)
		abiEosioToken := accnt.GetAbi()
		var abis abi_serializer.AbiSerializer
		abis.SetAbi(abiEosioToken, chain.AbiSerializerMaxTime)

		actionTypeName := abis.GetActionType(transfer)

		act := types.Action{}
		act.Account = eosioToken
		act.Name = transfer
		act.Authorization =[]types.PermissionLevel{{tester, common.DefaultConfig.ActiveName}}
		data := common.Variants{"from": tester, "to": tester2,"quantity": "3.0000 CUR", "memo": "hi"}
		act.Data = abis.VariantToBinary(actionTypeName, &data, chain.AbiSerializerMaxTime)

		trx := types.SignedTransaction{}
		trx.Actions = append(trx.Actions, &act)

		chain.SetTransactionHeaders(&trx.Transaction, 30, 10)
		pk := chain.getPrivateKey(tester, "active")
		chainId := chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(18, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "99.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "1.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "96.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "4.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}

func TestCancelDelay(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		var ids []common.TransactionIdType

		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(10, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 10),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)


		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)


		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			10)
		ids = append(ids, trace.ID)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, TakeCountOf(generatedIndex))
		assert.Equal(t, 0, len(trace.ActionTraces))

		it := generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)

		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId: trxId}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s = "999900.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),
		}
		auth := []types.PermissionLevel{{ tester, common.N("first")}}

		CheckThrowMsg(t, "transaction declares authority", func() {
			trace = chain.PushAction4(
				&eosio,
				&actName,
				&auth,
				&updateAuthData,
				30,
				7,
			)
		})

		// this transaction will be delayed 20 blocks
		updateAuthData = common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 0),
		}
		trace = chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			30,
			10,
		)

		ids = append(ids, trace.ID)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2 )
		assert.Equal(t,  len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(16, false)
		fmt.Println("******")

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 20 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			10)
		for i, v := range ids {
			fmt.Println(i)
			fmt.Println(v)
		}
		ids = append(ids, trace.ID)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),3)
		assert.Equal(t, len(trace.ActionTraces),0)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// send canceldelay for first delayed transaction
		trx := types.SignedTransaction{}
		pl := []types.PermissionLevel{{tester, common.PermissionName(common.N("active"))}}
		cancelDelay := CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("active"))}, ids[0]}
		act := newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk := chain.getPrivateKey(tester, "active")
		chainId := chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2)


		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)



		chain.ProduceBlocks(1, false)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2)

		chain.ProduceBlocks(1,false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),2)

		chain.ProduceBlocks(1,false)
		// update auth will finally be performed

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)


		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// this transfer is performed right away since delay is removed
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "90.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "10.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		chain.ProduceBlocks(15, false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "90.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "10.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// second transfer finally is performed
		chain.ProduceBlocks(1, false)

		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "85.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "15.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		chain.close()
	})
}

// test canceldelay action under different permission levels
func TestCancelDelay2(t *testing.T) {
	t.Run("", func(t *testing.T) {
		_, chain := initializeValidatingTester()
		var ids []common.TransactionIdType

		tester := common.AccountName(common.N("tester"))
		tester2 := common.AccountName(common.N("tester2"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(eosioToken, eosio, false, true)
		chain.ProduceBlocks(1, false)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(common.N("tester"), eosio, false, true)
		chain.CreateAccount(common.N("tester2"), eosio, false, true)
		chain.ProduceBlocks(1, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 5),

		}
		actName :=UpdateAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		updateAuthData2 := common.Variants{
			"account": tester,
			"permission": common.N("second"),
			"parent": common.N("first"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "second"), 0),

		}
		chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData2,
			chain.DefaultExpirationDelta,
			0,
		)

		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)


		chain.ProduceBlocks(1,false)
		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.ProduceBlocks(1,false)
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": eosioToken, "quantity": "1000000.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)

		transfer := common.N("transfer")
		chain.ProduceBlocks(1,false)
		trace := chain.PushAction2(
			&eosioToken,
			&transfer,
			eosioToken,
			&common.Variants{"from": eosioToken, "to": tester,"quantity": "100.0000 CUR", "memo": "hi"},
			chain.DefaultExpirationDelta,
			0)
		assert.Equal(t, types.TransactionStatusExecuted,trace.Receipt.Status)
		
		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex), 0)

		chain.ProduceBlocks(1, false)

		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
		eosioTokenAccount := common.AccountName(eosioToken)
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &eosioTokenAccount)
		s :="999900.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		assert.Equal(t, expected, liquidBalance)

		// this transaction will be delayed 10 blocks
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "1.0000 CUR", "memo": "hi"},
			30,
			5)
		trxId := trace.ID

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, 1, generatedIndex)
		assert.Equal(t, 0, len(trace.ActionTraces))


		it := generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)
		
		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId:trace.ID}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// attempt canceldelay with wrong canceling_auth for delayed transfer of 1.0000 CUR
		trx := types.SignedTransaction{}
		pl := []types.PermissionLevel{{tester, common.PermissionName(common.N("active"))}}
		cancelDelay := CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("active"))}, trxId}
		act := newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk := chain.getPrivateKey(tester, "active")
		chainId := chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		CheckThrowMsg(t, "canceling_auth in canceldelay action was not found as authorization in the original delayed transaction", func() {
			trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)
		})


		// attempt canceldelay with "second" permission for delayed transfer of 1.0000 CUR
		trx = types.SignedTransaction{}
		pl = []types.PermissionLevel{{tester, common.PermissionName(common.N("second"))}}
		cancelDelay = CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("first"))}, trxId}
		act = newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk = chain.getPrivateKey(tester, "second")
		chainId = chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		CheckThrowMsg(t,  "canceldelay action declares irrelevant authority", func() {
			trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)
		})

		// canceldelay with "active" permission for delayed transfer of 1.0000 CUR
		trx = types.SignedTransaction{}
		pl = []types.PermissionLevel{{tester, common.PermissionName(common.N("active"))}}
		cancelDelay = CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("first"))}, trxId}
		act = newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk = chain.getPrivateKey(tester, "active")
		chainId = chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0 )


		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)
		
		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId: trxId}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		chain.ProduceBlocks(10, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		LinkAuthData = common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("second"),
		}
		linkName = LinkAuth{}.GetName()
		chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			30,
			5,
		)

		chain.ProduceBlocks(11, false)

		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": tester2,"quantity": "5.0000 CUR", "memo": "hi"},
			30,
			5)
		trxId = trace.ID
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)
		assert.Equal(t, len(trace.ActionTraces),0)

		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)

		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId: trxId}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		chain.ProduceBlocks(1, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// canceldelay with "first" permission for delayed transfer of 5.0000 CUR
		trx = types.SignedTransaction{}
		pl = []types.PermissionLevel{{tester, common.PermissionName(common.N("first"))}}
		cancelDelay = CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("second"))}, ids[0]}
		act = newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk = chain.getPrivateKey(tester, "first")
		chainId = chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0)

		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)

		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId: trxId}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		chain.ProduceBlocks(10, false)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		// this transaction will be delayed 10 blocks
		auth := []types.PermissionLevel{{ tester, common.PermissionName(common.N("owner"))}}
		trace = chain.PushAction4(
			&eosioToken,
			&transfer,
			&auth,
			&common.Variants{"from": tester, "to": tester2,"quantity": "10.0000 CUR", "memo": "hi"},
			30,
			5)
		trxId = trace.ID
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)
		assert.Equal(t, len(trace.ActionTraces), 0)

		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)

		//generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		//in :=entity.GeneratedTransaction{TrxId: trxId}
		//err := generatedIndex.Find(in, &in)
		//assert.Equal(t, err == nil, true)

		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

		// attempt canceldelay with "active" permission for delayed transfer of 10.0000 CUR
		trx = types.SignedTransaction{}
		pl = []types.PermissionLevel{{tester, common.PermissionName(common.N("active"))}}
		cancelDelay = CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("owner"))}, ids[0]}
		act = newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk = chain.getPrivateKey(tester, "active")
		chainId = chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		// attempt canceldelay with "active" permission for delayed transfer of 10.0000 CUR

		trx = types.SignedTransaction{}
		pl = []types.PermissionLevel{{tester, common.PermissionName(common.N("owner"))}}
		cancelDelay = CancelDelay{types.PermissionLevel{tester, common.PermissionName(common.N("owner"))}, ids[0]}
		act = newAction(pl, &cancelDelay)
		trx.Actions = append(trx.Actions, act)

		chain.SetTransactionHeaders(&trx.Transaction, chain.DefaultExpirationDelta, 0)
		pk = chain.getPrivateKey(tester, "owner")
		chainId = chain.Control.GetChainId()
		trx.Sign(&pk, &chainId)
		trace = chain.PushTransaction(&trx, common.MaxTimePoint(), chain.DefaultBilledCpuTimeUs)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0)

		it = generatedIndex.Begin()
		assert.Equal(t, it != generatedIndex.End(), true)



		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s = "100.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)
		liquidBalance = chain.GetCurrencyBalance(&eosioToken, &symbol, &tester2)
		s = "0.0000 CUR"
		expected = asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)

	})
}


func TestMaxTransactionDelayCreate(t *testing.T) {
	t.Run("", func(t *testing.T) {
		//assuming max transaction delay is 45 days (default in config.hpp)
		_, chain := initializeValidatingTester()

		tester := common.AccountName(common.N("tester"))

		chain.ProduceBlocks(1, false)
		chain.CreateAccount(tester, eosio, false, true)
		chain.ProduceBlocks(10, false)

		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 50*86400),// 50 days delay

		}
		actName :=UpdateAuth{}.GetName()
		CheckThrowMsg(t, "Cannot set delay longer than max_transacton_delay", func() {
			chain.PushAction2(
				&eosio,
				&actName,
				tester,
				&updateAuthData,
				chain.DefaultExpirationDelta,
				0,
			)
		})
		chain.close()
	})
}


func TestMaxTransactionDelayExecute(t *testing.T) {
	t.Run("", func(t *testing.T) {

		_, chain := initializeValidatingTester()

		tester := common.AccountName(common.N("tester"))

		chain.CreateAccount(eosioToken, eosio, false, true)

		eosioTokenWasm := "test_contracts/eosio.token.wasm"
		eosioTokenAbi := "test_contracts/eosio.token.abi"
		code, _ := ioutil.ReadFile(eosioTokenWasm)
		abi, _ := ioutil.ReadFile(eosioTokenAbi)
		chain.SetCode(eosioToken, code, nil)
		chain.SetAbi(eosioToken,abi,nil)

		chain.CreateAccount(common.N("tester"), eosio, false, true)


		create := common.N("create")
		chain.PushAction2(
			&eosioToken,
			&create,
			eosioToken,
			&common.Variants{"issuer": eosioToken, "maximum_supply": "9000000.0000 CUR"},
			chain.DefaultExpirationDelta,
			0)

		issue := common.N("issue")
		chain.PushAction2(
			&eosioToken,
			&issue,
			eosioToken,
			&common.Variants{"to": tester, "quantity": "100.0000 CUR", "memo": "for stuff"},
			chain.DefaultExpirationDelta,
			0)


		//create a permission level with delay 30 days and associate it with token transfer
		updateAuthData := common.Variants{
			"account": tester,
			"permission": common.N("first"),
			"parent": common.N("active"),
			"auth": types.NewAuthority(chain.getPublicKey(tester, "first"), 30*86400),// 30 days delay

		}
		actName :=UpdateAuth{}.GetName()
		trace :=chain.PushAction2(
			&eosio,
			&actName,
			tester,
			&updateAuthData,
			chain.DefaultExpirationDelta,
			0,
		)
		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
		LinkAuthData := common.Variants{
			"account": tester,
			"code": eosioToken,
			"type": common.N("transfer"),
			"requirement": common.N("first"),
		}
		linkName := LinkAuth{}.GetName()

		trace = chain.PushAction2(
			&eosio,
			&linkName,
			tester,
			&LinkAuthData,
			chain.DefaultExpirationDelta,
			0,
		)

		assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		//change max_transaction_delay to 60 sec ???
		chain.Control.DB.Modify(chain.Control.GetGlobalProperties(), func(gprops *entity.GlobalPropertyObject) {
			gprops.Configuration.MaxTrxDelay = 60
		})
		chain.ValidatingControl.DB.Modify(chain.ValidatingControl.GetGlobalProperties(), func(gprops *entity.GlobalPropertyObject) {
			gprops.Configuration.MaxTrxDelay = 60
		})

		chain.ProduceBlocks(1, false)

		//should be able to create transaction with delay 60 sec, despite permission delay being 30 days, because max_transaction_delay is 60 sec
		transfer := common.N("transfer")
		trace = chain.PushAction2(
			&eosioToken,
			&transfer,
			tester,
			&common.Variants{"from": tester, "to": eosioToken,"quantity": "9.0000 CUR", "memo": ""},
			120,
			60)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)

		chain.ProduceBlocks(1, false)

		generatedIndex, _ := chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),1)
		assert.Equal(t, len(trace.ActionTraces),0)

		//check that the delayed transaction executed after after 60 sec
		chain.ProduceBlocks(120, false)
		generatedIndex, _ = chain.Control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
		assert.Equal(t, TakeCountOf(generatedIndex),0)
		symbol := common.Symbol{Precision: 4, Symbol: "CUR"}
	    fmt.Println(chain.GetCurrencyBalance(&eosioToken, &symbol, &tester))

		//check that the transfer really happened
		symbol = common.Symbol{Precision: 4, Symbol: "CUR"}
		liquidBalance := chain.GetCurrencyBalance(&eosioToken, &symbol, &tester)
		s :="91.0000 CUR"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		assert.Equal(t, expected, liquidBalance)


		chain.close()

	})
}


func TestValidatingDelay(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newValidatingTester(true, SPECULATIVE)
		b.ProduceBlocks(2,false)
		//b.pro

		trx := types.SignedTransaction{}

		newco := common.N("newco")
		creator := common.DefaultConfig.SystemAccountName
		ownerAuth := types.NewAuthority(b.getPublicKey(newco, "owner"), uint32(b.AbiSerializerMaxTime))

		pl := []types.PermissionLevel{{common.DefaultConfig.SystemAccountName, common.PermissionName(common.N("active"))}}
		a := NewAccount{
			creator,
			newco,
			ownerAuth,
			types.NewAuthority(b.getPublicKey(newco, "active"), uint32(b.AbiSerializerMaxTime))}
		act := newAction(pl,&a)
		trx.Actions = append(trx.Actions,act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta,0)
		trx.DelaySec = 3
		trx.Expiration = common.NewTimePointSecTp(b.Control.HeadBlockTime().AddUs(common.Microseconds(1000000)))
		privKey := b.getPrivateKey(creator,"active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey,&chainId)

		trace := b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)

		sb := b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		sb = b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

		assert.Equal(t, types.TransactionStatusDelayed, trace.Receipt.Status)
		b.ProduceEmptyBlock(common.Milliseconds(610*1000), 0)
		sb = b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		assert.Equal(t, 1, len(sb.Transactions))
		assert.Equal(t, types.TransactionStatusExpired, sb.Transactions[0].Status)

		b.CreateAccount(common.N("tester"), eosio, false, true)


		b.close()

		})
	}







