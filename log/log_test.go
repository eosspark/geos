package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Info("found block for id at num %d %s", 100, "walker")
	Debug("found block for id at num %d %s", 100, "walker")
	Warn("found block for id at num %d %s", 100, "walker")
	Error("found block for id at num %d %s", 100, "walker")
}

func TestLog2(t *testing.T) {
	h, _ := FileHandler("./log.log", TerminalFormat(true))
	Root().SetHandler(h)
	Info("found block for id at num %d %s", 100, "walker")

}

func TestLog3(t *testing.T) {
	srvlog := New()
	srvlog.Error("found block for id at num %d %s", 100, "walker")
}
