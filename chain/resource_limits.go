package chain

import "github.com/eosspark/eos-go/database"

var IsActiveRc bool

var rcInstance *ResourceLimitsManager

type ResourceLimitsManager struct {
	db *database.LDataBase `json:"db"`
}

func GetResourceLimitsManager() *ResourceLimitsManager {
	return &ResourceLimitsManager{}
	//if !IsActiveRc {
	//	rcInstance = newResourceLimitsManager()
	//}
	//return rcInstance
}

//func newResourceLimitsManager() *ResourceLimitsManager {
//	IsActiveRc = true
//	control := GetControllerInstance()
//	db := control.DataBase()
//	return &ResourceLimitsManager{db: db}
//}

//func (rlm *ResourceLimitsManager) InitializeDatabase() {
//	config := entity.NewResourceLimitsConfigObject()
//	rlm.db.Insert(&config)
//	state := entity.ResourceLimitsStateObject{}
//	rlm.db.Find("byId", &state)
//	rlm.db.Modify(&state, func(data interface{}) error {
//		ref := reflect.ValueOf(data).Elem()
//		if ref.CanSet() {
//			ref.FieldByName("VirtualCpuLimit").SetUint(config.CpuLimitParameters.Max)
//			ref.FieldByName("VirtualNetLimit").SetUint(config.NetLimitParameters.Max)
//		} else {
//			// log ?
//		}
//		//state.VirtualCpuLimit = config.CpuLimitParameters.Max
//		//state.VirtualNetLimit = config.NetLimitParameters.Max
//		return nil
//
//	})
//}
//
//func (rlm *ResourceLimitsManager) InitializeAccount(account common.AccountName) {
//	rlo := entity.ResourceLimitsObject{}
//	rlo.Owner = account
//	rlm.db.Insert(&rlo)
//
//	ruo := entity.ResourceUsageObject{}
//	ruo.Owner = account
//	rlm.db.Insert(&ruo)
//}
//
//func (rlm *ResourceLimitsManager) SetBlockParameters(cpuLimitParameters types.ElasticLimitParameters, netLimitParameters types.ElasticLimitParameters) {
//	cpuLimitParameters.Validate()
//	netLimitParameters.Validate()
//	config := entity.ResourceLimitsConfigObject{}
//	rlm.db.Find("byId", &config)
//	rlm.db.Modify(&config, func(data interface{}) error {
//		//ref := reflect.ValueOf(data).Elem()
//		//if ref.CanSet() {
//		//	ref.FieldByName("CpuLimitParameters").Set(reflect.ValueOf(cpuLimitParameters))
//		//	ref.FieldByName("NetLimitParameters").Set(reflect.ValueOf(netLimitParameters))
//		//} else {
//		//	// log ?
//		//}
//		config.CpuLimitParameters = cpuLimitParameters
//		config.NetLimitParameters = netLimitParameters
//		return nil
//	})
//}
//
//func (rlm *ResourceLimitsManager) UpdateAccountUsage(account []common.AccountName, timeSlot uint32) { //待定
//	config := entity.ResourceLimitsConfigObject{}
//	rlm.db.Find("byId", &config)
//	ruo := entity.ResourceUsageObject{}
//	for _, a := range account {
//		ruo.Owner = a
//		rlm.db.Find("ByOwner", &ruo)
//		rlm.db.Modify(&ruo, func(data interface{}) error {
//			ruo.NetUsage.Add(0, timeSlot, config.AccountNetUsageAverageWindow)
//			ruo.CpuUsage.Add(0, timeSlot, config.AccountCpuUsageAverageWindow)
//			return nil
//		})
//	}
//}
//
//func (rlm *ResourceLimitsManager) AddTransactionUsage(account []common.AccountName, cpuUsage uint64, netUsage uint64, timeSlot uint32) {
//	state := entity.ResourceLimitsStateObject{}
//	rlm.db.Find("byId", &state)
//	config := entity.ResourceLimitsConfigObject{}
//	rlm.db.Find("byId", &config)
//	for _, a := range account {
//		ruo := entity.ResourceUsageObject{}
//		ruo.Owner = a
//		rlm.db.Find("byOwner", &ruo)
//		var unUsed, netWeight, cpuWeight int64
//		rlm.GetAccountLimits(a, &unUsed, &netWeight, &cpuWeight)
//		rlm.db.Modify(&ruo, func(data interface{}) error {
//			ruo.CpuUsage.Add(netUsage, timeSlot, config.AccountNetUsageAverageWindow)
//			ruo.NetUsage.Add(cpuUsage, timeSlot, config.AccountCpuUsageAverageWindow)
//			return nil
//		})
//
//		if cpuWeight >= 0 && state.TotalCpuWeight > 0 {
//			windowSize := new(big.Int).SetUint64(uint64(config.AccountCpuUsageAverageWindow))
//			virtualNetworkCapacityInWindow := new(big.Int).Mul(windowSize, new(big.Int).SetUint64(state.VirtualCpuLimit))
//			cpuUsedInWindow := new(big.Int).Div(
//				new(big.Int).Mul(windowSize, new(big.Int).SetUint64(ruo.CpuUsage.ValueEx)),
//				new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
//
//			userWeight := new(big.Int).SetInt64(cpuWeight)
//			allUserWeight := new(big.Int).SetUint64(state.TotalCpuWeight)
//
//			maxUserUseInWindow := new(big.Int).Div(
//				new(big.Int).Mul(virtualNetworkCapacityInWindow, userWeight), allUserWeight)
//			if cpuUsedInWindow.Cmp(maxUserUseInWindow) == 1 {
//				fmt.Println("error")
//			}
//		}
//
//		if netWeight >= 0 && state.TotalNetWeight > 0 {
//			windowSize := new(big.Int).SetUint64(uint64(config.AccountNetUsageAverageWindow))
//			virtualNetworkCapacityInWindow := new(big.Int).Mul(windowSize, new(big.Int).SetUint64(state.VirtualNetLimit))
//			netUsedInWindow := new(big.Int).Div(
//				new(big.Int).Mul(windowSize, new(big.Int).SetUint64(ruo.NetUsage.ValueEx)),
//				new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
//
//			userWeight := new(big.Int).SetInt64(netWeight)
//			allUserWeight := new(big.Int).SetUint64(state.TotalNetWeight)
//
//			maxUserUseInWindow := new(big.Int).Div(
//				new(big.Int).Mul(virtualNetworkCapacityInWindow, userWeight), allUserWeight)
//			if netUsedInWindow.Cmp(maxUserUseInWindow) == 1 {
//				fmt.Println("error")
//			}
//		}
//	}
//
//	rlm.db.Modify(&state, func(data interface{}) error {
//		state.PendingCpuUsage += cpuUsage
//		state.PendingNetUsage += netUsage
//		return nil
//	})
//
//}
//
//func (rlm *ResourceLimitsManager) AddPendingRamUsage(account common.AccountName, ramDelta int64) {
//	if ramDelta == 0 {
//		return
//	}
//
//	ruo := entity.ResourceUsageObject{}
//	ruo.Owner = account
//	rlm.db.Find("byOwner", &ruo)
//
//	if ramDelta > 0 && math.MaxUint64-ruo.RamUsage < uint64(ramDelta) {
//		fmt.Println("error")
//	}
//	if ramDelta < 0 && ruo.RamUsage < uint64(-ramDelta) {
//		fmt.Println("error")
//	}
//
//	rlm.db.Modify(&ruo, func(data interface{}) error {
//		ruo.RamUsage += uint64(ramDelta)
//		return nil
//	})
//}
//
//func (rlm *ResourceLimitsManager) VerifyAccountRamUsage(account common.AccountName) {
//	var ramBytes, netWeight, cpuWeight int64
//	rlm.GetAccountLimits(account, &ramBytes, &netWeight, &cpuWeight)
//	ruo := entity.ResourceUsageObject{}
//
//	rlm.db.Find("byOwner", &ruo)
//
//	if ramBytes >= 0 {
//		if int64(ruo.RamUsage) > ramBytes {
//			fmt.Println("error")
//		}
//	}
//}
//
//func (rlm *ResourceLimitsManager) GetAccountRamUsage(account common.AccountName) int64 {
//	ruo := ResourceUsageObject{}
//	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, account}, &ruo)
//	return int64(ruo.RamUsage)
//}
//
//func (rlm *ResourceLimitsManager) SetAccountLimits(account common.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) bool { //for test
//	pendingRlo := ResourceLimitsObject{}
//	err := rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, true}, &pendingRlo)
//	if err != nil {
//		rlo := ResourceLimitsObject{}
//		rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, false}, &rlo)
//		pendingRlo.Rlo = RloIndex{rlo.ID, rlo.Owner, true}
//		pendingRlo.ID = rlo.ID
//		pendingRlo.Owner = rlo.Owner
//		pendingRlo.Pending = true
//		pendingRlo.CpuWeight = rlo.CpuWeight
//		pendingRlo.NetWeight = rlo.NetWeight
//		pendingRlo.RamBytes = rlo.RamBytes
//		rlm.db.Insert(&pendingRlo)
//	}
//	decreasedLimit := false
//	if ramBytes >= 0 {
//		decreasedLimit = pendingRlo.RamBytes < 0 || ramBytes < pendingRlo.RamBytes
//	}
//
//	rlm.db.Update(&pendingRlo, func(data interface{}) error {
//		ref := reflect.ValueOf(data).Elem()
//		if ref.CanSet() {
//			ref.FieldByName("RamBytes").SetInt(ramBytes)
//			ref.FieldByName("NetWeight").SetInt(netWeight)
//			ref.FieldByName("CpuWeight").SetInt(cpuWeight)
//		}
//		return nil
//	})
//	return decreasedLimit
//}
//
//func (rlm *ResourceLimitsManager) GetAccountLimits(account common.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) {
//	pendingRlo := ResourceLimitsObject{}
//	err := rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, true}, &pendingRlo)
//	if err == nil {
//		*ramBytes = pendingRlo.RamBytes
//		*netWeight = pendingRlo.NetWeight
//		*cpuWeight = pendingRlo.CpuWeight
//	} else {
//		rlo := ResourceLimitsObject{}
//		rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, false}, &rlo)
//		*ramBytes = rlo.RamBytes
//		*netWeight = rlo.NetWeight
//		*cpuWeight = rlo.CpuWeight
//	}
//}
//
//func (rlm *ResourceLimitsManager) ProcessAccountLimitUpdates() {
//	updateStateAndValue := func(total *uint64, value *int64, pendingValue int64, debugWhich string) {
//		if *value > 0 {
//			if *total < uint64(*value) {
//				fmt.Println("error")
//			}
//			*total -= uint64(*value)
//		}
//
//		if pendingValue > 0 {
//			if math.MaxUint64-*total < uint64(pendingValue) {
//				fmt.Println("error")
//			}
//			*total += uint64(pendingValue)
//		}
//
//		*value = pendingValue
//	}
//	var pendingRlo []ResourceLimitsObject
//	rlm.db.Get("Pending", true, &pendingRlo)
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	rlm.db.Update(&state, func(data interface{}) error {
//		for _, itr := range pendingRlo {
//			rlo := ResourceLimitsObject{}
//			rlm.db.Find("Rlo", RloIndex{ResourceLimits, itr.Owner, false}, &rlo)
//			rlm.db.Update(&rlo, func(data interface{}) error {
//				updateStateAndValue(&state.TotalRamBytes, &rlo.RamBytes, itr.RamBytes, "ram_bytes")
//				updateStateAndValue(&state.TotalCpuWeight, &rlo.CpuWeight, itr.CpuWeight, "cpu_weight")
//				updateStateAndValue(&state.TotalNetWeight, &rlo.NetWeight, itr.NetWeight, "net_weight")
//				return nil
//			})
//		}
//		return nil
//	})
//}
//
//func (rlm *ResourceLimitsManager) ProcessBlockUsage(blockNum uint32) {
//	config := ResourceLimitsConfigObject{}
//	rlm.db.Find("ID", ResourceLimitsConfig, &config)
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	rlm.db.Update(&state, func(data interface{}) error {
//
//		state.AverageBlockCpuUsage.add(state.PendingCpuUsage, blockNum, config.CpuLimitParameters.Periods)
//		state.UpdateVirtualCpuLimit(config)
//		state.PendingCpuUsage = 0
//
//		state.AverageBlockNetUsage.add(state.PendingNetUsage, blockNum, config.NetLimitParameters.Periods)
//		state.UpdateVirtualNetLimit(config)
//		state.PendingNetUsage = 0
//
//		return nil
//	})
//}
//
//func (rlm *ResourceLimitsManager) GetVirtualBlockCpuLimit() uint64 {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	return state.VirtualCpuLimit
//}
//
//func (rlm *ResourceLimitsManager) GetVirtualBlockNetLimit() uint64 {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	return state.VirtualNetLimit
//}
//
//func (rlm *ResourceLimitsManager) GetBlockCpuLimit() uint64 {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	config := ResourceLimitsConfigObject{}
//	rlm.db.Find("ID", ResourceLimitsConfig, &config)
//	return config.CpuLimitParameters.Max - state.PendingCpuUsage
//}
//
//func (rlm *ResourceLimitsManager) GetBlockNetLimit() uint64 {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	config := ResourceLimitsConfigObject{}
//	rlm.db.Find("ID", ResourceLimitsConfig, &config)
//	return config.NetLimitParameters.Max - state.PendingNetUsage
//}
//
//func (rlm *ResourceLimitsManager) GetAccountCpuLimit(name common.AccountName, elastic bool) int64 {
//	arl := rlm.GetAccountCpuLimitEx(name, elastic)
//	return arl.Available
//}
//
//func (rlm *ResourceLimitsManager) GetAccountCpuLimitEx(name common.AccountName, elastic bool) AccountResourceLimit {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	config := ResourceLimitsConfigObject{}
//	rlm.db.Find("ID", ResourceLimitsConfig, &config)
//	ruo := ResourceUsageObject{}
//	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, name}, &ruo)
//
//	var cpuWeight, x, y int64
//	rlm.GetAccountLimits(name, &x, &y, &cpuWeight)
//
//	if cpuWeight < 0 || state.TotalCpuWeight == 0 {
//		return AccountResourceLimit{-1, -1, -1}
//	}
//
//	arl := AccountResourceLimit{}
//	windowSize := new(big.Int).SetUint64(uint64(config.AccountCpuUsageAverageWindow))
//	virtualCpuCapacityInWindow := new(big.Int)
//	if elastic {
//		virtualCpuCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(state.VirtualCpuLimit), windowSize)
//	} else {
//		virtualCpuCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(config.CpuLimitParameters.Max), windowSize)
//	}
//	userWeight := new(big.Int).SetUint64(uint64(cpuWeight))
//	allUserWeight := new(big.Int).SetUint64(state.TotalCpuWeight)
//
//	maxUserUseInWindow := new(big.Int).Div(new(big.Int).Mul(virtualCpuCapacityInWindow, userWeight), allUserWeight)
//	cpuUsedInWindow := IntegerDivideCeil(
//		new(big.Int).Mul(new(big.Int).SetUint64(ruo.CpuUsage.ValueEx), windowSize),
//		new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
//
//	if maxUserUseInWindow.Cmp(cpuUsedInWindow) != 1 {
//		arl.Available = 0
//	} else {
//		arl.Available = DowngradeCast(new(big.Int).Sub(maxUserUseInWindow, cpuUsedInWindow))
//	}
//
//	arl.Used = DowngradeCast(cpuUsedInWindow)
//	arl.Max = DowngradeCast(maxUserUseInWindow)
//	return arl
//}
//
//func (rlm *ResourceLimitsManager) GetAccountNetLimit(name common.AccountName, elastic bool) int64 {
//	arl := rlm.GetAccountNetLimitEx(name, elastic)
//	return arl.Available
//}
//
//func (rlm *ResourceLimitsManager) GetAccountNetLimitEx(name common.AccountName, elastic bool) AccountResourceLimit {
//	state := ResourceLimitsStateObject{}
//	rlm.db.Find("ID", ResourceLimitsState, &state)
//	config := ResourceLimitsConfigObject{}
//	rlm.db.Find("ID", ResourceLimitsConfig, &config)
//	ruo := ResourceUsageObject{}
//	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, name}, &ruo)
//
//	var netWeight, x, y int64
//	rlm.GetAccountLimits(name, &x, &y, &netWeight)
//
//	if netWeight < 0 || state.TotalNetWeight == 0 {
//		return AccountResourceLimit{-1, -1, -1}
//	}
//
//	arl := AccountResourceLimit{}
//	windowSize := new(big.Int).SetUint64(uint64(config.AccountNetUsageAverageWindow))
//	virtualNetCapacityInWindow := new(big.Int)
//	if elastic {
//		virtualNetCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(state.VirtualNetLimit), windowSize)
//	} else {
//		virtualNetCapacityInWindow = new(big.Int).Mul(new(big.Int).SetUint64(config.NetLimitParameters.Max), windowSize)
//	}
//	userWeight := new(big.Int).SetUint64(uint64(netWeight))
//	allUserWeight := new(big.Int).SetUint64(state.TotalNetWeight)
//
//	maxUserUseInWindow := new(big.Int).Div(new(big.Int).Mul(virtualNetCapacityInWindow, userWeight), allUserWeight)
//	netUsedInWindow := IntegerDivideCeil(
//		new(big.Int).Mul(new(big.Int).SetUint64(ruo.NetUsage.ValueEx), windowSize),
//		new(big.Int).SetUint64(uint64(common.DefaultConfig.RateLimitingPrecision)))
//
//	if maxUserUseInWindow.Cmp(netUsedInWindow) != 1 {
//		arl.Available = 0
//	} else {
//		arl.Available = DowngradeCast(new(big.Int).Sub(maxUserUseInWindow, netUsedInWindow))
//	}
//
//	arl.Used = DowngradeCast(netUsedInWindow)
//	arl.Max = DowngradeCast(maxUserUseInWindow)
//	return arl
//}
