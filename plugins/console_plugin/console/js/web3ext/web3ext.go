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
         call:'/v1/chain/get_info',
     }),
     new web3._extend.Method({
         name:'pushTransaction',
         call:'/v1/chain/push_transaction',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'pushTransactions',
         call:'/v1/chain/push_transactions',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'abiJsonToBin',
         call:'/v1/chain/abi_json_to_bin',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getBlock',
         call:'/v1/chain/get_block',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getBlockHeaderState',
         call:'/v1/chain/get_block_header_state',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getAccount',
         call:'/v1/chain/get_account',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getTableRows',
         call:'/v1/chain/get_table_rows',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCode',
         call:'/v1/chain/get_code',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getAbi',
         call:'/v1/chain/get_abi',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getRawCodeAndAbi',
         call:'/v1/chain/get_raw_code_and_abi',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCurrencyBalance',
         call:'/v1/chain/get_currency_balance',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getCurrencyStats',
         call:'/v1/chain/get_currency_stats',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getProducers',
         call:'/v1/chain/get_producers',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getProducerSchedule',
         call:'/v1/chain/get_producer_schedule',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getRequiredKeys',
         call:'/v1/chain/get_required_keys',
         params:1,
         inputFormatter:[null],
     }),

]
});
`

//new web3._extend.Method({
//name:'createNewAccount',
//call:'createNewAccount',
//params:4,
//inputFormatter:[null,null,null,null],
//}),

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
        call:'/v1/wallet/create',
        params:1,
        inputFormatter:[null],
     }),
     new web3._extend.Method({
        name:'open',
        call:'/v1/wallet/open',
        params:1,
        inputFormatter:[null],
     }),
     new web3._extend.Method({
        name:'listWallets',
        call:'/v1/wallet/list_wallets',
     }),
     new web3._extend.Method({
        name:'listKeys',
        call:'/v1/wallet/list_keys',
        params:2,
        inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'getPublicKeys',
         call:'/v1/wallet/get_public_keys',
     }),
     new web3._extend.Method({
         name:'lock',
         call:'/v1/wallet/lock',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'lockAll',
         call:'/v1/wallet/lock_all',
     }),
     new web3._extend.Method({
         name:'unlock',
         call:'/v1/wallet/unlock',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'importKey',
         call:'/v1/wallet/import_key',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'removeKey',
         call:'/v1/wallet/remove_key',
         params:3,
         inputFormatter:[null,null,null],
     }),
     new web3._extend.Method({
         name:'createKey',
         call:'/v1/wallet/create_key',
         params:2,
         inputFormatter:[null,null],
     }),
     new web3._extend.Method({
         name:'signTransaction',
         call:'/v1/wallet/sign_transaction',
         params:3,
         inputFormatter:[null,null,null],
     }),
]
});
`
