package chain_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	math "github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/database"
	. "github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"strconv"
	"strings"
)

const (
	KEYi64 = "i64"
	i64    = "i64"
	i128   = "i128"
	i256   = "i256"
	f64    = "float64"
	f128   = "float128"
	sha256 = "sha256"
)

type ReadOnly struct {
	db                   *chain.Controller
	abiSerializerMaxTime common.Microseconds
	shortenAbiErrors     bool
}

func NewReadOnly(db *chain.Controller, abiSerializerMaxTime common.Microseconds) *ReadOnly {
	return &ReadOnly{db: db, abiSerializerMaxTime: abiSerializerMaxTime}
}

func GetAbi(db *chain.Controller, account common.Name) abi_serializer.AbiDef {
	d := db.DataBase()
	codeAccnt := &AccountObject{Name: account}
	err := d.Find("byName", codeAccnt, codeAccnt)
	EosAssert(err == nil, &AccountQueryException{}, "Fail to retrieve account for %s", account)
	var abi abi_serializer.AbiDef
	abi_serializer.ToABI(codeAccnt.Abi, &abi)
	return abi
}

func GetTableType(abi *abi_serializer.AbiDef, tableName common.TableName) string {
	for _, t := range abi.Tables {
		if t.Name == tableName {
			return t.IndexType
		}
	}
	EosAssert(false, &ContractTableQueryException{}, "Table %s is not specified in the ABI", tableName)
	return ""
}

func CopyInlineRow(obj *KeyValueObject, data *[]byte) {
	*data = make([]byte, obj.Value.Size())
	copy(*data, obj.Value)
}

func (ro *ReadOnly) SetShortenAbiErrors(f bool) {
	ro.shortenAbiErrors = f
}

func (ro *ReadOnly) WalkKeyValueTable(code, scope, table common.Name, f func(KeyValueObject) bool) {
	db := ro.db.DataBase()
	tid := TableIdObject{Code: code, Scope: scope, Table: table}
	err := db.Find("byCodeScopeTable", tid, &tid)
	if err == nil { //TODO: check miss or error
		idx, err := db.GetIndex("byScopePrimary", KeyValueObject{})
		Throw(err)

		nextTid := tid.ID + 1

		lower, err := idx.LowerBound(KeyValueObject{TId: tid.ID})
		Throw(err)

		upper, err := idx.UpperBound(KeyValueObject{TId: nextTid})
		Throw(err)

		for itr := lower; !idx.CompareIterator(itr, upper); itr.Next() {
			data := KeyValueObject{}
			if err := itr.Data(&data); err != nil {
				Throw(err)
			}
			if !f(data) {
				break
			}
		}
	}
}

func (ro *ReadOnly) GetInfo() *GetInfoResult {
	rm := ro.db.GetMutableResourceLimitsManager()
	return &GetInfoResult{
		ServerVersion:            strconv.FormatUint(app.App().GetVersion(), 10),
		ChainID:                  ro.db.GetChainId(),
		HeadBlockNum:             ro.db.ForkDbHeadBlockNum(),
		LastIrreversibleBlockNum: ro.db.LastIrreversibleBlockNum(),
		LastIrreversibleBlockID:  ro.db.LastIrreversibleBlockId(),
		HeadBlockID:              ro.db.ForkDbHeadBlockId(),
		HeadBlockTime:            ro.db.ForkDbHeadBlockTime(),
		HeadBlockProducer:        ro.db.ForkDbHeadBlockProducer(),
		VirtualBlockCPULimit:     rm.GetVirtualBlockCpuLimit(),
		VirtualBlockNetLimit:     rm.GetVirtualBlockNetLimit(),
		BlockCPULimit:            rm.GetBlockCpuLimit(),
		BlockNetLimit:            rm.GetBlockNetLimit(),
		ServerVersionString:      app.App().VersionString(),
	}
}

