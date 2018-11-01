package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Trace("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)
	Debug("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)
	Info("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)
	Error("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)
	Warn("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)
	Crit("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)

}

func TestLog2(t *testing.T) {
	h, _ := FileHandler("./log.log", TerminalFormat(true))
	Root().SetHandler(h)
	Warn("Sanitizing cache to Go's GC limits", "provided", 100, "updated", true)

}

func TestLog3(t *testing.T) {
	srvlog := New("module", "app/server")
	srvlog.Warn("net Plugin", "rate", 100, "low", 9, "high", 3.4)
}
