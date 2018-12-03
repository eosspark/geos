// package web3ext contains geth specific web3.js extensions.
package web3ext

var Modules = map[string]string{
	"rpc":    RPC_JS,
	"api":    API_JS,
	"net":    NET_JS,
	"chain":  CHAIN_JS,
	"wallet": WALLET_JS,
}

const RPC_JS = `
web3._extend({
	property: 'rpc',
	properties: [
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const API_JS = `
web3._extend({
     property:'api',
     methods:[
     new web3._extend.Method({
         name:'createKey',
         call:'api_createKey',
     }),
     new web3._extend.Method({
         name:'forking',
         call:'api_forking',
         params:1,
         inputFormatter:[null],
     }),
     ],
     properties:[]
});
`

//chainFuncBase           string = "/v1/chain"
//getInfoFunc             string = chainFuncBase + "/get_info"
//pushTxnFunc             string = chainFuncBase + "/push_transaction"
//pushTxnsFunc            string = chainFuncBase + "/push_transactions"
//jsonToBinFunc           string = chainFuncBase + "/abi_json_to_bin"
//getBlockFunc            string = chainFuncBase + "/get_block"
//getBlockHeaderStateFunc string = chainFuncBase + "/get_block_header_state"
//getAccountFunc          string = chainFuncBase + "/get_account"
//getTableFunc            string = chainFuncBase + "/get_table_rows"
//getCodeFunc             string = chainFuncBase + "/get_code"
//getAbiFunc              string = chainFuncBase + "/get_abi"
//getRawCodeAndAbiFunc    string = chainFuncBase + "/get_raw_code_and_abi"
//getCurrencyBalanceFunc  string = chainFuncBase + "/get_currency_balance"
//getCurrencyStatsFunc    string = chainFuncBase + "/get_currency_stats"
//getProducersFunc        string = chainFuncBase + "/get_producers"
//getScheduleFunc         string = chainFuncBase + "/get_producer_schedule"
//getRequiredKeys         string = chainFuncBase + "/get_required_keys"

const CHAIN_JS = `
web3._extend({
     property:'chain',
     methods:[
     new web3._extend.Method({
         name:'getInfo',
         call:'chain_getInfo',
     }),
     new web3._extend.Method({
         name:'pushTransaction',
         call:'chain_pushTransaction',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'pushTransactions',
         call:'chain_pushTransactions',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'abiJsonToBin',
         call:'chain_abiJsonToBin',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getBlock',
         call:'chain_getBlock',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getBlockHeaderState',
         call:'chain_GetBlockHeaderState',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getAccount',
         call:'chain_getAccount',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getTableRows',
         call:'chain_getTableRows',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCode',
         call:'chain_getCode',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getAbi',
         call:'chain_getAbi',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getRawCodeAndAbi',
         call:'chain_getRawCodeAndAbi',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCurrencyBalance',
         call:'chain_getCurrencyBalance',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCurrencyStats',
         call:'chain_getCurrencyStats',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getProducers',
         call:'chain_getProducers',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getProducerSchedule',
         call:'chain_getProducerSchedule',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getRequiredKeys',
         call:'chain_getRequiredKeys',
         params:1,
         inputFormatter:[null],
     }),
]
});
`

//netFuncBase    string = "/v1/net"
//netConnect     string = netFuncBase + "/connect"
//netDisconnect  string = netFuncBase + "/disconnect"
//netStatus      string = netFuncBase + "/status"
//netConnections string = netFuncBase + "/connections"
const NET_JS = `
web3._extend({
     property:'net',
     methods:[
     new web3._extend.Method({
         name:'connect',
         call:'net_connect',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'disconnect',
         call:'net_disconnect',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'status',
         call:'net_status',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'connections',
         call:'net_connections',
     }),
]
});
`

//walletFuncBase   string = "/v1/wallet"
//walletCreate     string = walletFuncBase + "/create"
//walletOpen       string = walletFuncBase + "/open"
//walletList       string = walletFuncBase + "/list_wallets"
//walletListKeys   string = walletFuncBase + "/list_keys"
//walletPublicKeys string = walletFuncBase + "/get_public_keys"
//walletLock       string = walletFuncBase + "/lock"
//walletLockAll    string = walletFuncBase + "/lock_all"
//walletUnlock     string = walletFuncBase + "/unlock"
//walletImportKey  string = walletFuncBase + "/import_key"
//walletRemoveKey  string = walletFuncBase + "/remove_key"
//walletCreateKey  string = walletFuncBase + "/create_key"
//walletSignTrx    string = walletFuncBase + "/sign_transaction"
const WALLET_JS = `
web3._extend({
     property:'wallet',
     methods:[
     new web3._extend.Method({
        name:'create',
        call:'wallet_create',
        params:1,
        inputFormatter:[null],
     }),
     new web3._extend.Method({
        name:'open',
        call:'wallet_open',
        params:1,
        inputFormatter:[null],
     }),
     new web3._extend.Method({
        name:'listWallets',
        call:'wallet_listWallets',
     }),
     new web3._extend.Method({
        name:'listKeys',
        call:'wallet_listKeys',
        params:2,
        inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'getPublicKeys',
         call:'wallet_getPublicKeys',
     }),
     new web3._extend.Method({
         name:'lock',
         call:'wallet_lock',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'lockAll',
         call:'wallet_lockAllwallets',
     }),
     new web3._extend.Method({
         name:'unlock',
         call:'wallet_unlock',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'importKey',
         call:'wallet_importKey',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'removeKey',
         call:'wallet_removeKey',
         params:3,
         inputFormatter:[null,null,null],
     }),
     new web3._extend.Method({
         name:'createKey',
         call:'wallet_createKey',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'signTransaction',
         call:'wallet_signTransaction',
         params:3,
         inputFormatter:[null,null,null],
     }),
]
});
`
