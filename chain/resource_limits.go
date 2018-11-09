package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"math"
)

type ResourceLimitsManager struct {
	db database.DataBase `json:"db"`
}

func newResourceLimitsManager(control *Controller) *ResourceLimitsManager {
	rcInstance := ResourceLimitsManager{}
	rcInstance.db = control.DB
	return &rcInstance
}

func (r *ResourceLimitsManager) InitializeDatabase() {
	config := entity.NewResourceLimitsConfigObject()
	err := r.db.Insert(&config)
	if err != nil {
		log.Error("InitializeDatabase is error: %s", err)
	}

	state := entity.DefaultResourceLimitsStateObject
	state.VirtualCpuLimit = config.CpuLimitParameters.Max
	state.VirtualNetLimit = config.NetLimitParameters.Max
	err = r.db.Insert(&state)
	if err != nil {
		log.Error("InitializeDatabase is error: %s", err)
	}
}

func (r *ResourceLimitsManager) InitializeAccount(account common.AccountName) {
	bl := entity.NewResourceLimitsObject()
	bl.Owner = account
	err := r.db.Insert(&bl)
	if err != nil {
		log.Error("InitializeAccount is error: %s", err)
	}

	bu := entity.ResourceUsageObject{}
	bu.Owner = account
	err = r.db.Insert(&bu)
	if err != nil {
		log.Error("InitializeAccount is error: %s", err)
	}
}

func (r *ResourceLimitsManager) SetBlockParameters(cpuLimitParameters types.ElasticLimitParameters, netLimitParameters types.ElasticLimitParameters) {
	cpuLimitParameters.Validate()
	netLimitParameters.Validate()
	config := entity.DefaultResourceLimitsConfigObject
	err := r.db.Find("id", config, &config)
	if err != nil {
		log.Error("SetBlockParameters is error: %s", err)
	}
	err = r.db.Modify(&config, func(c *entity.ResourceLimitsConfigObject) {
		c.CpuLimitParameters = cpuLimitParameters
		c.NetLimitParameters = netLimitParameters
	})
	if err != nil {
		log.Error("SetBlockParameters is error: %s", err)
	}
}

func (r *ResourceLimitsManager) UpdateAccountUsage(account *common.FlatSet, timeSlot uint32) { //待定
	config := entity.DefaultResourceLimitsConfigObject
	err := r.db.Find("id", config, &config)
	if err != nil {
		log.Error("UpdateAccountUsage is error: %s", err)
	}
	usage := entity.ResourceUsageObject{}
	for _, a := range account.Data {
		usage.Owner = *a.(*common.AccountName)
		err := r.db.Find("byOwner", usage, &usage)
		if err != nil {
			log.Error("UpdateAccountUsage is error: %s", err)
		}
		err = r.db.Modify(&usage, func(bu *entity.ResourceUsageObject) {
			bu.NetUsage.Add(0, timeSlot, config.AccountNetUsageAverageWindow)
			bu.CpuUsage.Add(0, timeSlot, config.AccountCpuUsageAverageWindow)
		})
		if err != nil {
			log.Error("UpdateAccountUsage is error: %s", err)
		}
	}
}