func (ro *ReadOnly) GetBlock(params GetBlockParams) *GetBlockResult {
	var block *types.SignedBlock
	var blockNum uint64

	EosAssert(len(params.BlockNumOrID) != 0 && len(params.BlockNumOrID) <= 64, &BlockIdTypeException{},
		"Invalid Block number or ID,must be greater than 0 and less than 64 characters ")

	num, ok := math.ParseUint64(params.BlockNumOrID)
	if ok {
		blockNum = num
		block = ro.db.FetchBlockByNumber(uint32(blockNum))
	} else {
		Try(func() {
			block = ro.db.FetchBlockById(*crypto.NewSha256String(params.BlockNumOrID))
		}).EosRethrowExceptions(&BlockIdTypeException{}, "Invalid block ID: %s", params.BlockNumOrID).End()
	}

	EosAssert(block != nil, &UnknownBlockException{}, "Could not find block: %s", params.BlockNumOrID)

	refBlockPrefix := uint32(block.BlockID().Hash[1])

	return &GetBlockResult{
		SignedBlock:    block,
		ID:             block.BlockID(),
		BlockNum:       block.BlockNumber(),
		RefBlockPrefix: refBlockPrefix,
	}
}

func (ro *ReadOnly) GetBlockHeaderState(params GetBlockHeaderStateParams) GetBlockHeaderStateResult {
	var blockNum uint64
	var b *types.BlockState

	num, ok := math.ParseUint64(params.BlockNumOrID)
	if ok {
		blockNum = num
		b = ro.db.FetchBlockStateByNumber(uint32(blockNum))
	} else {
		Try(func() {
			b = ro.db.FetchBlockStateById(*crypto.NewSha256String(params.BlockNumOrID))
		}).EosRethrowExceptions(&BlockIdTypeException{}, "Invalid block ID: %s", params.BlockNumOrID).End()
	}
	fmt.Println(blockNum, b)
	EosAssert(b != nil, &UnknownBlockException{}, "Could not find reversible block: %s", params.BlockNumOrID)

	return b.BlockHeaderState
}

type RamMarketExchangeState struct {
	Ignore1    common.Asset
	Ignore2    common.Asset
	Ignore3    common.Asset
	CoreSymbol common.Asset
	Ignore4    common.Asset
}

func (ro *ReadOnly) ExtractCoreSymbol() common.Symbol {
	coreSymbol := common.Symbol{} // Default to CORE_SYMBOL if the appropriate data structure cannot be found in the system contract table data

	// The following code makes assumptions about the contract deployed on eosio account (i.e. the system contract) and how it stores its data.
	d := ro.db.DataBase()
	tid := TableIdObject{Code: common.N("eosio"), Scope: common.N("eosio"), Table: common.N("rammarket")}
	err := d.Find("byCodeScopeTable", tid, &tid)
	if err == nil {
		idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
		Throw(err)

		it := KeyValueObject{TId: tid.ID, PrimaryKey: common.StringToSymbol(4, "RAMCORE")}
		err = idx.Find(it, &it)
		if err == nil {
			ramMarketExchangeState := RamMarketExchangeState{}

			err := rlp.DecodeBytes(it.Value, &ramMarketExchangeState)
			if err != nil {
				return coreSymbol
			}

			if ramMarketExchangeState.CoreSymbol.Symbol.Valid() {
				coreSymbol = ramMarketExchangeState.CoreSymbol.Symbol
			}
		}
	}
	return coreSymbol
}

