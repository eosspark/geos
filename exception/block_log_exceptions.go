package exception

type BlockLogException struct{ logMessage }

func (e *BlockLogException) ChainExceptions()    {}
func (e *BlockLogException) BlockLogExceptions() {}
func (e *BlockLogException) Code() ExcTypes      { return 3190000 }
func (e *BlockLogException) What() string        { return "Block log exception" }

type BlockLogUnsupportedVersion struct{ logMessage }

func (e *BlockLogUnsupportedVersion) ChainExceptions()    {}
func (e *BlockLogUnsupportedVersion) BlockLogExceptions() {}
func (e *BlockLogUnsupportedVersion) Code() ExcTypes      { return 3190001 }
func (e *BlockLogUnsupportedVersion) What() string        { return "unsupported version of block log" }

type BlockLogAppendFail struct{ logMessage }

func (e *BlockLogAppendFail) ChainExceptions()    {}
func (e *BlockLogAppendFail) BlockLogExceptions() {}
func (e *BlockLogAppendFail) Code() ExcTypes      { return 3190002 }
func (e *BlockLogAppendFail) What() string        { return "fail to append block to the block log" }

type BlockLogNotFound struct{ logMessage }

func (e *BlockLogNotFound) ChainExceptions()    {}
func (e *BlockLogNotFound) BlockLogExceptions() {}
func (e *BlockLogNotFound) Code() ExcTypes      { return 3190003 }
func (e *BlockLogNotFound) What() string        { return "block log can not be found" }

type BlockLogBackupDirExist struct{ logMessage }

func (e *BlockLogBackupDirExist) ChainExceptions()    {}
func (e *BlockLogBackupDirExist) BlockLogExceptions() {}
func (e *BlockLogBackupDirExist) Code() ExcTypes      { return 3190004 }
func (e *BlockLogBackupDirExist) What() string        { return "block log backup dir already exists" }
