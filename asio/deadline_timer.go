package gosio

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

func (d *DeadlineTimer) Expires(t time.Time) {
	d.duration = t.Sub(time.Now())
}

func (d *DeadlineTimer) AsyncWait(callback func(ec ErrorCode)) {
	d.internal = time.AfterFunc(d.duration, func() {
		d.ctx.GetService().push(callback, NewErrorCode(nil))
	})
}

func (d *DeadlineTimer) Cancel() {
	if d.internal != nil {
		d.internal.Stop()
		d.internal = nil
	}
}
