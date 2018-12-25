package chain_plugin

import (
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	math "github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"strconv"
)

type ReadOnly struct {
	db                   *chain.Controller
	abiSerializerMaxTime common.Microseconds
	shortenAbiErrors     bool
}

func NewReadOnly(db *chain.Controller, abiSerializerMaxTime common.Microseconds) *ReadOnly {
	return &ReadOnly{db: db, abiSerializerMaxTime: abiSerializerMaxTime}
}

func (ro *ReadOnly) SetShortenAbiErrors(f bool) {
	ro.shortenAbiErrors = f
}

func (ro *ReadOnly) WalkKeyValueTable(code, scope, table common.Name, f func(interface{}) bool) {
	db := ro.db.DataBase()
	tid := entity.TableIdObject{Code: common.AccountName(code),
		Scope: common.ScopeName(scope),
		Table: common.TableName(table),
	}

	err := db.Find("byCodeScopeTable", tid, &tid)
	if err == nil { //TODO: check miss or error
		idx, e := db.GetIndex("byScopePrimary", entity.KeyValueObject{})
		EosAssert(e == nil, &DatabaseException{}, e.Error())
		newTid := tid.ID + 1
		lower, e1 := idx.LowerBound(tid)
		EosAssert(e1 == nil, &DatabaseException{}, e.Error())
		upper, e2 := idx.UpperBound(newTid)
		EosAssert(e2 == nil, &DatabaseException{}, e.Error())

		//TODO lower_bound & upper_bound
		for lower != upper {
			lower.Next()
		}
	}

}

