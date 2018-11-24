package wasmgo

import (
//"fmt"
)

func readActionData(w *WasmGo, memory int, bufferSize int) int {

	if bufferSize > (1<<16) || memory+bufferSize > (1<<16) {
		w.ilog.Error("access violation")
		return 0
	}

	data := w.context.GetActionData()
	s := len(data)
	if bufferSize == 0 {
		return s
	}

	copySize := min(bufferSize, s)
	setMemory(w, memory, data, 0, copySize)

	//w.ilog.Info("action data:%v size:%d", data, copySize)
	return copySize
}

func actionDataSize(w *WasmGo) int {
	// size := len(w.context.GetActionData())
	// w.ilog.Info("actionDataSize:%d", size)
	return len(w.context.GetActionData())
}

func currentReceiver(w *WasmGo) int64 {
	//    receiver := w.context.GetReceiver()
	//    w.ilog.Info("currentReceiver:%v", receiver)
	//    return int64(receiver)
	return int64(w.context.GetReceiver())
}
