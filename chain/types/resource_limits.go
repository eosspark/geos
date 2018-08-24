package types

import (
	"github.com/eosspark/eos-go/base"
	"github.com/eosspark/eos-go/common"
	chainConfig "github.com/eosspark/eos-go/chain/config"
	"fmt"
	"math"
	"reflect"
	"github.com/eosspark/eos-go/db"
)

type Ratio struct {
	Numerator   uint64 `json:"numerator"`
	Denominator uint64 `json:"denominator"`
}

type ElasticLimitParameters struct {
	Target        uint64 `json:"target"`
	Max           uint64 `json:"max"`
	Periods       uint32 `json:"periods"`
	MaxMultiplier uint32 `json:"max_multiplier"`
	ContractRate  Ratio  `json:"contract_rate"`
	ExpandRate    Ratio  `json:"expand_rate"`
}

type AccountResourceLimit struct {
	Used      int64 `json:"used"`
 	Available int64 `json:"available"`
	Max       int64 `json:"max"`
}

type ResourceLimitsManager struct {
	db *eosiodb.Session `json:"db"`
}

func NewResourceLimitsManager(db *eosiodb.Session) *ResourceLimitsManager{
	return &ResourceLimitsManager{db:db}
}

func UpdateElasticLimit(currentLimit uint64, averageUsage uint64, params ElasticLimitParameters) uint64{
	result := currentLimit
	if averageUsage > params.Target {
		result = result * params.ContractRate.Numerator / params.ContractRate.Denominator
	} else {
		result = result * params.ExpandRate.Numerator / params.ExpandRate.Denominator
	}
	return base.Min(base.Min(result, params.Max), uint64(params.Max * uint64(params.MaxMultiplier)))
}

func (elp ElasticLimitParameters) Validate(){

}

func (state *ResourceLimitsStateObject) UpdateVirtualCpuLimit(cfg ResourceLimitsConfigObject){
	state.VirtualCpuLimit = UpdateElasticLimit(state.VirtualCpuLimit, state.AverageBlockCpuUsage.Average(), cfg.CpuLimitParameters)
}

func (state *ResourceLimitsStateObject) UpdateVirtualNetLimit(cfg ResourceLimitsConfigObject){
	state.VirtualNetLimit = UpdateElasticLimit(state.VirtualNetLimit, state.AverageBlockNetUsage.Average(), cfg.NetLimitParameters)
}

func (rlm *ResourceLimitsManager) AddIndices(){
	rlm.db.Insert(&ResourceLimitsObject{Id:ResourceLimits})
	rlm.db.Insert(&ResourceUsageObject{Id:ResourceUsage})
	rlm.db.Insert(&ResourceLimitsConfigObject{Id:ResourceLimitsConfig})
	rlm.db.Insert(&ResourceLimitsStateObject{Id:ResourceLimitsState})
}

func (rlm *ResourceLimitsManager) InitializeDatabase(){
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	rlm.db.Update(&state, func(data interface{}) error {
		//ref := reflect.ValueOf(data).Elem()
		//if ref.CanSet() {
		//	ref.FieldByName("VirtualCpuLimit").SetUint(config.CpuLimitParameters.Max)
		//	ref.FieldByName("VirtualNetLimit").SetUint(config.NetLimitParameters.Max)
		//} else {
		//	// log ?
		//}
		state.VirtualCpuLimit = config.CpuLimitParameters.Max
		state.VirtualNetLimit = config.NetLimitParameters.Max
		return nil

	})
}

func (rlm *ResourceLimitsManager) InitializeAccount(account common.AccountName){
	var rlo ResourceLimitsObject
	rlm.db.Find("Rlo", RloIndex{ResourceLimits, 0, false}, &rlo)
	rlo.Id = ResourceLimits
	rlo.Owner = account
	rlo.Pending = false
	rlo.Rlo = RloIndex{ResourceLimits, account, false}
	rlm.db.Insert(&rlo)

	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, 0}, &ruo)
	ruo.Id = ResourceUsage
	ruo.Owner = account
	ruo.Ruo = RuoIndex{ResourceUsage, account}
	rlm.db.Insert(&ruo)
}

