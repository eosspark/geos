package exception

type ControllerEmitSignalException struct{ logMessage }

func (ControllerEmitSignalException) ChainExceptions()                {}
func (ControllerEmitSignalException) ControllerEmitSignalExceptions() {}
func (ControllerEmitSignalException) Code() ExcTypes                  { return 3140000 }
func (ControllerEmitSignalException) What() string {
	return "Exceptions that are allowed to bubble out of emit calls in controller"
}

type CheckpointException struct{ logMessage }

func (CheckpointException) ChainExceptions()                {}
func (CheckpointException) ControllerEmitSignalExceptions() {}
func (CheckpointException) Code() ExcTypes                  { return 3140001 }
func (CheckpointException) What() string {
	return "Block does not match checkpoint"
}
