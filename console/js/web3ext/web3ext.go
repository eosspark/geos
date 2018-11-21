// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

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
			name: 'getInfo',
			call: 'api_getInfo'
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
         call:'net_getInfo',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'pushTransaction',
         call:'net_pushTransaction',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'pushTransactions',
         call:'net_pushTransactions',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'abiJsonToBin',
         call:'net_abiJsonToBin',
     }),

     new web3._extend.Method({
         name:'getBlock',
         call:'net_getBlock',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getBlockHeaderState',
         call:'netGetBlockHeaderState',
         params:1,
         inputFormatter:[null],
     }),
     new web3._extend.Method({
         name:'getAccount',
         call:'net_getAccount',
     }),
     new web3._extend.Method({
         name:'getTableRows',
         call:'net_getTableRows',
     }),
     new web3._extend.Method({
         name:'getCode',
         call:'net_getCode',
     }),
     new web3._extend.Method({
         name:'getAbi',
         call:'net_getAbi',
     }),
     new web3._extend.Method({
         name:'getRawCodeAndAbi',
         call:'net_getRawCodeAndAbi',
     }),
     new web3._extend.Method({
         name:'getCurrencyBalance',
         call:'net_getCurrencyBalance',
     }),
     new web3._extend.Method({
         name:'getCurrencyStats',
         call:'net_getCurrencyStats',
     }),
     new web3._extend.Method({
         name:'getProducers',
         call:'net_getProducers',
     }),
     new web3._extend.Method({
         name:'getProducerSchedule',
         call:'net_getProducerSchedule',
     }),
     new web3._extend.Method({
         name:'getRequiredKeys',
         call:'net_getRequiredKeys',
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
         call:'wallet_lockAll',
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
