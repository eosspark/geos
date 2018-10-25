package asio

import (
	"os"
	"os/signal"
)

type SignalSet struct {
	ctx    *IoContext
	notify chan os.Signal
	cancel chan struct{}
}

func NewSignalSet(ctx *IoContext, sig ...os.Signal) *SignalSet {
	s := new(SignalSet)
	s.ctx = ctx

	s.notify = make(chan os.Signal, 1)
	signal.Notify(s.notify, sig...)

	s.cancel = make(chan struct{}, 1)

	return s
}

func (s *SignalSet) AsyncWait(op func(ec ErrorCode)) {
	go func() {
		for {
			select {
			case <-s.notify:
				// push into io_service
				// operation will be executed in the correct time
				s.ctx.GetService().post(op, ErrorCode{})
				break
			case <-s.cancel:
				return
			}
		}
	}()
}

func (s *SignalSet) Cancel() {
	s.cancel <- struct{}{}
}
