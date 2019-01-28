package console

import "github.com/eosspark/eos-go/common"

type ConsoleInterface interface {
	getOptions() *StandardTransactionOptions
}

type StandardTransactionOptions struct {
	Expiration        uint64   `json:"expiration"`
	TxForceUnique     bool     `json:"force_unique"`
	TxSkipSign        bool     `json:"skip_sign"`
	TxPrintJson       bool     `json:"json"`
	TxDontBroadcast   bool     `json:"dont_broadcast"`
	TxReturnPacked    bool     `json:"return_packed"`
	TxRefBlockNumOrId string   `json:"ref_block"`
	TxPermission      []string `json:"p"`
	TxMaxCpuUsage     uint8    `json:"max_cpu_usage_ms"`
	TxMaxNetUsage     uint32   `json:"max_net_usage"`
	DelaySec          uint32   `json:"delay_sec"`
}

func (s *StandardTransactionOptions) getOptions() *StandardTransactionOptions {
	return s
}

type CreateAccountParams struct {
	Creator   common.Name `json:"creator"`
	Name      common.Name `json:"name"`
	OwnerKey  string      `json:"owner"`
	ActiveKey string      `json:"active"`
	StandardTransactionOptions
}

type PushAction struct {
	ContractAccount string `json:"account"`
	Action          string `json:"action"`
	Data            string `json:"data"`
	StandardTransactionOptions
}

type SetContractParams struct {
	Account                string `json:"account"`
	ContractPath           string `json:"code_file"`
	AbiPath                string `json:"abi_file"`
	ContractClear          bool   `json:"clear"`
	SuppressDuplicateCheck bool   `json:"suppress_duplicate_check"`
	StandardTransactionOptions
}

type SetAccountPermissionParams struct {
	Account             string `json:"account"`
	Permission          string `json:"permission"`
	AuthorityJsonOrFile string `json:"authority"`
	Parent              string `json:"parent"`
	StandardTransactionOptions
}

type SetActionPermissionParams struct {
	Account     string `json:"account"`
	Code        string `json:"code"`
	TypeStr     string `json:"type"`
	Requirement string `json:"requirement"`
	StandardTransactionOptions
}

type TransferParams struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
	Memo      string `json:"memo"`
	PayRam    bool   `json:"pay_ram"`
	StandardTransactionOptions
}

//system
type NewAccountParams struct {
	Creator             common.Name `json:"creator"`
	Name                common.Name `json:"name"`
	OwnerKey            string      `json:"owner"`
	ActiveKey           string      `json:"active"`
	StakeNet            string      `json:"stake_net"`
	StakeCpu            string      `json:"stake_cpu"`
	BuyRamBytesInKbytes uint32      `json:"buy_ram_kbytes"`
	BuyRamBytes         uint32      `json:"buy_ram_bytes"`
	BuyRamEos           string      `json:"buy_ram"`
	Transfer            bool        `json:"transfer"`
	StandardTransactionOptions
}

type RegisterProducer struct {
	Producer string `json:"producer"`
	Key      string `json:"key"`
	Url      string `json:"url"`
	Loc      uint16 `json:"loc"`
	StandardTransactionOptions
}

type UnregrodParams struct {
	Producer string `json:"producer"`
	StandardTransactionOptions
}

type Proxy struct {
	Voter string `json:"voter"`
	Proxy string `json:"proxy"`
	StandardTransactionOptions
}

type Prods struct {
	Voter         string `json:"voter"`
	ProducerNames Names  `json:"producer_names"`
	StandardTransactionOptions
}

type Approve struct {
	Voter        common.Name `json:"voter"`
	ProducerName common.Name `json:"producer_name"`
	StandardTransactionOptions
}

type UnapproveProducer struct {
	Voter        common.Name `json:"voter"`
	ProducerName common.Name `json:"producer_name"`
	StandardTransactionOptions
}

type ListproducersParams struct {
	PrintJson bool   `json:"print_json"`
	Limit     uint32 `json:"limit"`
	Lower     string `json:"lower"`
}

