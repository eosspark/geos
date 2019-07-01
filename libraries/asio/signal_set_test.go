package asio

import (
	"syscall"
	"testing"
	"time"
)

func TestNewSignalSet(t *testing.T) {
	ioc := NewIoContext()
	sigint := NewSignalSet(ioc, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		ioc.Stop()
		sigint.Cancel()
	})

	sigterm := NewSignalSet(ioc, syscall.SIGTERM)
	sigterm.AsyncWait(func(err error) {
		ioc.Stop()
		sigterm.Cancel()
	})

	delay := NewDeadlineTimer(ioc)
	delay.ExpiresFromNow(time.Millisecond)
	delay.AsyncWait(func(err error) {
		sigint.notify <- syscall.SIGINT
		sigterm.notify <- syscall.SIGTERM
	})

	ioc.Run()
}
