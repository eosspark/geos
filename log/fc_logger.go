package log

import (
	"os"
	"path/filepath"
	"runtime"
)

type FcLoggerMap struct {
	initialed bool
	path      string
	m         map[string]*logger
}

var loggerMap = FcLoggerMap{}

func (fl *FcLoggerMap) appendLogger(names ...string) {
	for _, name := range names {
		path := fl.path + name + ".log"

		h, err := FileHandler(path, LogfmtFormat())
		if err != nil {
			root.Warn("logger[%s] add failed: %s", name, err.Error())
			return
		}

		l := &logger{name: name, h: new(swapHandler), enable: true}
		l.SetHandler(h)
		fl.m[name] = l
	}
}

func GetLoggerMap() map[string]*logger {
	if !loggerMap.initialed {
		loggerMap.path = getDefaultLogDir() + "/"
		os.MkdirAll(loggerMap.path, os.ModePerm)
		loggerMap.m = make(map[string]*logger)
		loggerMap.appendLogger(fcLogs...)
		loggerMap.initialed = true
	}
	return loggerMap.m
}

func getDefaultLogDir() string {

	home := os.Getenv("HOME")
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Application Support", "eosgo", "log")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "eosgo", "log")
		} else {
			return filepath.Join(home, ".clef")
		}
	}
	return "./log"
}
