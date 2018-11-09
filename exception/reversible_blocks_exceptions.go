package exception

import . "github.com/eosspark/eos-go/log"

type ReversibleBlocksException struct{ LogMessage }

func (ReversibleBlocksException) ChainExceptions()            {}
func (ReversibleBlocksException) ReversibleBlocksExceptions() {}
func (ReversibleBlocksException) Code() ExcTypes              { return 3180000 }
func (ReversibleBlocksException) What() string {
	return "Reversible Blocks exception"
}

type InvalidReversibleBlocksDir struct{ LogMessage }

func (InvalidReversibleBlocksDir) ChainExceptions()            {}
func (InvalidReversibleBlocksDir) ReversibleBlocksExceptions() {}
func (InvalidReversibleBlocksDir) Code() ExcTypes              { return 3180001 }
func (InvalidReversibleBlocksDir) What() string {
	return "Invalid reversible blocks directory"
}

type ReversibleBlocksBackupDirExist struct{ LogMessage }

func (ReversibleBlocksBackupDirExist) ChainExceptions()            {}
func (ReversibleBlocksBackupDirExist) ReversibleBlocksExceptions() {}
func (ReversibleBlocksBackupDirExist) Code() ExcTypes              { return 3180002 }
func (ReversibleBlocksBackupDirExist) What() string {
	return "Backup directory for reversible blocks already exist"
}

type GapInReversibleBlocksDb struct{ LogMessage }

func (GapInReversibleBlocksDb) ChainExceptions()            {}
func (GapInReversibleBlocksDb) ReversibleBlocksExceptions() {}
func (GapInReversibleBlocksDb) Code() ExcTypes              { return 3180003 }
func (GapInReversibleBlocksDb) What() string {
	return "Gap in the reversible blocks database"
}