func (r *ResourceLimitsManager) AddTransactionUsage(account *common.FlatSet, cpuUsage uint64, netUsage uint64, timeSlot uint32) {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("AddTransactionUsage is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("AddTransactionUsage is error: %s", err)
	}
	for _, a := range account.Data {
		usage := entity.ResourceUsageObject{}
		usage.Owner = *a.(*common.AccountName)
		err := r.db.Find("byOwner", usage, &usage)
		if err != nil {
			log.Error("AddTransactionUsage is error: %s", err)
		}
		var unUsed, netWeight, cpuWeight int64
		r.GetAccountLimits(*a.(*common.AccountName), &unUsed, &netWeight, &cpuWeight)
		err = r.db.Modify(&usage, func(bu *entity.ResourceUsageObject) {
			bu.CpuUsage.Add(cpuUsage, timeSlot, config.AccountCpuUsageAverageWindow)
			bu.NetUsage.Add(netUsage, timeSlot, config.AccountNetUsageAverageWindow)
		})
		if err != nil {
			log.Error("AddTransactionUsage is error: %s", err)
		}

		if cpuWeight >= 0 && state.TotalCpuWeight > 0 {
			windowSize := uint64(config.AccountCpuUsageAverageWindow)
			virtualNetworkCapacityInWindow := arithmeticTypes.MulUint64(state.VirtualCpuLimit, windowSize)
			cpuUsedInWindow := arithmeticTypes.MulUint64(usage.CpuUsage.ValueEx, windowSize)
			cpuUsedInWindow, _ = cpuUsedInWindow.Div(arithmeticTypes.Uint128{Low: uint64(common.DefaultConfig.RateLimitingPrecision)})
			userWeight := arithmeticTypes.Uint128{Low: uint64(cpuWeight)}
			allUserWeight := arithmeticTypes.Uint128{Low: state.TotalCpuWeight}

			maxUserUseInWindow := virtualNetworkCapacityInWindow.Mul(userWeight)
			maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)
			EosAssert(cpuUsedInWindow.Compare(maxUserUseInWindow) < 1, &TxCpuUsageExceeded{},
				"authorizing account %s has insufficient cpu resources for this transaction,\n cpu_used_in_window: %s,\n max_user_use_in_window: %s",
				a, cpuUsedInWindow, maxUserUseInWindow)
		}

		if netWeight >= 0 && state.TotalNetWeight > 0 {
			windowSize := uint64(config.AccountNetUsageAverageWindow)
			virtualNetworkCapacityInWindow := arithmeticTypes.MulUint64(state.VirtualNetLimit, windowSize)
			netUsedInWindow := arithmeticTypes.MulUint64(usage.NetUsage.ValueEx, windowSize)
			netUsedInWindow, _ = netUsedInWindow.Div(arithmeticTypes.Uint128{Low: uint64(common.DefaultConfig.RateLimitingPrecision)})
			userWeight := arithmeticTypes.Uint128{Low: uint64(netWeight)}
			allUserWeight := arithmeticTypes.Uint128{Low: state.TotalNetWeight}

			maxUserUseInWindow := virtualNetworkCapacityInWindow.Mul(userWeight)
			maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)
			EosAssert(netUsedInWindow.Compare(maxUserUseInWindow) < 1, &TxNetUsageExceeded{},
				"authorizing account %s has insufficient net resources for this transaction,\n net_used_in_window: %s,\n max_user_use_in_window: %s",
				a, netUsedInWindow, maxUserUseInWindow)
		}
	}

	err = r.db.Modify(&state, func(rls *entity.ResourceLimitsStateObject) {
		rls.PendingCpuUsage += cpuUsage
		rls.PendingNetUsage += netUsage
	})
	if err != nil {
		log.Error("AddTransactionUsage is error: %s", err)
	}

	EosAssert(state.PendingCpuUsage <= config.CpuLimitParameters.Max, &BlockResourceExhausted{}, "Block has insufficient cpu resources")
	EosAssert(state.PendingNetUsage <= config.NetLimitParameters.Max, &BlockResourceExhausted{}, "Block has insufficient net resources")
}

func (r *ResourceLimitsManager) AddPendingRamUsage(account common.AccountName, ramDelta int64) {
	if ramDelta == 0 {
		return
	}

	usage := entity.ResourceUsageObject{}
	usage.Owner = account
	err := r.db.Find("byOwner", usage, &usage)
	if err != nil {
		log.Error("AddPendingRamUsage is error: %s", err)
	}

	EosAssert(ramDelta <= 0 || math.MaxUint64-usage.RamUsage >= uint64(ramDelta), &TransactionException{},
		"Ram usage delta would overflow UINT64_MAX")
	EosAssert(ramDelta >= 0 || usage.RamUsage >= uint64(-ramDelta), &TransactionException{},
		"Ram usage delta would underflow UINT64_MAX")

	err = r.db.Modify(&usage, func(u *entity.ResourceUsageObject) {
		u.RamUsage += uint64(ramDelta)
	})
	if err != nil {
		log.Error("AddPendingRamUsage is error: %s", err)
	}
}

func (r *ResourceLimitsManager) VerifyAccountRamUsage(account common.AccountName) {
	var ramBytes, netWeight, cpuWeight int64
	r.GetAccountLimits(account, &ramBytes, &netWeight, &cpuWeight)
	usage := entity.ResourceUsageObject{}
	usage.Owner = account
	err := r.db.Find("byOwner", usage, &usage)
	if err != nil {
		log.Error("VerifyAccountRamUsage is error: %s", err)
	}

	if ramBytes >= 0 {
		EosAssert(usage.RamUsage <= uint64(ramBytes), &RamUsageExceeded{},
			"account %s has insufficient ram; needs %d bytes has %d bytes", account, usage.RamUsage, ramBytes)
	}
}