func (ro *ReadOnly) GetAccount(params GetAccountParams) GetAccountResult {
	var result GetAccountResult
	result.AccountName = params.AccountName

	d := ro.db.DataBase()
	rm := ro.db.GetMutableResourceLimitsManager()

	result.HeadBlockNum = ro.db.HeadBlockNum()
	result.HeadBlockTime = ro.db.HeadBlockTime()

	rm.GetAccountLimits(result.AccountName, &result.RAMQuota, &result.NetWeight, &result.CPUWeight)

	a := ro.db.GetAccount(result.AccountName)

	result.Privileged = a.Privileged
	result.LastCodeUpdate = a.LastCodeUpdate
	result.Created = a.CreationDate.ToTimePoint()

	grelisted := ro.db.IsResourceGreylisted(result.AccountName)
	result.NetLimit = rm.GetAccountNetLimitEx(result.AccountName, !grelisted)
	result.CpuLimit = rm.GetAccountCpuLimitEx(result.AccountName, !grelisted)
	result.RAMUsage = rm.GetAccountRamUsage(result.AccountName)

	permissions, err := d.GetIndex("byOwner", PermissionObject{})
	Throw(err)
	lower, err := permissions.LowerBound(PermissionObject{Owner: params.AccountName})
	Throw(err)

	for !permissions.CompareEnd(lower) {
		perm := PermissionObject{}
		err = lower.Data(&perm)
		Throw(err)

		if perm.Owner != params.AccountName {
			break
		}
		/// TODO: lookup perm->parent name
		var parent common.Name

		// Don't lookup parent if null
		if perm.Parent > 0 {
			p := PermissionObject{ID: perm.Parent}
			err = d.Find("id", p, &p)
			if err == nil {
				EosAssert(perm.Owner == p.Owner, &InvalidParentPermission{}, "Invalid parent permission")
				parent = p.Name
			}
		}

		result.Permissions = append(result.Permissions, Permission{perm.Name, parent, perm.Auth.ToAuthority()})
		lower.Next()
	}

	codeAccount := AccountObject{Name: common.DefaultConfig.SystemAccountName}
	err = d.Find("byName", codeAccount, &codeAccount)

	var abi abi_serializer.AbiDef
	if abi_serializer.ToABI(codeAccount.Abi, &abi) {
		abis := abi_serializer.NewAbiSerializer(&abi, ro.abiSerializerMaxTime)

		tokenCode := common.N("eosio.token")

		/// core balance
		coreSymbol := ro.ExtractCoreSymbol()

		tid := TableIdObject{Code: tokenCode, Scope: params.AccountName, Table: common.N("accounts")}
		err = d.Find("byCodeScopeTable", tid, &tid)
		if err == nil {
			idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
			Throw(err)
			it := KeyValueObject{TId: tid.ID, PrimaryKey: coreSymbol.ToSymbolCode()}
			err = idx.Find(it, &it)
			if err == nil && it.Value.Size() >= common.SizeofAsset {
				bal := common.Asset{}
				err := rlp.DecodeBytes(it.Value, &bal)
				Throw(err)

				if bal.Symbol.Valid() && bal.Symbol == coreSymbol {
					result.CoreLiquidBalance = bal
				}
			}

		}

		tid = TableIdObject{Code: common.DefaultConfig.SystemAccountName, Scope: params.AccountName, Table: common.N("userres")}
		err = d.Find("byCodeScopeTable", tid, &tid)
		if err == nil {
			idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
			Throw(err)
			it := KeyValueObject{TId: tid.ID, PrimaryKey: uint64(params.AccountName)}
			err = idx.Find(it, &it)
			if err == nil {
				var data []byte
				CopyInlineRow(&it, &data)
				result.TotalResources = abis.BinaryToVariant("user_resources", data, ro.abiSerializerMaxTime, false)
			}
		}

		tid = TableIdObject{Code: common.DefaultConfig.SystemAccountName, Scope: params.AccountName, Table: common.N("delband")}
		err = d.Find("byCodeScopeTable", tid, &tid)
		if err == nil {
			idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
			Throw(err)
			it := KeyValueObject{TId: tid.ID, PrimaryKey: uint64(params.AccountName)}
			err = idx.Find(it, &it)
			if err == nil {
				var data []byte
				CopyInlineRow(&it, &data)
				result.SelfDelegatedBandwidth = abis.BinaryToVariant("delegated_bandwidth", data, ro.abiSerializerMaxTime, false)
			}
		}

		tid = TableIdObject{Code: common.DefaultConfig.SystemAccountName, Scope: params.AccountName, Table: common.N("refunds")}
		err = d.Find("byCodeScopeTable", tid, &tid)
		if err == nil {
			idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
			Throw(err)
			it := KeyValueObject{TId: tid.ID, PrimaryKey: uint64(params.AccountName)}
			err = idx.Find(it, &it)
			if err == nil {
				var data []byte
				CopyInlineRow(&it, &data)
				result.RefundRequest = abis.BinaryToVariant("refund_request", data, ro.abiSerializerMaxTime, false)
			}
		}

		tid = TableIdObject{Code: common.DefaultConfig.SystemAccountName, Scope: common.DefaultConfig.SystemAccountName, Table: common.N("voters")}
		err = d.Find("byCodeScopeTable", tid, &tid)
		if err == nil {
			idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
			Throw(err)
			it := KeyValueObject{TId: tid.ID, PrimaryKey: uint64(params.AccountName)}
			err = idx.Find(it, &it)
			if err == nil {
				var data []byte
				CopyInlineRow(&it, &data)
				result.VoterInfo = abis.BinaryToVariant("voter_info", data, ro.abiSerializerMaxTime, false)
			}
		}
	}

	return result
}