func (rlm *ResourceLimitsManager) SetBlockParameters(cpuLimitParameters ElasticLimitParameters, netLimitParameters ElasticLimitParameters){
	cpuLimitParameters.Validate()
	netLimitParameters.Validate()
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	rlm.db.Update(&config, func(data interface{}) error {
		//ref := reflect.ValueOf(data).Elem()
		//if ref.CanSet() {
		//	ref.FieldByName("CpuLimitParameters").Set(reflect.ValueOf(cpuLimitParameters))
		//	ref.FieldByName("NetLimitParameters").Set(reflect.ValueOf(netLimitParameters))
		//} else {
		//	// log ?
		//}
		config.CpuLimitParameters = cpuLimitParameters
		config.NetLimitParameters = netLimitParameters
		return nil
	})
}

func (rlm *ResourceLimitsManager) UpdateAccountUsage(account []common.AccountName, timeSlot uint32){     //待定
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	var ruo ResourceUsageObject
	for _, a := range account {
		rlm.db.Find("Ruo", RuoIndex{ResourceUsage, a}, &ruo)
		rlm.db.Update(&ruo, func(data interface{}) error {
			ruo.NetUsage.add(0, timeSlot, config.AccountNetUsageAverageWindow)
			ruo.CpuUsage.add(0, timeSlot, config.AccountCpuUsageAverageWindow)
			return nil
		})
	}
}

func (rlm *ResourceLimitsManager) AddTransactionUsage(account []common.AccountName, cpuUsage uint64, netUsage uint64, timeSlot uint32){
	var state ResourceLimitsStateObject
	var config ResourceLimitsConfigObject
	for _, a := range account{
		var ruo ResourceUsageObject
		rlm.db.Find("Ruo", RuoIndex{ResourceUsage, a}, &ruo)
		var unused, netWeight, cpuWeight int64
		rlm.GetAccountLimits(a, &unused, &netWeight, &cpuWeight)
		rlm.db.Update(&ruo, func(data interface{}) error {
			ruo.CpuUsage.add(netUsage, timeSlot, config.AccountNetUsageAverageWindow)
			ruo.NetUsage.add(cpuUsage, timeSlot, config.AccountCpuUsageAverageWindow)
			return nil
		})

		if cpuWeight >=0 && state.TotalCpuWeight > 0 {
			windowSize := uint64(config.AccountCpuUsageAverageWindow)
			virtualNetworkCapacityInWindow := state.VirtualCpuLimit * windowSize
			cpuUsedInWindow := ruo.CpuUsage.ValueEx * windowSize / uint64(chainConfig.RateLimitingPrecision)

			userWeight := cpuWeight
			allUserWeight := state.TotalCpuWeight

			maxUserUseInWindow := virtualNetworkCapacityInWindow * uint64(userWeight) /  allUserWeight

			if cpuUsedInWindow > maxUserUseInWindow {
				fmt.Println("error")
			}
		}

		if netWeight >= 0 && state.TotalNetWeight > 0 {
			windowSize := uint64(config.AccountNetUsageAverageWindow)
			virtualNetworkCapacityInWindow := state.VirtualNetLimit * windowSize
			netUsedInWindow := ruo.NetUsage.ValueEx * windowSize / uint64(chainConfig.RateLimitingPrecision)

			userWeight := netWeight
			allUserWeight := state.TotalNetWeight

			maxUserUseInWindow := virtualNetworkCapacityInWindow * uint64(userWeight) /  allUserWeight

			if netUsedInWindow > maxUserUseInWindow {
				fmt.Println("error")
			}
		}
	}

	rlm.db.Update(&state, func(data interface{}) error {
		state.PendingCpuUsage += cpuUsage
		state.PendingNetUsage += netUsage
		return nil
	})

}

func (rlm *ResourceLimitsManager) AddPendingRamUsage(account common.AccountName, ramDelta int64){
	if ramDelta == 0 {
		return
	}

	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, account}, &ruo)

	if ramDelta > 0 && math.MaxUint64 - ruo.RamUsage < uint64(ramDelta) {
		fmt.Println("error")
	}
	if ramDelta < 0 && ruo.RamUsage < uint64(-ramDelta) {
		fmt.Println("error")
	}

	rlm.db.Update(&ruo, func(data interface{}) error {
		ruo.RamUsage += uint64(ramDelta)
		return nil
	})
}

func (rlm *ResourceLimitsManager) VerifyAccountRamUsage(account common.AccountName){
	var ramBytes, netWeight, cpuWeight int64
	rlm.GetAccountLimits(account, &ramBytes, &netWeight, &cpuWeight)
	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, account}, &ruo)

	if ramBytes >=0 {
		if int64(ruo.RamUsage) > ramBytes {
			fmt.Println("error")
		}
	}
}

