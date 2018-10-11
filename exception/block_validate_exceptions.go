package exception

type BlockValidateException struct{ logMessage }

func (e *BlockValidateException) ChainExceptions()         {}
func (e *BlockValidateException) BlockValidateExceptions() {}
func (e *BlockValidateException) Code() ExcTypes           { return 3030000 }
func (e *BlockValidateException) What() string             { return "Action validate exception" }

