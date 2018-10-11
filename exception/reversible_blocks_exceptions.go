package exception

type ReversibleBlocksException struct{ logMessage }

func (e *ReversibleBlocksException) ChainExceptions()            {}
func (e *ReversibleBlocksException) ReversibleBlocksExceptions() {}
func (e *ReversibleBlocksException) Code() ExcTypes              { return 3180000 }
func (e *ReversibleBlocksException) What() string {
	return "Reversible Blocks exception"
}
