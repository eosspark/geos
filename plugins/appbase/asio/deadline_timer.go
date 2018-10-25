package asio

import "time"

type DeadlineTimer struct {
	ctx *IoContext
	internal *time.Timer
	duration time.Duration
}

func NewDeadlineTimer(ctx *IoContext) *DeadlineTimer {
	d := new(DeadlineTimer)
	d.ctx = ctx
	return d
}

func (d *DeadlineTimer) ExpiresFromNow(duration time.Duration) {
	d.duration = duration
}

func (d *DeadlineTimer) ExpiresAt(t time.Time) {
	d.ExpiresFromNow(t.Sub(time.Now()))
}

func (d *DeadlineTimer) AsyncWait(op func(ec ErrorCode)) {
	// use go-timers to receive time event in new goroutine
	d.internal = time.AfterFunc(d.duration, func() {
		d.ctx.GetService().post(op, NewErrorCode(nil))
	})
}

func (d *DeadlineTimer) Cancel() {
	if d.internal != nil {
		d.internal.Stop()
		d.internal = nil
	}
}
