package common

import (
	"fmt"
	"github.com/eosspark/eos-go/libraries/asio"
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"testing"
)

func Test_TimePoint(t *testing.T) {
	assert.Equal(t, "294247-01-10T04:00:54.775", MaxTimePoint().String(), "error max time")
	assert.Equal(t, "1970-01-01T00:00:00", MinTimePoint().String(), "error min time")

	assert.Equal(t, "2106-02-07T06:28:15", MaxTimePointSec().String(), "error max sec")
	assert.Equal(t, "1970-01-01T00:00:00", MinTimePointSec().String(), "error min sec")
}

func Test_FromIsoString(t *testing.T) {
	s := "2006-01-02T15:05:05.500"

	tp, e := FromIsoString(s)
	assert.NoError(t, e, "error create TimePoint from string")
	assert.Equal(t, "2006-01-02T15:05:05.5", tp.String(), "TimePoint from string is wrong")

	tps, err := FromIsoStringSec(s)
	assert.NoError(t, err, "error create TimePointSec from string")
	assert.Equal(t, "2006-01-02T15:05:05", tps.String(), "TimePointSec from string is wrong")
}

func TestTimer(t *testing.T) {
	timer := NewTimer((*asio.IoContext)(nil))
	timer.ExpiresFromNow(Milliseconds(1))
	timer.AsyncWait(func(err error) {})
}

func TestTimePoint_MarshalJSON(t *testing.T) {
	s := "2006-01-02T15:05:05.521"
	tp, e := FromIsoString(s)
	assert.NoError(t, e)
	json, e := tp.MarshalJSON()
	assert.NoError(t, e)
	fmt.Println(string(json))
	assert.Equal(t, "\"2006-01-02T15:05:05.521\"", string(json))
	//
	valid := Now()
	e = valid.UnmarshalJSON(json)
	assert.NoError(t, e)
	assert.Equal(t, valid, tp)

	tps, e := FromIsoStringSec(s)
	assert.NoError(t, e)
	json, e = tps.MarshalJSON()
	assert.NoError(t, e)
	fmt.Println(string(json))
	assert.Equal(t, "\"2006-01-02T15:05:05\"", string(json))
	//
	valids := TimePointSec(0)
	e = valids.UnmarshalJSON(json)
	assert.NoError(t, e)
	assert.Equal(t, valids, tps)
}
