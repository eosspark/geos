package unittests

import (
	"fmt"
	"math"
	"testing"

	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/entity"

	"github.com/stretchr/testify/assert"
)

func TestAccountsExists(t *testing.T) {
	test := newValidatingTester(true, chain.SPECULATIVE)
	chain1_db := test.Control.DB

	tmp := entity.AccountObject{}
	nobody := entity.AccountObject{Name: common.DefaultConfig.NullAccountName}
	err := chain1_db.Find("byName", nobody, &tmp)
	assert.NotNil(t, tmp)

	activeTmp := entity.PermissionObject{Owner: common.DefaultConfig.NullAccountName, Name: common.DefaultConfig.ActiveName}
	NobodyActiveAuthority := entity.PermissionObject{}
	err = chain1_db.Find("byOwner", activeTmp, &NobodyActiveAuthority)
	assert.Equal(t, int(NobodyActiveAuthority.Auth.Threshold), 1)
	assert.Equal(t, len(NobodyActiveAuthority.Auth.Accounts), 0)
	assert.Equal(t, len(NobodyActiveAuthority.Auth.Keys), 0)

	ownerTmp := entity.PermissionObject{Owner: common.DefaultConfig.NullAccountName, Name: common.DefaultConfig.OwnerName}
	NobodyOwnerAuthority := entity.PermissionObject{}
	err = chain1_db.Find("byOwner", ownerTmp, &NobodyOwnerAuthority)
	assert.Equal(t, int(NobodyOwnerAuthority.Auth.Threshold), 1)
	assert.Equal(t, len(NobodyOwnerAuthority.Auth.Accounts), 0)
	assert.Equal(t, len(NobodyOwnerAuthority.Auth.Keys), 0)

	producerTmp := entity.AccountObject{Name: common.DefaultConfig.ProducersAccountName}
	producers := entity.AccountObject{}
	err = chain1_db.Find("byName", producerTmp, &producers)
	if err != nil {
		fmt.Println(err)
	}
	assert.NotNil(t, producers)

	ActiveProducers := test.Control.HeadBlockState().ActiveSchedule
	proActAuthTmp := entity.PermissionObject{Owner: common.DefaultConfig.ProducersAccountName, Name: common.DefaultConfig.ActiveName}
	producersActiveAuthority := entity.PermissionObject{}
	err = chain1_db.Find("byOwner", proActAuthTmp, &producersActiveAuthority)
	expectedThreshold := (len(ActiveProducers.Producers)*2)/3 + 1
	assert.Equal(t, int(producersActiveAuthority.Auth.Threshold), expectedThreshold)
	assert.Equal(t, len(producersActiveAuthority.Auth.Accounts), len(ActiveProducers.Producers))
	assert.Equal(t, len(producersActiveAuthority.Auth.Keys), 0)

	activeAuth := make([]common.AccountName, 0)
	for _, v := range producersActiveAuthority.Auth.Accounts {
		activeAuth = append(activeAuth, v.Permission.Actor)
	}

	diff := make([]common.AccountName, 0)
	for i := 0; i < int(math.Max(float64(len(activeAuth)), float64(len(ActiveProducers.Producers)))); i++ {
		var n1 common.AccountName
		if i < len(activeAuth) {
			n1 = activeAuth[i]
		} else {
			n1 = common.AccountName(0)
		}
		var n2 common.AccountName
		if i < len(ActiveProducers.Producers) {
			n2 = ActiveProducers.Producers[i].ProducerName
		} else {
			n2 = common.AccountName(0)
		}
		if n1 != n2 {
			diff = append(diff, common.AccountName(uint64(n2)-uint64(n1)))
		}
	}
	assert.Equal(t, len(diff), 0)

	proOwnAuthTmp := entity.PermissionObject{Owner: common.DefaultConfig.ProducersAccountName, Name: common.DefaultConfig.OwnerName}
	producersOwnerAuthority := entity.PermissionObject{}
	err = chain1_db.Find("byOwner", proOwnAuthTmp, &producersOwnerAuthority)
	assert.Equal(t, int(producersOwnerAuthority.Auth.Threshold), 1)
	assert.Equal(t, len(producersOwnerAuthority.Auth.Accounts), 0)
	assert.Equal(t, len(producersOwnerAuthority.Auth.Keys), 0)

	//TODO: Add checks on the other permissions of the producers account
	test.close()
}