func (ro *ReadOnly) GetAbi(params GetAbiParams) GetAbiResult {
	result := GetAbiResult{}
	result.AccountName = params.AccountName

	d := ro.db.DataBase()

	account := AccountObject{Name: params.AccountName}
	if err := d.Find("byName", account, &account); err != nil {
		EosThrow(&DatabaseException{}, err.Error())
	}

	var abi abi_serializer.AbiDef
	if abi_serializer.ToABI(account.Abi, &abi) {
		result.Abi = abi
	}

	return result
}

func (ro *ReadOnly) GetCode(params GetCodeParams) GetCodeResult {
	result := GetCodeResult{AccountName: params.AccountName}
	d := ro.db.DataBase()

	account := AccountObject{Name: params.AccountName}
	if err := d.Find("byName", account, &account); err != nil {
		EosThrow(&DatabaseException{}, err.Error())
	}

	EosAssert(params.CodeAsWasm, &UnsupportedFeature{}, "Returning WAST from get_code is no longer supported")

	if account.Code.Size() > 0 {
		result.Wasm = string(account.Code)
		result.CodeHash = *crypto.Hash256(account.Code)
	}

	var abi abi_serializer.AbiDef
	if abi_serializer.ToABI(account.Abi, &abi) {
		result.Abi = abi
	}

	return result
}

func (ro *ReadOnly) GetRequiredKeys(params GetRequiredKeysParams) GetRequiredKeysResult {
	trx := &types.Transaction{}
	common.FromVariant(&params.Transaction, trx)

	return GetRequiredKeysResult{RequiredKeys: ro.db.GetAuthorizationManager().GetRequiredKeys(trx, &params.AvailableKeys, 0)}
}

func (ro *ReadOnly) GetTableIndexName(p GetTableRowsParams, primary *bool) uint64 {
	// see multi_index packing of index name
	table := p.Table
	index := uint64(table) & 0xFFFFFFFFFFFFFFF0
	EosAssert(index == uint64(table), &ContractTableQueryException{}, "Unsupported table name: %s", p.Table)

	*primary = false
	pos := uint64(0)
	if p.IndexPosition == "" || p.IndexPosition == "first" || p.IndexPosition == "primary" || p.IndexPosition == "one" {
		*primary = true
	} else if strings.HasPrefix(p.IndexPosition, "sec") || strings.HasPrefix(p.IndexPosition, "two") { // second, secondary
		// pos 0
	} else if strings.HasPrefix(p.IndexPosition, "ter") || strings.HasPrefix(p.IndexPosition, "th") { // tertiary, ternary, third, three
		pos = 1
	} else if strings.HasPrefix(p.IndexPosition, "fou") { // four, fourth
		pos = 2
	} else if strings.HasPrefix(p.IndexPosition, "fi") { // five, fifth
		pos = 3
	} else if strings.HasPrefix(p.IndexPosition, "six") { // six, sixth
		pos = 4
	} else if strings.HasPrefix(p.IndexPosition, "sev") { // seven, seventh
		pos = 5
	} else if strings.HasPrefix(p.IndexPosition, "eig") { // eight, eighth
		pos = 6
	} else if strings.HasPrefix(p.IndexPosition, "nin") { // nine, ninth
		pos = 7
	} else if strings.HasPrefix(p.IndexPosition, "ten") { // ten, tenth
		pos = 8
	} else {
		Try(func() {
			math.MustParseUint64(p.IndexPosition)
		}).Catch(func(interface{}) {
			EosAssert(false, &ContractTableQueryException{}, "Invalid index_position: %s", p.IndexPosition)
		}).End()
		if pos < 2 {
			*primary = true
			pos = 0
		} else {
			pos -= 2
		}
	}
	index |= pos & 0x000000000000000F
	return index
}

