package chain

type ChannelsType int

const (
	PreAcceptedBlock = ChannelsType(iota + 1)
	RejectedBlock
	AcceptedBlockHeader
	AcceptedBlock
	IrreversibleBlock
	AcceptedTransaction
	AppliedTransaction
	AcceptedConfirmation
)

type MethodsType int

const (
	GetBlockByNumber = MethodsType(iota + 1)
	GetBlockById
	GetHeadBlockId
	GetLibBlockId

	GetLastIrreversibleBlockNumber
)

type IncomingChannels int

const (
	Block = IncomingChannels(iota + 1)
	Transaction
)

type IncomingMethods int

const (
	BlockSync = IncomingMethods(iota + 1)
	TransactionAsync
)

type CompatChannels int

const (
	TransactionAck = CompatChannels(iota + 1)
)