func (rlm *ResourceLimitsManager) GetAccountRamUsage(account common.AccountName) int64{
	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, account}, &ruo)
	return int64(ruo.RamUsage)
}

func (rlm *ResourceLimitsManager) SetAccountLimits(account common.AccountName, ramBytes int64, netWeight int64, cpuWeight int64) bool{ //for test
	var pendingRlo ResourceLimitsObject
	err := rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, true}, &pendingRlo)
	if err != nil {
		var rlo ResourceLimitsObject
		rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, false}, &rlo)
		pendingRlo.Rlo = RloIndex{rlo.Id, rlo.Owner, true}
		pendingRlo.Id = rlo.Id
		pendingRlo.Owner = rlo.Owner
		pendingRlo.Pending = true
		pendingRlo.CpuWeight = rlo.CpuWeight
		pendingRlo.NetWeight = rlo.NetWeight
		pendingRlo.RamBytes = rlo.RamBytes
		rlm.db.Insert(&pendingRlo)
	}
	decreasedLimit := false
	if ramBytes >= 0 {
		decreasedLimit = pendingRlo.RamBytes < 0 || ramBytes < pendingRlo.RamBytes
	}

	rlm.db.Update(&pendingRlo, func(data interface{}) error {
		ref := reflect.ValueOf(data).Elem()
		if ref.CanSet(){
			ref.FieldByName("RamBytes").SetInt(ramBytes)
			ref.FieldByName("NetWeight").SetInt(netWeight)
			ref.FieldByName("CpuWeight").SetInt(cpuWeight)
		}
		return nil
	})
	return decreasedLimit
}

func (rlm *ResourceLimitsManager) GetAccountLimits(account common.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) {
	var pendingRlo ResourceLimitsObject
	err := rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, true}, &pendingRlo)
	if err == nil{
		*ramBytes = pendingRlo.RamBytes
		*netWeight = pendingRlo.NetWeight
		*cpuWeight = pendingRlo.CpuWeight
	} else {
		var rlo ResourceLimitsObject
		rlm.db.Find("Rlo", RloIndex{ResourceLimits, account, false}, &rlo)
		*ramBytes = rlo.RamBytes
		*netWeight = rlo.NetWeight
		*cpuWeight = rlo.CpuWeight
	}
}

func (rlm *ResourceLimitsManager) ProcessAccountLimitUpdates(){
	updateStateAndValue := func(total *uint64, value *int64, pendingValue int64, debugWhich string){
		if *value > 0 {
			if *total < uint64(*value) {
				fmt.Println("error")
			}
			*total -= uint64(*value)
		}

		if pendingValue > 0 {
			if math.MaxUint64 - *total < uint64(pendingValue) {
				fmt.Println("error")
			}
			*total += uint64(pendingValue)
		}

		*value = pendingValue
	}
	var pendingRlo []ResourceLimitsObject
	rlm.db.Get("Pending", true, &pendingRlo)
	var state ResourceLimitsStateObject
	rlm.db.Update(&state, func(data interface{}) error {
		for _, itr := range pendingRlo {
			var rlo ResourceLimitsObject
			rlm.db.Find("Rlo", RloIndex{ResourceLimits, itr.Owner, false}, &rlo)
			rlm.db.Update(&rlo, func(data interface{}) error {
				updateStateAndValue(&state.TotalRamBytes, &rlo.RamBytes, itr.RamBytes, "ram_bytes")
				updateStateAndValue(&state.TotalCpuWeight, &rlo.CpuWeight, itr.CpuWeight, "cpu_weight")
				updateStateAndValue(&state.TotalNetWeight, &rlo.NetWeight, itr.NetWeight, "net_weight")
				return nil
			})
		}
		return nil
	})
}

func (rlm *ResourceLimitsManager) ProcessBlockUsage(blockNum uint32){
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	rlm.db.Update(&state, func(data interface{}) error {

		state.AverageBlockCpuUsage.add(state.PendingCpuUsage, blockNum, config.CpuLimitParameters.Periods)
		state.UpdateVirtualCpuLimit(config)
		state.PendingCpuUsage = 0

		state.AverageBlockNetUsage.add(state.PendingNetUsage, blockNum, config.NetLimitParameters.Periods)
		state.UpdateVirtualNetLimit(config)
		state.PendingNetUsage = 0

		return nil
	})
}

