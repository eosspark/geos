package log

import (
	"testing"
)

func TestLogGlobal(t *testing.T) {
	Debug("log terminalHandler test %s", "walker")
	Info("log terminalHandler test %s", "walker")
	Warn("log terminalHandler test %s", "walker")
	Error("log terminalHandler test %s", "walker")
}

func Test_log(t *testing.T) {
	newlog := New("test")
	newlog.SetHandler(TerminalHandler)

	newlog.Debug("log terminalHandler test %s", "walker")
	newlog.Info("log terminalHandler test %s", "walker")
	newlog.Warn("log terminalHandler test %s", "walker")
	newlog.Error("log terminalHandler test %s", "walker")
}

func TestFile(t *testing.T) {
	newlog := New("test")
	h, _ := FileHandler("./test.log", LogfmtFormat())
	root.SetHandler(h)

	newlog.Debug("log terminalHandler test %s", "walker")
	newlog.Info("log terminalHandler test %s", "walker")
	newlog.Warn("log terminalHandler test %s", "walker")
	newlog.Error("log terminalHandler test %s", "walker")
}

func TestFilterLog(t *testing.T) {
	newlog := New("test")
	root.SetHandler(LvlFilterHandler(LvlWarn, TerminalHandler))

	newlog.Debug("log terminalHandler test %s", "walker")
	newlog.Info("log terminalHandler test %s", "walker")
	newlog.Warn("log terminalHandler test %s", "walker")
	newlog.Error("log terminalHandler test %s", "walker")
}
