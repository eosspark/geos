package exception

import _ "github.com/eosspark/eos-go/log"

type BlockLogException struct{ ELog }

func (BlockLogException) ChainExceptions()    {}
func (BlockLogException) BlockLogExceptions() {}
func (BlockLogException) Code() ExcTypes      { return 3190000 }
func (BlockLogException) What() string        { return "Block log exception" }

type BlockLogUnsupportedVersion struct{ ELog }

func (BlockLogUnsupportedVersion) ChainExceptions()    {}
func (BlockLogUnsupportedVersion) BlockLogExceptions() {}
func (BlockLogUnsupportedVersion) Code() ExcTypes      { return 3190001 }
func (BlockLogUnsupportedVersion) What() string        { return "unsupported version of block log" }

type BlockLogAppendFail struct{ ELog }

func (BlockLogAppendFail) ChainExceptions()    {}
func (BlockLogAppendFail) BlockLogExceptions() {}
func (BlockLogAppendFail) Code() ExcTypes      { return 3190002 }
func (BlockLogAppendFail) What() string        { return "fail to append block to the block log" }

type BlockLogNotFound struct{ ELog }

func (BlockLogNotFound) ChainExceptions()    {}
func (BlockLogNotFound) BlockLogExceptions() {}
func (BlockLogNotFound) Code() ExcTypes      { return 3190003 }
func (BlockLogNotFound) What() string        { return "block log can not be found" }

type BlockLogBackupDirExist struct{ ELog }

func (BlockLogBackupDirExist) ChainExceptions()    {}
func (BlockLogBackupDirExist) BlockLogExceptions() {}
func (BlockLogBackupDirExist) Code() ExcTypes      { return 3190004 }
func (BlockLogBackupDirExist) What() string        { return "block log backup dir already exists" }
