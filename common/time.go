package common

import (
	"time"
		"strings"
	"strconv"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
)

const format = "2006-01-02T15:04:05"

type Microseconds int64

func MaxMicroseconds() Microseconds { return Microseconds(0x7fffffffffffffff) }
func MinMicroseconds() Microseconds { return Microseconds(0) }

func (ms Microseconds) ToSeconds() int64        { return int64(ms / 1e6) }
func (ms Microseconds) Count() int64            { return int64(ms) }
func (ms Microseconds) String() string          { return TimePoint(ms).String() }

func Seconds(s int64) Microseconds      { return Microseconds(s * 1e6) }
func Milliseconds(s int64) Microseconds { return Microseconds(s * 1e3) }
func Minutes(m int64) Microseconds      { return Seconds(60 * m) }
func Hours(h int64) Microseconds        { return Minutes(60 * h) }
func Days(d int64) Microseconds         { return Hours(24 * d) }

type TimePoint Microseconds

func Now() TimePoint          { return TimePoint(time.Now().UTC().UnixNano() / 1e3) }
func MaxTimePoint() TimePoint { return TimePoint(MaxMicroseconds()) }
func MinTimePoint() TimePoint { return TimePoint(MinMicroseconds()) }

func (tp TimePoint) TimeSinceEpoch() Microseconds { return Microseconds(tp) }
func (tp TimePoint) SecSinceEpoch() uint32        { return uint32(tp) / 1e6 }
func (tp TimePoint) String() string {
	return time.Unix(int64(tp)/1e6, int64(tp)%1e6*1000).UTC().String()
}

func (tp TimePoint) MarshalJSON() ([]byte, error) {
	return []byte(tp.String()), nil
}

func (tp *TimePoint) UnmarshalJSON(data []byte) (err error) {
	*tp, err = FromIsoString(string(data))
	return
}

func FromIsoString(s string) (TimePoint, error) {
	if strings.IndexByte(s, '.') < 0 {
		tps, err := FromIsoStringSec(s)
		if err != nil {
			return 0, err
		}
		return tps.ToTimePoint(), nil
	} else {
		tps, err := FromIsoStringSec(strings.Split(s, ".")[0])
		if err != nil {
			return 0, err
		}
		subs := []byte(strings.Split(s, ".")[1])
		for len(subs) < 3 {
			subs = append(subs, '0')
		}
		ms, err2 := strconv.Atoi("1" + string(subs))
		if err2 != nil {
			return 0, err2
		}
		return tps.ToTimePoint().AddUs(Milliseconds(int64(ms) - 1000)), nil
	}
}

func (tp TimePoint) AddUs(m Microseconds) TimePoint     { return TimePoint(Microseconds(tp) + m) }
func (tp TimePoint) SubUs(m Microseconds) TimePoint     { return TimePoint(Microseconds(tp) - m) }
func (tp TimePoint) Sub(t TimePoint) Microseconds       { return Microseconds(tp - t) }
func (tp TimePoint) SubTps(t TimePointSec) Microseconds { return tp.Sub(t.ToTimePoint()) }

/**
 *  A lower resolution time_point accurate only to seconds from 1970
 */
type TimePointSec uint32

func NewTimePointSecTp(t TimePoint) TimePointSec { return TimePointSec(t.TimeSinceEpoch() / 1e6) }

func MaxTimePointSec() TimePointSec { return TimePointSec(0xffffffff) }
func MinTimePointSec() TimePointSec { return TimePointSec(0) }

func (tp TimePointSec) ToTimePoint() TimePoint  { return TimePoint(Seconds(int64(tp))) }
func (tp TimePointSec) SecSinceEpoch() uint32   { return uint32(tp) }
func (tp TimePointSec) String() string          { return tp.ToTimePoint().String() }

func (tp TimePointSec) MarshalJSON() ([]byte, error) {
	return []byte(tp.String()), nil
}

func (tp *TimePointSec) UnmarshalJSON(data []byte) (err error) {
	*tp, err = FromIsoStringSec(string(data))
	return
}

func FromIsoStringSec(s string) (TimePointSec, error) {
	pt, err := time.Parse(format, s)
	return TimePointSec(pt.Unix()), err
}

func (tp TimePointSec) AddSec(offset uint32) TimePointSec { return TimePointSec(uint32(tp) + offset) }
func (tp TimePointSec) SubSec(offset uint32) TimePointSec { return TimePointSec(uint32(tp) + offset) }
func (tp TimePointSec) AddUs(m Microseconds) TimePoint    { return tp.ToTimePoint().AddUs(m) }
func (tp TimePointSec) SubUs(m Microseconds) TimePoint    { return tp.ToTimePoint().SubUs(m) }
func (tp TimePointSec) Sub(t TimePointSec) Microseconds   { return tp.ToTimePoint().Sub(t.ToTimePoint()) }

/**
 * using asio.DeadlineTimer
 */
type Timer asio.DeadlineTimer

func NewTimer(ctx *asio.IoContext) *Timer {
	return (*Timer)(asio.NewDeadlineTimer(ctx))
}

func (t *Timer) ExpiresUntil(time TimePoint)  { t.ExpiresFromNow(time.Sub(Now())) }
func (t *Timer) ExpiresAt(epoch Microseconds) { t.ExpiresUntil(TimePoint(epoch)) }
func (t *Timer) ExpiresFromNow(m Microseconds) {
	(*asio.DeadlineTimer)(t).ExpiresFromNow(time.Microsecond * time.Duration(m))
}

func (t *Timer) Cancel()                      { (*asio.DeadlineTimer)(t).Cancel() }
func (t *Timer) AsyncWait(op func(err error)) { (*asio.DeadlineTimer)(t).AsyncWait(op) }
