package exception

type ControllerEmitSignalException struct{ logMessage }

func (e *ControllerEmitSignalException) ChainExceptions()                {}
func (e *ControllerEmitSignalException) ControllerEmitSignalExceptions() {}
func (e *ControllerEmitSignalException) Code() ExcTypes                  { return 3140000 }
func (e *ControllerEmitSignalException) What() string {
	return "Exceptions that are allowed to bubble out of emit calls in controller"
}

type CheckpointException struct{ logMessage }

func (e *CheckpointException) ChainExceptions()                {}
func (e *CheckpointException) ControllerEmitSignalExceptions() {}
func (e *CheckpointException) Code() ExcTypes                  { return 3140001 }
func (e *CheckpointException) What() string {
	return "Block does not match checkpoint"
}
