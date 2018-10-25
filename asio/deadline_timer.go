package asio

import "time"

type DeadlineTimer struct {
	ctx *IoContext
	internal *time.Timer

	Duration time.Duration
}

func NewDeadlineTimer(ctx *IoContext) *DeadlineTimer {
	d := new(DeadlineTimer)
	d.ctx = ctx
	return d
}

func (d *DeadlineTimer) Expires(t time.Time) {
	d.Duration = t.Sub(time.Now())
}

func (d *DeadlineTimer) AsyncWait(op func(ec ErrorCode)) {
	// use go-timers to receive time event in new goroutine
	d.internal = time.AfterFunc(d.Duration, func() {
		d.ctx.GetService().push(op, NewErrorCode(nil))
	})
}

func (d *DeadlineTimer) Cancel() {
	if d.internal != nil {
		d.internal.Stop()
		d.internal = nil
	}
}