type DelegatebwParams struct {
	From               string `json:"from"`
	Receiver           string `json:"receiver"`
	StakeNetAmount     string `json:"stake_net_amount"`
	StakeCpuAmount     string `json:"stake_cpu_amount"`
	StakeStorageAmount string `json:"stake_storage_amount"`
	BuyRamAmount       string `json:"buy_ram_amount"`
	BuyRamBytes        uint32 `json:"buy_ram_bytes"`
	Transfer           bool   `json:"transfer"`
	StandardTransactionOptions
}

type UndelegatebwParams struct {
	From             string `json:"from"`
	Receive          string `json:"receive"`
	UnstakeNetAmount string `json:"unstake_net_amount"`
	UnstakeCpuAmount string `json:"unstake_cpu_amount"`
	StandardTransactionOptions
}

type ListbwParams struct {
	Account   common.Name `json:"name"`
	PrintJson bool        `json:"print_json"`
}

type BidnameParams struct {
	Bidder    string `json:"bidder"`
	NewName   string `json:"newname"`
	BidAmount string `json:"bid_amount"`
	StandardTransactionOptions
}

type BidNameinfoParams struct {
	PrintJson bool   `json:"print_json"`
	Newname   string `json:"newname"`
}

type BuyramParams struct {
	Payer     string `json:"payer"`
	Receiver  string `json:"receiver"`
	Amount    string `json:"amount"`
	Kbytes    bool   `json:"kbytes"`
	BytesFlag bool   `json:"bytes"`
	StandardTransactionOptions
}

type SellRamParams struct {
	From     string `json:"from"`
	Receiver string `json:"receiver"`
	Amount   uint64 `json:"amount"`
	StandardTransactionOptions
}

type ClaimrewardsParams struct {
	Owner string `json:"owner"`
	StandardTransactionOptions
}

type RegproxyParams struct {
	Proxy string `json:"proxy"`
	StandardTransactionOptions
}

type CanceldelayParams struct {
	CancelingAccount   string `json:"canceling_account"`
	CanclingPermission string `json:"canceling_permission"`
	TrxID              string `json:"trx_id"`
	StandardTransactionOptions
}

//chain
type GetCodeParams struct {
	AccountName  string `json:"name"`
	CodeFileName string `json:"code"`
	AbiFileName  string `json:"abi"`
	CodeAsWasm   bool   `json:"wasm"`
}
type GetTableParams struct {
	Code          string `json:"code"`
	Scope         string `json:"scope"`
	Table         string `json:"table"`
	Binary        bool   `json:"binary"`
	Limit         uint32 `json:"limit"` //default =10
	TableKey      string `json:"key"`
	Lower         string `json:"lower"`
	Upper         string `json:"upper"`
	IndexPosition string `json:"index"`
	KeyType       string `json:"key_type"`
	EncodeType    string `json:"encode_type"` //default ='dec'
}
type GetScopeParams struct {
	Code  string `json:"code"`
	Table string `json:"table"`
	Limit uint32 `json:"limit"` //default =10
	Lower string `json:"lower"`
	Upper string `json:"upper"`
}

type Names []common.Name

func (n Names) Len() int {
	return len(n)
}
func (n Names) Less(i, j int) bool {
	return uint64(n[i]) < uint64(n[j])
}
func (n Names) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type ProposeParams struct {
	ProposalName            string `json:"proposal_name"`
	RequestedPerm           string `json:"requested_permissions"`
	TransactionPerm         string `json:"trx_permissions"`
	ProposedContract        string `json:"contract"`
	ProposedAction          string `json:"action"`
	ProposedTransaction     string `json:"data"`
	Proposer                string `json:"proposer"`
	ProposalExpirationHours int    `json:"proposal_expiration"` //TODO default 24
	StandardTransactionOptions
}

type ProposeTrxParams struct {
	ProposalName  string `json:"proposal_name"`
	RequestedPerm string `json:"requested_permissions"`
	Proposer      string `json:"proposer"`
	TrxToPush     string `json:"transaction"`
	StandardTransactionOptions
}

type ReviewParams struct {
	ProposalName string `json:"proposal_name"`
	Proposer     string `json:"proposer"`
}

type ApproveAndUnapproveParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Perm         string `json:"permissions"`
	StandardTransactionOptions
}

type CancelParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Canceler     string `json:"canceler"`
	StandardTransactionOptions
}

type ExecuteParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Executer     string `json:"executer"`
	StandardTransactionOptions
}