func (rlm *ResourceLimitsManager) GetVirtualBlockCpuLimit() uint64{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	return state.VirtualCpuLimit
}

func (rlm *ResourceLimitsManager) GetVirtualBlockNetLimit() uint64{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	return state.VirtualNetLimit
}

func (rlm *ResourceLimitsManager) GetBlockCpuLimit() uint64{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	return config.CpuLimitParameters.Max - state.PendingCpuUsage
}

func (rlm *ResourceLimitsManager) GetBlockNetLimit() uint64{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	return config.NetLimitParameters.Max - state.PendingNetUsage
}

func (rlm *ResourceLimitsManager) GetAccountCpuLimit(name common.AccountName, elastic bool) int64{
	arl := rlm.GetAccountCpuLimitEx(name, elastic)
	return arl.Available
}

func (rlm *ResourceLimitsManager) GetAccountCpuLimitEx(name common.AccountName, elastic bool) AccountResourceLimit{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, name}, &ruo)

	var cpuWeight, x, y int64
	rlm.GetAccountLimits(name, &x, &y, &cpuWeight)

	if cpuWeight < 0 || state.TotalCpuWeight == 0 {
		return AccountResourceLimit{-1, -1, -1}
	}

	var arl AccountResourceLimit
	windowSize := uint64(config.AccountCpuUsageAverageWindow)
	var virtualCpuCapacityInWindow uint64
	if elastic {
		virtualCpuCapacityInWindow = state.VirtualCpuLimit * windowSize
	} else {
		virtualCpuCapacityInWindow = config.CpuLimitParameters.Max * windowSize
	}
	userWeight := uint64(cpuWeight)
	allUserWeight := state.TotalCpuWeight

	maxUserUseInWindow := virtualCpuCapacityInWindow * userWeight / allUserWeight
	cpuUsedInWindow := IntegerDivideCeil(ruo.CpuUsage.ValueEx * windowSize, uint64(chainConfig.RateLimitingPrecision))

	if maxUserUseInWindow <= cpuUsedInWindow {
		arl.Available = 0
	} else {
		arl.Available = DowngradeCast(maxUserUseInWindow - cpuUsedInWindow)
	}

	arl.Used = DowngradeCast(cpuUsedInWindow)
	arl.Max = DowngradeCast(maxUserUseInWindow)
	return arl
}

func (rlm *ResourceLimitsManager) GetAccountNetLimit(name common.AccountName, elastic bool) int64 {
	arl := rlm.GetAccountNetLimitEx(name, elastic)
	return arl.Available
}

func (rlm *ResourceLimitsManager) GetAccountNetLimitEx(name common.AccountName, elastic bool) AccountResourceLimit{
	var state ResourceLimitsStateObject
	rlm.db.Find("Id", ResourceLimitsState, &state)
	var config ResourceLimitsConfigObject
	rlm.db.Find("Id", ResourceLimitsConfig, &config)
	var ruo ResourceUsageObject
	rlm.db.Find("Ruo", RuoIndex{ResourceUsage, name}, &ruo)

	var netWeight, x, y int64
	rlm.GetAccountLimits(name, &x, &y, &netWeight)

	if netWeight < 0 || state.TotalNetWeight == 0 {
		return AccountResourceLimit{-1, -1, -1}
	}

	var arl AccountResourceLimit
	windowSize := uint64(config.AccountNetUsageAverageWindow)
	var virtualNetCapacityInWindow uint64
	if elastic {
		virtualNetCapacityInWindow = state.VirtualNetLimit * windowSize
	} else {
		virtualNetCapacityInWindow = config.NetLimitParameters.Max * windowSize
	}
	userWeight := uint64(netWeight)
	allUserWeight := state.TotalNetWeight

	maxUserUseInWindow := virtualNetCapacityInWindow * userWeight / allUserWeight
	netUsedInWindow := IntegerDivideCeil(ruo.NetUsage.ValueEx * windowSize, uint64(chainConfig.RateLimitingPrecision))

	if maxUserUseInWindow <= netUsedInWindow {
		arl.Available = 0
	} else {
		arl.Available = DowngradeCast(maxUserUseInWindow - netUsedInWindow)
	}

	arl.Used = DowngradeCast(netUsedInWindow)
	arl.Max = DowngradeCast(maxUserUseInWindow)
	return arl
}