func (r *ResourceLimitsManager) GetAccountRamUsage(account common.AccountName) int64 {
	usage := entity.ResourceUsageObject{}
	usage.Owner = account
	err := r.db.Find("byOwner", usage, &usage)
	if err != nil {
		log.Error("GetAccountRamUsage is error: %s", err)
	}
	return int64(usage.RamUsage)
}

func (r *ResourceLimitsManager) SetAccountLimits(account common.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) bool {

	findOrCreatePendingLimits := func() entity.ResourceLimitsObject {
		pendingLimits := entity.ResourceLimitsObject{}
		pendingLimits.Owner = account
		pendingLimits.Pending = true
		err := r.db.Find("byOwner", pendingLimits, &pendingLimits)
		if err != nil {
			limits := entity.ResourceLimitsObject{}
			limits.Owner = account
			limits.Pending = false
			r.db.Find("byOwner", limits, &limits)
			pendingLimits.Owner = limits.Owner
			pendingLimits.RamBytes = limits.RamBytes
			pendingLimits.NetWeight = limits.NetWeight
			pendingLimits.CpuWeight = limits.CpuWeight
			pendingLimits.Pending = true
			r.db.Insert(&pendingLimits)
			return pendingLimits
		} else {
			return pendingLimits
		}
	}

	limits := findOrCreatePendingLimits()
	decreasedLimit := false
	if ramBytes >= 0 {
		decreasedLimit = limits.RamBytes < 0 || ramBytes < limits.RamBytes
	}

	err := r.db.Modify(&limits, func(pendingLimits *entity.ResourceLimitsObject) {
		pendingLimits.RamBytes = ramBytes
		pendingLimits.NetWeight = netWeight
		pendingLimits.CpuWeight = cpuWeight
	})
	if err != nil {
		log.Error("SetAccountLimits is error: %s", err)
	}
	return decreasedLimit
}

func (r *ResourceLimitsManager) GetAccountLimits(account common.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) {
	pendingBuo := entity.ResourceLimitsObject{}
	pendingBuo.Owner = account
	pendingBuo.Pending = true
	err := r.db.Find("byOwner", pendingBuo, &pendingBuo)
	if err == nil {
		*ramBytes = pendingBuo.RamBytes
		*netWeight = pendingBuo.NetWeight
		*cpuWeight = pendingBuo.CpuWeight
	} else {
		buo := entity.ResourceLimitsObject{}
		buo.Owner = account
		buo.Pending = false
		r.db.Find("byOwner", buo, &buo)
		*ramBytes = buo.RamBytes
		*netWeight = buo.NetWeight
		*cpuWeight = buo.CpuWeight
	}
}

func (r *ResourceLimitsManager) ProcessAccountLimitUpdates() {

	byOwnerIndex, err := r.db.GetIndex("byOwner", entity.ResourceLimitsObject{})
	if err != nil {
		log.Error("ProcessAccountLimitUpdates is error: %s", err)
	}

	updateStateAndValue := func(total *uint64, value *int64, pendingValue int64, debugWhich string) {
		if *value > 0 {
			EosAssert(*total >= uint64(*value), &RateLimitingStateInconsistent{}, "underflow when reverting old value to %s", debugWhich)
			*total -= uint64(*value)
		}

		if pendingValue > 0 {
			EosAssert(math.MaxUint64-*total >= uint64(pendingValue), &RateLimitingStateInconsistent{}, "overflow when applying new value to %s", debugWhich)
			*total += uint64(pendingValue)
		}

		*value = pendingValue
	}
	state := entity.DefaultResourceLimitsStateObject
	err = r.db.Find("id", state, &state)
	if err != nil {
		log.Error("ProcessAccountLimitUpdates is error: %s", err)
	}
	err = r.db.Modify(&state, func(rso *entity.ResourceLimitsStateObject) {
		limit := entity.ResourceLimitsObject{}
		for !byOwnerIndex.Empty() {
			itr, err := byOwnerIndex.LowerBound(entity.ResourceLimitsObject{Pending: true})
			if err != nil {
				break
			}
			itr.Data(&limit)
			if byOwnerIndex.CompareEnd(itr) || limit.Pending != true {
				break
			}

			actualEntry := entity.ResourceLimitsObject{}
			actualEntry.Pending = false
			actualEntry.Owner = limit.Owner
			err = r.db.Find("byOwner", actualEntry, &actualEntry)
			if err != nil {
				log.Error("ProcessAccountLimitUpdates is error: %s", err)
			}
			err = r.db.Modify(&actualEntry, func(rlo *entity.ResourceLimitsObject) {
				updateStateAndValue(&rso.TotalRamBytes, &rlo.RamBytes, limit.RamBytes, "ram_bytes")
				updateStateAndValue(&rso.TotalCpuWeight, &rlo.CpuWeight, limit.CpuWeight, "cpu_weight")
				updateStateAndValue(&rso.TotalNetWeight, &rlo.NetWeight, limit.NetWeight, "net_weight")
			})
			if err != nil {
				log.Error("ProcessAccountLimitUpdates is error: %s", err)
			}
			err = r.db.Remove(&limit)
			if err != nil {
				log.Error("ProcessAccountLimitUpdates is error: %s", err)
			}
			itr.Release()
		}
	})
	if err != nil {
		log.Error("ProcessAccountLimitUpdates is error: %s", err)
	}
}

