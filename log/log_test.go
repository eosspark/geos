package log

import (
	"bytes"
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

func TestLongdata(t *testing.T) {
	testlog := New("test")
	testlog.SetHandler(TerminalHandler)
	data := bytes.Repeat([]byte{100, 255, 0, '\n'}, 100000)
	testlog.Debug("%#v", data)
}

func TestFcLogMessage(t *testing.T) {
	makeLog := func(msg string) {
		message := FcLogMessage(LvlInfo, "test FC_LOG_MESSAGE %s", msg)
		fmt.Println(message.GetMessage())
		fmt.Println(message.GetContext().String())
		//assert.Equal(t, "log_test.go:92", message.GetContext().String())
		//assert.Equal(t, "test FC_LOG_MESSAGE message", message.GetMessage())
	}
	makeLog("message")
}