type GetInfoResult struct {
	ServerVersion            string             `json:"server_version"`
	ChainID                  common.ChainIdType `json:"chain_id"`
	HeadBlockNum             uint32             `json:"head_block_num"`
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIdType `json:"last_irreversible_block_id"`
	HeadBlockID              common.BlockIdType `json:"head_block_id"`
	HeadBlockTime            common.TimePoint   `json:"head_block_time"`
	HeadBlockProducer        common.AccountName `json:"head_block_producer"`

	VirtualBlockCPULimit uint64 `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit uint64 `json:"virtual_block_net_limit"`
	BlockCPULimit        uint64 `json:"block_cpu_limit"`
	BlockNetLimit        uint64 `json:"block_net_limit"`
	ServerVersionString  string `json:"server_version_string"`
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

type GetBlockParams struct {
	BlockNumOrID string `json:"block_num_or_id"`
}

type GetBlockResult struct {
	SignedBlock    *types.SignedBlock `json:"signed_block"`
	ID             common.BlockIdType `json:"id"`
	BlockNum       uint32             `json:"block_num"`
	RefBlockPrefix uint32             `json:"ref_block_prefix"`
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

type GetBlockHeaderStateParams struct {
	BlockNumOrID string `json:"block_num_or_id"`
}

type GetBlockHeaderStateResult = types.BlockHeaderState

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

type GetAccountParams struct {
	AccountName        common.AccountName `json:"account_name"`
	ExpectedCoreSymbol *common.Symbol     `json:"expected_core_symbol"`
}

type GetAccountResult struct {
	AccountName   common.AccountName `json:"account_name"`
	HeadBlockNum  uint32             `json:"head_block_num"`
	HeadBlockTime common.TimePoint   `json:"head_block_time"`

	Privileged     bool             `json:"privileged"`
	LastCodeUpdate common.TimePoint `json:"last_code_update"`
	Created        common.TimePoint `json:"created"`

	CoreLiquidBalance common.Asset `json:"core_liquid_balance"`

	RAMQuota  int64 `json:"ram_quota"`
	NetWeight int64 `json:"net_weight"`
	CPUWeight int64 `json:"cpu_weight"`

	NetLimit types.AccountResourceLimit `json:"net_limit"`
	CpuLimit types.AccountResourceLimit `json:"cpu_limit"`
	RAMUsage int64                      `json:"ram_usage"`

	Permissions []types.Permission `json:"permissions"`

	TotalResources         common.Variant `json:"total_resources"`
	SelfDelegatedBandwidth common.Variant `json:"self_delegated_bandwidth"`
	RefundRequest          common.Variant `json:"refund_request"`
	VoterInfo              common.Variant `json:"voter_info"`
}

func (ro *ReadOnly) GetAccount(params GetAccountParams) GetAccountResult {
	var result GetAccountResult
	result.AccountName = params.AccountName

	//d := ro.db.DataBase()
	rm := ro.db.GetMutableResourceLimitsManager()

	result.HeadBlockNum = ro.db.HeadBlockNum()
	result.HeadBlockTime = ro.db.HeadBlockTime()

	rm.GetAccountLimits(result.AccountName, &result.RAMQuota, &result.NetWeight, &result.CPUWeight)

	a := ro.db.GetAccount(result.AccountName)

	result.Privileged = a.Privileged
	result.LastCodeUpdate = a.LastCodeUpdate
	result.Created = a.CreationDate.ToTimePoint()

	grelisted := ro.db.IsResourceGreylisted(&result.AccountName)
	result.NetLimit = rm.GetAccountNetLimitEx(result.AccountName, !grelisted)
	result.CpuLimit = rm.GetAccountCpuLimitEx(result.AccountName, !grelisted)
	result.RAMUsage = rm.GetAccountRamUsage(result.AccountName)

	//TODO permissions
	/*
			  const auto& permissions = d.get_index<permission_index,by_owner>();
		      auto perm = permissions.lower_bound( boost::make_tuple( params.account_name ) );
		      while( perm != permissions.end() && perm->owner == params.account_name ) {
		         /// T0D0: lookup perm->parent name
		         name parent;

		         // Don't lookup parent if null
		         if( perm->parent._id ) {
		            const auto* p = d.find<permission_object,by_id>( perm->parent );
		            if( p ) {
		               EOS_ASSERT(perm->owner == p->owner, invalid_parent_permission, "Invalid parent permission");
		               parent = p->name;
		            }
		         }

		         result.permissions.push_back( permission{ perm->name, parent, perm->auth.to_authority() } );
		         ++perm;
		      }
	*/

	//TODO token, delegated_bandwidth, refund, vote
	/*
			  const auto& code_account = db.db().get<account_object,by_name>( config::system_account_name );

		   abi_def abi;
		   if( abi_serializer::to_abi(code_account.abi, abi) ) {
		      abi_serializer abis( abi, abi_serializer_max_time );

		      const auto token_code = N(eosio.token);

		      auto core_symbol = extract_core_symbol();

		      if (params.expected_core_symbol.valid())
		         core_symbol = *(params.expected_core_symbol);

		      const auto* t_id = d.find<chain::table_id_object, chain::by_code_scope_table>(boost::make_tuple( token_code, params.account_name, N(accounts) ));
		      if( t_id != nullptr ) {
		         const auto &idx = d.get_index<key_value_index, by_scope_primary>();
		         auto it = idx.find(boost::make_tuple( t_id->id, core_symbol.to_symbol_code() ));
		         if( it != idx.end() && it->value.size() >= sizeof(asset) ) {
		            asset bal;
		            fc::datastream<const char *> ds(it->value.data(), it->value.size());
		            fc::raw::unpack(ds, bal);

		            if( bal.get_symbol().valid() && bal.get_symbol() == core_symbol ) {
		               result.core_liquid_balance = bal;
		            }
		         }
		      }

		      t_id = d.find<chain::table_id_object, chain::by_code_scope_table>(boost::make_tuple( config::system_account_name, params.account_name, N(userres) ));
		      if (t_id != nullptr) {
		         const auto &idx = d.get_index<key_value_index, by_scope_primary>();
		         auto it = idx.find(boost::make_tuple( t_id->id, params.account_name ));
		         if ( it != idx.end() ) {
		            vector<char> data;
		            copy_inline_row(*it, data);
		            result.total_resources = abis.binary_to_variant( "user_resources", data, abi_serializer_max_time, shorten_abi_errors );
		         }
		      }

		      t_id = d.find<chain::table_id_object, chain::by_code_scope_table>(boost::make_tuple( config::system_account_name, params.account_name, N(delband) ));
		      if (t_id != nullptr) {
		         const auto &idx = d.get_index<key_value_index, by_scope_primary>();
		         auto it = idx.find(boost::make_tuple( t_id->id, params.account_name ));
		         if ( it != idx.end() ) {
		            vector<char> data;
		            copy_inline_row(*it, data);
		            result.self_delegated_bandwidth = abis.binary_to_variant( "delegated_bandwidth", data, abi_serializer_max_time, shorten_abi_errors );
		         }
		      }

		      t_id = d.find<chain::table_id_object, chain::by_code_scope_table>(boost::make_tuple( config::system_account_name, params.account_name, N(refunds) ));
		      if (t_id != nullptr) {
		         const auto &idx = d.get_index<key_value_index, by_scope_primary>();
		         auto it = idx.find(boost::make_tuple( t_id->id, params.account_name ));
		         if ( it != idx.end() ) {
		            vector<char> data;
		            copy_inline_row(*it, data);
		            result.refund_request = abis.binary_to_variant( "refund_request", data, abi_serializer_max_time, shorten_abi_errors );
		         }
		      }

		      t_id = d.find<chain::table_id_object, chain::by_code_scope_table>(boost::make_tuple( config::system_account_name, config::system_account_name, N(voters) ));
		      if (t_id != nullptr) {
		         const auto &idx = d.get_index<key_value_index, by_scope_primary>();
		         auto it = idx.find(boost::make_tuple( t_id->id, params.account_name ));
		         if ( it != idx.end() ) {
		            vector<char> data;
		            copy_inline_row(*it, data);
		            result.voter_info = abis.binary_to_variant( "voter_info", data, abi_serializer_max_time, shorten_abi_errors );
		         }
		      }
		   }
	*/

	return result
}

type GetAbiParams struct {
	AccountName common.Name `json:"account_name"`
}

type GetAbiResult struct {
	AccountName common.Name           `json:"account_name"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

func (ro *ReadOnly) GetAbi(params GetAbiParams) GetAbiResult {
	result := GetAbiResult{}
	result.AccountName = params.AccountName

	d := ro.db.DataBase()

	account := entity.AccountObject{Name: params.AccountName}
	err := d.Find("byName", account, &account)
	EosAssert(err == nil, &DatabaseException{}, err.Error())

	var abi abi_serializer.AbiDef
	if abi_serializer.ToABI(account.Abi, &abi) {
		result.Abi = abi
	}

	return result
}

type GetCodeParams struct {
	AccountName common.Name `json:"account_name"`
	CodeAsWasm  bool        `json:"code_as_wasm"`
}

type GetCodeResult struct {
	AccountName common.Name           `json:"account_name"`
	Wast        string                `json:"wast"`
	Wasm        string                `json:"wasm"`
	CodeHash    crypto.Sha256         `json:"code_hash"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

func (ro *ReadOnly) GetCode(params GetCodeParams) GetCodeResult {
	result := GetCodeResult{AccountName: params.AccountName}
	d := ro.db.DataBase()

	account := entity.AccountObject{Name: params.AccountName}
	err := d.Find("byName", account, &account)
	EosAssert(err == nil, &DatabaseException{}, err.Error())
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

type GetRequiredKeysParams struct {
	Transaction   common.Variant  `json:"transaction"`
	AvailableKeys []ecc.PublicKey `json:"available_keys"`
}
type GetRequiredKeysResult struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

func (ro *ReadOnly) GetRequiredKeys(params GetRequiredKeysParams) GetRequiredKeysResult {
	trx := &types.Transaction{}
	EosAssert(common.FromVariant(&params.Transaction, trx) == nil, &TransactionTypeException{}, "Invalid transaction")

	candidateKeys := treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	for _, key := range params.AvailableKeys {
		candidateKeys.Add(key)
	}

	keys := ro.db.GetAuthorizationManager().GetRequiredKeys(trx, candidateKeys, 0)
	result := make([]ecc.PublicKey, 0, keys.Size())
	keys.Each(func(index int, value interface{}) {
		result = append(result, value.(ecc.PublicKey))
	})

	return GetRequiredKeysResult{RequiredKeys: result}
}

//rekey = {"available_keys":[],"transaction":{"expiration":19991,"ref_block_num":90,"ref_block_prefix":888,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":899,"context_free_actions":"hello","actions":null,"transaction_extensions":null,"signatures":[],"context_free_data":[]}}
type GetCurrencyBalanceParams struct {
	Code    common.Name
	Account common.Name
	Symbol  string
}

type GetCurrencyBalanceResult = []common.Asset

func (ro *ReadOnly) GetCurrencyBalance(params GetCurrencyBalanceParams) GetCurrencyBalanceResult {
	return GetCurrencyBalanceResult{} //TODO: get_currency_balance_result
}

type GetProducerScheduleResult struct {
	Active   common.Variant
	Pending  common.Variant
	Proposed common.Variant
}

func (ro *ReadOnly) GetProducerSchedule() GetProducerScheduleResult {
	result := GetProducerScheduleResult{}

	if err := common.ToVariant(ro.db.ActiveProducers(), &result.Active); err != nil {
		EosThrow(&ParseErrorException{}, err.Error())
	}

	if len(ro.db.PendingProducers().Producers) > 0 {
		if err := common.ToVariant(ro.db.PendingProducers(), &result.Pending); err != nil {
			EosThrow(&ParseErrorException{}, err.Error())
		}
	}

	if proposed := ro.db.ProposedProducers(); !common.Empty(proposed) && len(proposed.Producers) > 0 {
		if err := common.ToVariant(&proposed, &result.Proposed); err != nil {
			EosThrow(&ParseErrorException{}, err.Error())
		}
	}

	return result
}
