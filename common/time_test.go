package common

import (
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"testing"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
)

func Test_TimePoint(t *testing.T) {
	assert.Equal(t, "294247-01-10 04:00:54.775807 +0000 UTC", MaxTimePoint().String(), "error max time")
	assert.Equal(t, "1970-01-01 00:00:00 +0000 UTC", MinTimePoint().String(), "error min time")

	assert.Equal(t, "2106-02-07 06:28:15 +0000 UTC", MaxTimePointSec(), "error max sec")
	assert.Equal(t, "1970-01-01 00:00:00 +0000 UTC", MinTimePointSec(), "error min sec")
}

func Test_FromIsoString(t *testing.T) {
	s := "2006-01-02T15:05:05.500"

	tp, e := FromIsoString(s)
	assert.NoError(t, e, "error create TimePoint from string")
	assert.Equal(t, "2006-01-02 15:05:05.5 +0000 UTC", tp.String(), "TimePoint from string is wrong")

	tps, err := FromIsoStringSec(s)
	assert.NoError(t, err, "error create TimePointSec from string")
	assert.Equal(t, "2006-01-02 15:05:05 +0000 UTC", tps.String(), "TimePointSec from string is wrong")
}

func TestTimer(t *testing.T) {
	timer := NewTimer((*asio.IoContext)(nil))
	timer.ExpiresFromNow(Milliseconds(1))
	timer.AsyncWait(func(err error) {})
}
