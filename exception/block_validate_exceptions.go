package exception

type BlockValidateException struct{ logMessage }

func (e *BlockValidateException) ChainExceptions()         {}
func (e *BlockValidateException) BlockValidateExceptions() {}
func (e *BlockValidateException) Code() ExcTypes           { return 3030000 }
func (e *BlockValidateException) What() string             { return "Action validate exception" }

type UnlinkableBlockException struct{ logMessage }

func (e *UnlinkableBlockException) ChainExceptions()         {}
func (e *UnlinkableBlockException) BlockValidateExceptions() {}
func (e *UnlinkableBlockException) Code() ExcTypes           { return 3030001 }
func (e *UnlinkableBlockException) What() string             { return "Unlinkable block" }

type BlockFromTheFuture struct{ logMessage }

func (e *BlockFromTheFuture) ChainExceptions()         {}
func (e *BlockFromTheFuture) BlockValidateExceptions() {}
func (e *BlockFromTheFuture) Code() ExcTypes           { return 3030007 }
func (e *BlockFromTheFuture) What() string             { return "Block is from the future" }

type WrongSigningKey struct{ logMessage }

func (e *WrongSigningKey) ChainExceptions()         {}
func (e *WrongSigningKey) BlockValidateExceptions() {}
func (e *WrongSigningKey) Code() ExcTypes           { return 3030008 }
func (e *WrongSigningKey) What() string             { return "Block is not signed with expected key" }

type WrongProducer struct{ logMessage }

func (e *WrongProducer) ChainExceptions()         {}
func (e *WrongProducer) BlockValidateExceptions() {}
func (e *WrongProducer) Code() ExcTypes           { return 3030009 }
func (e *WrongProducer) What() string             { return "Block is not signed by expected producer" }
