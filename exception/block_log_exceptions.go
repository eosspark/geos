package exception

type BlockLogException struct{ logMessage }

func (e *BlockLogException) ChainExceptions()    {}
func (e *BlockLogException) BlockLogExceptions() {}
func (e *BlockLogException) Code() ExcTypes      { return 3190000 }
func (e *BlockLogException) What() string        { return "Block log exception" }
