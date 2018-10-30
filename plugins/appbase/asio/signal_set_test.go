package asio

import (
	"testing"
	"syscall"
)

func TestNewSignalSet(t *testing.T) {
	ioc := NewIoContext()
	sigint := NewSignalSet(ioc, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		ioc.Stop()
		sigint.Cancel()
	})

	ioc.Run()
}