func (ro *ReadOnly) GetTableRowsEx(p GetTableRowsParams, abi *abi_serializer.AbiDef) GetTableRowsResult {
	result := GetTableRowsResult{}
	d := ro.db.DataBase()
	scope := convertToUint64(p.Scope, "scope")

	abis := abi_serializer.AbiSerializer{}
	abis.SetAbi(abi, ro.abiSerializerMaxTime)
	tid := TableIdObject{Code: p.Code, Scope: common.ScopeName(scope), Table: p.Table}
	if d.Find("byCodeScopeTable", tid, &tid) == nil {
		//TODO
		idx, err := d.GetIndex("byScopePrimary", KeyValueObject{})
		Throw(err)
		nextTid := tid.ID + 1
		lower, err := idx.LowerBound(KeyValueObject{TId: tid.ID})
		Throw(err)
		upper, err := idx.UpperBound(KeyValueObject{TId: nextTid})
		Throw(err)

		if len(p.LowerBound) > 0 {
			if p.KeyType == "name" {
				s := common.N(p.LowerBound)
				lower, err = idx.LowerBound(KeyValueObject{TId: tid.ID, PrimaryKey: uint64(s)})
				Throw(err)
			} else {
				lv := convertToUint64(p.LowerBound, "lower_bound")
				lower, err = idx.LowerBound(KeyValueObject{TId: tid.ID, PrimaryKey: lv})
				Throw(err)
			}
		}

		if len(p.UpperBound) > 0 {
			if p.KeyType == "name" {
				s := common.N(p.UpperBound)
				lower, err = idx.UpperBound(KeyValueObject{TId: tid.ID, PrimaryKey: uint64(s)})
				Throw(err)
			} else {
				uv := convertToUint64(p.UpperBound, "upper_bound")
				upper, err = idx.UpperBound(KeyValueObject{TId: tid.ID, PrimaryKey: uv})
				Throw(err)
			}
		}

		var data []byte
		end := common.Now().AddUs(common.Microseconds(1000 * 10)) /// 10ms max time

		count := uint32(0)
		itr := lower
		for ; !idx.CompareIterator(itr, upper); itr.Next() {
			obj := KeyValueObject{}
			err = itr.Data(&obj)
			Throw(err)
			CopyInlineRow(&obj, &data)

			if p.JSON {
				result.Rows = append(result.Rows, abis.BinaryToVariant(abis.GetTableType(p.Table), data, ro.abiSerializerMaxTime, false))
			} else {
				result.Rows = append(result.Rows, common.VariantsFromData(data))
			}

			if count = count + 1; count == p.Limit || common.Now() > end {
				itr.Next()
				break
			}
		}
		if !idx.CompareIterator(itr, upper) {
			result.More = true
		}
	}
	return result
}

func (ro *ReadOnly) GetTableRows(p GetTableRowsParams) GetTableRowsResult {
	abi := GetAbi(ro.db, p.Code)

	primary := false
	tableWithIndex := ro.GetTableIndexName(p, &primary)
	if primary {
		EosAssert(uint64(p.Table) == tableWithIndex, &ContractTableQueryException{}, "Invalid table name %s", p.Table)
		tableType := GetTableType(&abi, p.Table)
		if tableType == KEYi64 || p.KeyType == KEYi64 || p.KeyType == "name" {
			return ro.GetTableRowsEx(p, &abi)
		}
		EosAssert(false, &ContractTableQueryException{}, "Invalid table type %s", tableType)
	} else {
		EosAssert(len(p.KeyType) != 0, &ContractTableQueryException{}, "key type required for non-primary index")

		//todo: second key
	}

	return GetTableRowsResult{}
}

