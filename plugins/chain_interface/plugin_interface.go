package chain_interface

type ChannelsType int

const (
	PreAcceptedBlock = ChannelsType(iota)
	RejectedBlock
	AcceptedBlockHeader
	AcceptedBlock
	IrreversibleBlock
	AcceptedTransaction
	AppliedTransaction
	AcceptedConfirmation

	//incoming
	Block
	Transaction

	//compat
	TransactionAck
)

type MethodsType int

const (
	GetBlockByNumber = MethodsType(iota)
	GetBlockById
	GetHeadBlockId
	GetLibBlockId

	GetLastIrreversibleBlockNumber

	//incoming
	BlockSync
	TransactionAsync
)
