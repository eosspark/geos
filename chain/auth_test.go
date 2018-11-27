package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception/try"
	"testing"
)

func initializeAuthTest() (*AuthorizationManager, *BaseTester) {
	control := GetControllerInstance()
	am := newAuthorizationManager(control)
	bt := newBaseTester(control)
	return am, bt
}

func TestMissingSigs(t *testing.T) {
	a, b := initializeAuthTest()
	b.CreateAccounts([]common.AccountName{common.AccountName(common.N("alice"))}, false, false)
	b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	try.Try(func() {
		b.PushReqAuth(common.N("alice"), &[]types.PermissionLevel{{common.N("alice"), common.DefaultConfig.ActiveName}}, &[]ecc.PrivateKey{})
	})
	fmt.Println(a.GetPermission(&types.PermissionLevel{common.N("alice"), common.N("active")}))
	//trace := b.PushReqAuth2(common.N("alice"),"owner", false)
	//b.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs),0)
	//assert.Equal(t,true,b.ChainHasTransaction(&trace.ID))
	b.Control.Close()
}

func TestMissingAuths(t *testing.T) {
	fmt.Println(common.S(6138663577826885632))
	fmt.Println(common.S(3617214756542218240))
}

func TestDelegateAuth(t *testing.T) {

}

func TestCommonEmpty(t *testing.T) {
	a := types.Permission{}
	fmt.Println(a)
	fmt.Println(common.Empty(a))
}

func TestMakeAuthChecker(t *testing.T) {
	a, _ := initializeAuthTest()
	providedKeys := common.FlatSet{}
	providedPermissions := common.FlatSet{}
	pub1, _ := ecc.NewPublicKey("EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM")
	pub2, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	providedKeys.Insert(&pub1)
	providedKeys.Insert(&pub2)
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return a.GetPermission(p).Auth },
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		&providedKeys,
		&providedPermissions,
		0,
		nil)
	fmt.Println(checker)
}