func (r *ResourceLimitsManager) ProcessBlockUsage(blockNum uint32) {
	s := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", s, &s)
	if err != nil {
		log.Error("ProcessBlockUsage is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("ProcessBlockUsage is error: %s", err)
	}
	err = r.db.Modify(&s, func(state *entity.ResourceLimitsStateObject) {

		state.AverageBlockCpuUsage.Add(state.PendingCpuUsage, blockNum, config.CpuLimitParameters.Periods)
		state.UpdateVirtualCpuLimit(config)
		state.PendingCpuUsage = 0

		state.AverageBlockNetUsage.Add(state.PendingNetUsage, blockNum, config.NetLimitParameters.Periods)
		state.UpdateVirtualNetLimit(config)
		state.PendingNetUsage = 0
	})
	if err != nil {
		log.Error("ProcessBlockUsage is error: %s", err)
	}
}

func (r *ResourceLimitsManager) GetVirtualBlockCpuLimit() uint64 {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetVirtualBlockCpuLimit is error: %s", err)
	}
	return state.VirtualCpuLimit
}

func (r *ResourceLimitsManager) GetVirtualBlockNetLimit() uint64 {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetVirtualBlockCpuLimit is error: %s", err)
	}
	return state.VirtualNetLimit
}

func (r *ResourceLimitsManager) GetBlockCpuLimit() uint64 {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetBlockCpuLimit is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("GetBlockCpuLimit is error: %s", err)
	}
	return config.CpuLimitParameters.Max - state.PendingCpuUsage
}

func (r *ResourceLimitsManager) GetBlockNetLimit() uint64 {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetBlockNetLimit is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("GetBlockNetLimit is error: %s", err)
	}
	return config.NetLimitParameters.Max - state.PendingNetUsage
}

func (r *ResourceLimitsManager) GetAccountCpuLimit(name common.AccountName, elastic bool) int64 {
	arl := r.GetAccountCpuLimitEx(name, elastic)
	return arl.Available
}

