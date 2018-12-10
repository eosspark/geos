package http_plugin

const (
	chainUrl = "http://127.0.0.1:8888"
	// walletUrl = "http://127.0.0.1:8900"
	walletUrl = "http://127.0.0.1:8000"
	// walletUrl = "http://127.0.0.1:8765"
)

type Variants map[string]interface{}

const (
	chainFuncBase           string = "/v1/chain"
	getInfoFunc             string = chainFuncBase + "/get_info"
	pushTxnFunc             string = chainFuncBase + "/push_transaction"
	pushTxnsFunc            string = chainFuncBase + "/push_transactions"
	jsonToBinFunc           string = chainFuncBase + "/abi_json_to_bin"
	getBlockFunc            string = chainFuncBase + "/get_block"
	getBlockHeaderStateFunc string = chainFuncBase + "/get_block_header_state"
	getAccountFunc          string = chainFuncBase + "/get_account"
	getTableFunc            string = chainFuncBase + "/get_table_rows"
	getCodeFunc             string = chainFuncBase + "/get_code"
	getAbiFunc              string = chainFuncBase + "/get_abi"
	getRawCodeAndAbiFunc    string = chainFuncBase + "/get_raw_code_and_abi"
	getCurrencyBalanceFunc  string = chainFuncBase + "/get_currency_balance"
	getCurrencyStatsFunc    string = chainFuncBase + "/get_currency_stats"
	getProducersFunc        string = chainFuncBase + "/get_producers"
	getScheduleFunc         string = chainFuncBase + "/get_producer_schedule"
	getRequiredKeys         string = chainFuncBase + "/get_required_keys"

	historyFuncBase           string = "/v1/history"
	getActionsFunc            string = historyFuncBase + "/get_actions"
	getTransactionFunc        string = historyFuncBase + "/get_transaction"
	getKeyAccountsFunc        string = historyFuncBase + "/get_key_accounts"
	getControlledAccountsFunc string = historyFuncBase + "/get_controlled_accounts"

	accountHistoryFuncBase string = "/v1/account_history"
	getTransactionsFunc    string = accountHistoryFuncBase + "/get_transactions"

	netFuncBase    string = "/v1/net"
	netConnect     string = netFuncBase + "/connect"
	netDisconnect  string = netFuncBase + "/disconnect"
	netStatus      string = netFuncBase + "/status"
	netConnections string = netFuncBase + "/connections"

	walletFuncBase   string = "/v1/wallet"
	walletCreate     string = walletFuncBase + "/create"
	walletOpen       string = walletFuncBase + "/open"
	walletList       string = walletFuncBase + "/list_wallets"
	walletListKeys   string = walletFuncBase + "/list_keys"
	walletPublicKeys string = walletFuncBase + "/get_public_keys"
	walletLock       string = walletFuncBase + "/lock"
	walletLockAll    string = walletFuncBase + "/lock_all"
	walletUnlock     string = walletFuncBase + "/unlock"
	walletImportKey  string = walletFuncBase + "/import_key"
	walletRemoveKey  string = walletFuncBase + "/remove_key"
	walletCreateKey  string = walletFuncBase + "/create_key"
	walletSignTrx    string = walletFuncBase + "/sign_transaction"

	// keosdStop string = "/v1/keosd/stop"
)
