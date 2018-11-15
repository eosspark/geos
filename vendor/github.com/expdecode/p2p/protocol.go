package p2p

const (
	HandshakeMessageTag = iota
	ChainSizeMessageTag
	GoawayMessageTag
	TimeMessageTag
	NoticeMessageTag
	RequestMessageTag
	SyncRequestMessageTag
	SignedBlockTag
	PackedTransactionTag
)

type Pubkey struct {
	Tag     uint8 `rlp:"tag"`
	Storage [33]byte
}

type SignatureType struct {
	Tag       uint8 `rlp:"tag"`
	Signature [65]byte
}

type HandshakeMessage struct {
	NetworkVersion           uint16
	ChainId                  [4]uint64
	NodeId                   [4]uint64
	Key                      Pubkey
	Timestamp                uint64
	Token                    [32]byte
	Sig                      SignatureType
	P2PAddress               string
	LastIrreversibleBlockNum uint32
	LastIrreversibleBlockId  [4]uint64
	HeadNum                  uint32
	HeadId                   [4]uint64
	Os                       string
	Agent                    string
	Generation               uint16
}

type ChainSizeMessage struct {
	LastIrreversibbleBlockNum uint32
	LastIrrversibleBlockId    [4]uint64 //block_id_type
	HeadNum                   uint32
	HeadId                    [4]uint64 //block_id_type
}

type Time_message_struct struct {
	Org uint64 //!< origin timestamp
	Rec uint64 //!< receive timestamp
	Xmt uint64 //!< transmit timestamp
	Dst uint64 //!< destination timestamp
}

type ProducerKey struct {
	ProducerName    uint64   //account_name
	BlockSigningKey [33]byte //public_key_type
}

type ProducerScheduleType struct {
	Version   uint32
	Producers []*ProducerKey
}
type OptionalPST struct {
	// ProducerValid bool //eos raw.hpp 278
	// Pst
	Pst ProducerScheduleType
}
type ExtensionType struct {
	Num  uint16
	Data []byte
}

// type HeaderExtension map[uint16][]*byte

type StaticVariant struct {
	Tag     uint8 `rlp:"tag"`
	PackTRX PackedTransaction
}

type TransactionReceipt struct {
	Status        uint8
	CPUUasgeUs    uint32
	NetUsageWords uint8 //uint32  todo
	Trx           StaticVariant
}

type BlockExtension struct {
	Num  uint16
	Data []byte
}

type SignedBlock struct {
	Timestamp        uint32
	Producer         uint64
	Confirmed        uint16
	Previous         [4]uint64
	TransactionMroot [4]uint64
	ActionMroot      [4]uint64
	ScheduleVersion  uint32
	//TODO 激活后以后的NewProduer 是有值的
	// NewProducers     bool //ProducerScheduleType
	NewProducers *OptionalPST `eos:"optional"`

	HeaderExtensions []*ExtensionType

	//以上为block_header

	ProducerSignature SignatureType //[65]byte //fc::array<unsigned char,65>
	//以上为signed_block_header

	Transactoins    []*TransactionReceipt
	BlockExtensions []*BlockExtension
}

const ( //for signed block status
	Executed = iota // succeed, no error handler executed
	SoftFail        // objectively failed (not executed), error handler executed
	HardFail        // objectively failed and error handler objectively failed thus no state change
	Delayed         // transaction delayed/deferred/scheduled for future execution
	Expired         // transaction expired and storage space refuned to user
)

type PermissionLevel struct {
	Actor      uint64 //account_name
	Permission uint64 //permission_name
}
type Action struct {
	ActionAccount uint64 //account_name
	ActionName    uint64 //action_name
	Authorization []*PermissionLevel
	Data          []byte //bytes vector<char>
}

type Transaction struct {
	Expiration     uint32 //time_point_sec
	RefBlockNum    uint16
	RefBlockPrefix uint32
	NetUsageWords  uint
	MaxCPUUsageMs  uint8
	DelaySec       uint

	//以上为transaction_header
	ContextFreeActions   []*Action //vector<action>
	Actions              []*Action
	TransactionExtension ExtensionType
	//以上为transaction
	Signatures      []*SignatureType
	ConTextFreeData []*byte
}

// type OptionalTransactoin struct {
// 	Valid bool
// 	// Value StorageType //typedef typename std::aligned_storage<sizeof(T), alignof(T)>::type storage_type;
// 	// Value
// }

// type Compression uint8

const (
	CompressionNone = iota
	CompressionZlib
)

var CompressionToString = map[uint8]string{
	CompressionNone: "None",
	CompressionZlib: "zlib",
}

type PackedTransaction struct {
	Signatures            []*SignatureType
	TrxCompression        uint8
	PackedContextFreeData []byte //bytes vector<char>
	PackedTrx             []byte //bytes vector<char>
}

//type Reason uint32
const (
	//NoReason      Reason = iota // no reason to go away
	NoReason       = iota // no reason to go away
	Self                  //the connection is to itself
	Duplicate             //the connection is redundant
	WrongChain            //the peer's chain id doesn't match
	WrongVersion          //the peer's network version doesn't match
	Forked                //peer's irreversible blocks are different
	Unlinkable            //the peer sent a block we couldn't use
	BadTransaction        //the peer sent a block that failed validation
	Validation            //the peer sent a block that failed validation
	BenignOther           //reasons such as a timeout. not fatal but warrant resetting
	FatalOther            //a catch-all for errors we don't have discriminated
	Authentication        //peer failed authenicatio
)

var ReasonToString = map[uint32]string{
	NoReason:       "no reason",
	Self:           "self connect",
	Duplicate:      "duplicate",
	WrongChain:     "wrong chain",
	WrongVersion:   "wrong version",
	Forked:         "chain is forked",
	Unlinkable:     "unlinkable block received",
	BadTransaction: "bad transaction",
	Validation:     "invalid block",
	BenignOther:    "some other non-fatal condition",
	FatalOther:     "some other failure",
	Authentication: "authentication failure",
}

type GoAwayMessage struct {
	//TODO
	Reason uint32
	NodeId [4]uint64
}

//TODO 默认入参？
func (g *GoAwayMessage) Init(r uint32) {
	g.Reason = r
	g.NodeId = [4]uint64{0, 0, 0, 0}
}

const (
	None = iota
	CatchUp
	LastIrrCatchUp
	Normal
)

var ModrsTostring = map[uint32]string{
	None:           "none",
	CatchUp:        "catch up",
	LastIrrCatchUp: "last irreversible",
	Normal:         "normal",
	//"underfined mode"
}

type SyncRequestMessage struct {
	StartBlock uint32
	EndBlock   uint32
}

type NoticeMessage struct {
	KnownTrx    OrderedTxnIds
	KnownBlocks OrderedBLKIds
}

type transaction_id_type [32]byte

type OrderedTxnIds struct {
	Mode    uint32
	Pending uint32
	Ids     []*transaction_id_type
}

func (s *OrderedTxnIds) empty() bool {
	return s.Mode == None
}

type block_id_type [4]uint64
type OrderedBLKIds struct {
	Mode    uint32
	Pending uint32
	Ids     []*block_id_type
}

func (s *OrderedBLKIds) empty() bool {
	return s.Mode == None
}

type RequestMessage struct {
	ReqTrx    OrderedTxnIds
	ReqBlocks OrderedBLKIds
}

// #define MAX_NUM_ARRAY_ELEMENTS (1024*1024)
// #define MAX_SIZE_OF_BYTE_ARRAYS (20*1024*1024)
