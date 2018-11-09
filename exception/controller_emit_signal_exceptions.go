package exception

import . "github.com/eosspark/eos-go/log"

type ControllerEmitSignalException struct{ LogMessage }

func (ControllerEmitSignalException) ChainExceptions()                {}
func (ControllerEmitSignalException) ControllerEmitSignalExceptions() {}
func (ControllerEmitSignalException) Code() ExcTypes                  { return 3140000 }
func (ControllerEmitSignalException) What() string {
	return "Exceptions that are allowed to bubble out of emit calls in controller"
}

type CheckpointException struct{ LogMessage }

func (CheckpointException) ChainExceptions()                {}
func (CheckpointException) ControllerEmitSignalExceptions() {}
func (CheckpointException) Code() ExcTypes                  { return 3140001 }
func (CheckpointException) What() string {
	return "Block does not match checkpoint"
}
