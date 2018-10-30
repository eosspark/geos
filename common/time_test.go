package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
		_ "net/http/pprof"
		"testing"
		)

func Test_TimePoint(t *testing.T) {
	fmt.Println(MaxTimePoint())
	fmt.Println(MinTimePoint())
	now := Now()
	fmt.Println(now, now.TimeSinceEpoch())

	fmt.Println(MaxTimePointSec())
	fmt.Println(MinTimePointSec())
}

func Test_FromIsoString(t *testing.T) {
	s := "2006-01-02T15:04:05.500"
	tp, e := FromIsoString(s)
	assert.NoError(t, e)

	tps, err := FromIsoStringSec(s)
	assert.NoError(t, err)

	fmt.Println(tp)
	fmt.Println(tps)
}

func Test_BlockTimestamp(t *testing.T) {}
