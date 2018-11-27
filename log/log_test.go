package log

import (
	"fmt"
	"testing"
	"time"
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

func TestLogTimeOn(t *testing.T) {
	testlog := New("test")
	testlog.SetHandler(TerminalHandler)
	start := time.Now()
	for i := 0; i < 10000; {
		testlog.Info("log terminalHandler test %s", "walker")
		i += 1
	}
	end := time.Now()

	g := end.Sub(start)
	fmt.Println(g)

}
func TestLogTimeOff(t *testing.T) {
	testlog := New("test")
	testlog.SetHandler(DiscardHandler())
	testlog.SetEnable(false)
	start := time.Now()
	for i := 0; i < 10000; {
		testlog.Info("log terminalHandler test %s", "walker")
		i += 1
	}
	end := time.Now()

	g := end.Sub(start)
	fmt.Println(g)
}
