package producer_plugin

import (
	"fmt"
	"time"
)

type scheduleTimer struct {
	internal *time.Timer
	duration time.Duration
}

func (pt *scheduleTimer) expiresFromNow(d time.Duration) {
	pt.duration = d
}

func (pt *scheduleTimer) expiresAt(ex time.Time) {
	pt.expiresFromNow(time.Until(ex))
}

func (pt *scheduleTimer) asyncWait(valid func() bool, call func()) {
	pt.internal = time.NewTimer(pt.duration)
	<-pt.internal.C
	if valid() {
		go call()
	} else {
		fmt.Println("no call")
	}
}

func (pt *scheduleTimer) cancel() {
	if pt.internal != nil {
		pt.internal.Stop()
		pt.internal = nil
	}
}