func (r *ResourceLimitsManager) GetAccountCpuLimitEx(name common.AccountName, elastic bool) types.AccountResourceLimit {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetAccountCpuLimitEx is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("GetAccountCpuLimitEx is error: %s", err)
	}

	usage := entity.ResourceUsageObject{}
	usage.Owner = name
	err = r.db.Find("byOwner", usage, &usage)
	if err != nil {
		log.Error("GetAccountCpuLimitEx is error: %s", err)
	}
	var cpuWeight, x, y int64
	r.GetAccountLimits(name, &x, &y, &cpuWeight)

	if cpuWeight < 0 || state.TotalCpuWeight == 0 {
		return types.AccountResourceLimit{-1, -1, -1}
	}

	arl := types.AccountResourceLimit{}
	windowSize := uint64(config.AccountCpuUsageAverageWindow)
	virtualCpuCapacityInWindow := arithmeticTypes.Uint128{}
	if elastic {
		virtualCpuCapacityInWindow = arithmeticTypes.MulUint64(state.VirtualCpuLimit, windowSize)
	} else {
		virtualCpuCapacityInWindow = arithmeticTypes.MulUint64(config.CpuLimitParameters.Max, windowSize)
	}
	userWeight := arithmeticTypes.Uint128{Low: uint64(cpuWeight)}
	allUserWeight := arithmeticTypes.Uint128{Low: state.TotalCpuWeight}
	maxUserUseInWindow := virtualCpuCapacityInWindow.Mul(userWeight)
	maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)

	//cpuUsedInWindow := IntegerDivideCeilUint64(              //TODO: ValueEx * windowSize may > MaxUnit64
	//	usage.CpuUsage.ValueEx * windowSize,
	//	uint64(common.DefaultConfig.RateLimitingPrecision))
	cpuUsedInWindow := arithmeticTypes.MulUint64(usage.CpuUsage.ValueEx, windowSize)
	cpuUsedInWindow, _ = cpuUsedInWindow.Div(arithmeticTypes.Uint128{Low: uint64(common.DefaultConfig.RateLimitingPrecision)})
	if maxUserUseInWindow.Compare(cpuUsedInWindow) != 1 {
		arl.Available = 0
	} else {
		arl.Available = types.DowngradeCast(maxUserUseInWindow.Sub(cpuUsedInWindow))
	}

	arl.Used = types.DowngradeCast(cpuUsedInWindow)
	arl.Max = types.DowngradeCast(maxUserUseInWindow)
	return arl
}

func (r *ResourceLimitsManager) GetAccountNetLimit(name common.AccountName, elastic bool) int64 {
	arl := r.GetAccountNetLimitEx(name, elastic)
	return arl.Available
}

func (r *ResourceLimitsManager) GetAccountNetLimitEx(name common.AccountName, elastic bool) types.AccountResourceLimit {
	state := entity.DefaultResourceLimitsStateObject
	err := r.db.Find("id", state, &state)
	if err != nil {
		log.Error("GetAccountNetLimitEx is error: %s", err)
	}
	config := entity.DefaultResourceLimitsConfigObject
	err = r.db.Find("id", config, &config)
	if err != nil {
		log.Error("GetAccountNetLimitEx is error: %s", err)
	}

	usage := entity.ResourceUsageObject{}
	usage.Owner = name
	err = r.db.Find("byOwner", usage, &usage)
	if err != nil {
		log.Error("GetAccountNetLimitEx is error: %s", err)
	}

	var netWeight, x, y int64
	r.GetAccountLimits(name, &x, &netWeight, &y)

	if netWeight < 0 || state.TotalNetWeight == 0 {
		return types.AccountResourceLimit{-1, -1, -1}
	}

	arl := types.AccountResourceLimit{}
	windowSize := uint64(config.AccountCpuUsageAverageWindow)
	virtualNetworkCapacityInWindow := arithmeticTypes.Uint128{}
	if elastic {
		virtualNetworkCapacityInWindow = arithmeticTypes.MulUint64(state.VirtualNetLimit, windowSize)
	} else {
		virtualNetworkCapacityInWindow = arithmeticTypes.MulUint64(config.CpuLimitParameters.Max, windowSize)
	}
	userWeight := arithmeticTypes.Uint128{Low: uint64(netWeight)}
	allUserWeight := arithmeticTypes.Uint128{Low: state.TotalNetWeight}

	maxUserUseInWindow := virtualNetworkCapacityInWindow.Mul(userWeight)
	maxUserUseInWindow, _ = maxUserUseInWindow.Div(allUserWeight)
	//cpuUsedInWindow := IntegerDivideCeilUint64(              //TODO: ValueEx * windowSize may > MaxUnit64
	//	usage.CpuUsage.ValueEx * windowSize,
	//	uint64(common.DefaultConfig.RateLimitingPrecision))
	netUsedInWindow := arithmeticTypes.MulUint64(usage.NetUsage.ValueEx, windowSize)
	netUsedInWindow, _ = netUsedInWindow.Div(arithmeticTypes.Uint128{Low: uint64(common.DefaultConfig.RateLimitingPrecision)})

	if maxUserUseInWindow.Compare(netUsedInWindow) != 1 {
		arl.Available = 0
	} else {
		arl.Available = types.DowngradeCast(maxUserUseInWindow.Sub(netUsedInWindow))
	}

	arl.Used = types.DowngradeCast(netUsedInWindow)
	arl.Max = types.DowngradeCast(maxUserUseInWindow)
	return arl
}