func (ro *ReadOnly) GetTableByScope(p GetTableByScopeParams) GetTableByScopeResult {
	d := ro.db.DB
	idx, err := d.GetIndex("byCodeScopeTable", TableIdObject{})
	Throw(err)
	var lower database.Iterator
	var upper database.Iterator

	if len(p.LowerBound) > 0 {
		scope := convertToUint64(p.LowerBound, "lower_bound scope")
		lower, err = idx.LowerBound(TableIdObject{Code: p.Code, Scope: common.ScopeName(scope), Table: p.Table})
		Throw(err)
	} else {
		lower, err = idx.LowerBound(TableIdObject{Code: p.Code, Scope: 0, Table: p.Table})
	}
	if len(p.UpperBound) > 0 {
		scope := convertToUint64(p.UpperBound, "upper_bound scope")
		upper, err = idx.LowerBound(TableIdObject{Code: p.Code, Scope: common.ScopeName(scope), Table: 0})
		Throw(err)
	} else {
		upper, err = idx.LowerBound(TableIdObject{Code: p.Code + 1, Scope: 0, Table: 0})
	}

	end := common.Now().AddUs(common.Microseconds(1000 * 10))
	count := uint32(0)
	itr := lower
	result := GetTableByScopeResult{}
	for ; !idx.CompareIterator(itr, upper); itr.Next() {
		obj := TableIdObject{}
		Throw(itr.Data(&obj))
		if p.Table > 0 && obj.Table != p.Table {
			if common.Now() > end {
				break
			}
			continue
		}
		result.Rows = append(result.Rows, GetTableByScopeResultRow{obj.Code, obj.Scope, obj.Table, obj.Payer, obj.Count})
		if count++; count == p.Limit || common.Now() > end {
			itr.Next()
			break
		}
	}
	if !idx.CompareIterator(itr, upper) {
		obj := TableIdObject{}
		Throw(itr.Data(&obj))
		result.More = obj.Scope.String()
	}
	return result
}

func (ro *ReadOnly) GetCurrencyBalance(params GetCurrencyBalanceParams) GetCurrencyBalanceResult {
	abi := GetAbi(ro.db, params.Code)
	GetTableType(&abi, common.N("accounts"))

	var results []common.Asset
	ro.WalkKeyValueTable(params.Code, params.Account, common.N("accounts"), func(obj KeyValueObject) bool {
		EosAssert(obj.Value.Size() >= common.SizeofAsset, &AssetTypeException{}, "Invalid data on table")

		cursor := common.Asset{}
		err := rlp.DecodeBytes(obj.Value, &cursor)
		Throw(err)

		EosAssert(cursor.Symbol.Valid(), &AssetTypeException{}, "Invalid asset")

		if params.Symbol == "" || cursor.Symbol.Name() == params.Symbol {
			results = append(results, cursor)
		}

		// return false if we are looking for one and found it, true otherwise
		return !(params.Symbol != "" && cursor.Symbol.Name() == params.Symbol)
	})

	return results
}

func (ro *ReadOnly) GetCurrencyStats(params GetCurrencyStatsParams) map[string]GetCurrencyStatsResult {
	results := make(map[string]GetCurrencyStatsResult)

	abi := GetAbi(ro.db, params.Code)
	GetTableType(&abi, common.N("stat")) //assert if error

	scope := common.StringToSymbol(0, strings.ToUpper(params.Symbol)) >> 8

	ro.WalkKeyValueTable(params.Code, common.Name(scope), common.N("stat"), func(obj KeyValueObject) bool {
		EosAssert(obj.Value.Size() >= SizeofGetCurrencyStatsResult, &AssetTypeException{}, "Invalid data on table")

		ds := rlp.NewDecoder(obj.Value)
		result := GetCurrencyStatsResult{}

		ds.Decode(&result.Supply)
		ds.Decode(&result.MaxSupply)
		ds.Decode(&result.Issuer)

		results[result.Supply.Symbol.Name()] = result
		return true
	})

	return results
}

func (ro *ReadOnly) GetProducerSchedule() GetProducerScheduleResult {
	result := GetProducerScheduleResult{}

	common.ToVariant(ro.db.ActiveProducers(), &result.Active)

	if len(ro.db.PendingProducers().Producers) > 0 {
		common.ToVariant(ro.db.PendingProducers(), &result.Pending)
	}

	if proposed := ro.db.ProposedProducers(); !common.Empty(proposed) && len(proposed.Producers) > 0 {
		common.ToVariant(&proposed, &result.Proposed)
	}

	return result
}
