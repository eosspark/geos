package types

import (
	"fmt"
	"time"
	"github.com/eosspark/eos-go/common"
	."github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/common/math"
)

type BlockTimeStamp uint32

func NewBlockTimeStamp(t common.TimePoint) BlockTimeStamp {
	microSinceEpoch := t.TimeSinceEpoch()
	msecSinceEpoch := microSinceEpoch / 1000
	slot := (int64(msecSinceEpoch) - common.DefaultConfig.BlockTimestampEpochMs) / common.DefaultConfig.BlockIntervalMs
	return BlockTimeStamp(slot)
}

func NewBlockTimeStampSec(ts common.TimePointSec) BlockTimeStamp {
	secSinceEpoch := ts.SecSinceEpoch()
	slot := (int64(secSinceEpoch*1000) - common.DefaultConfig.BlockTimestampEpochMs) / common.DefaultConfig.BlockIntervalMs
	return BlockTimeStamp(slot)
}

const blockTimestampFormat = "2006-01-02T15:04:05.000"

func (t BlockTimeStamp) Next() BlockTimeStamp {
	EosAssert(math.MaxUint32 - t >= 1, &OverflowException{}, "block timestamp overflow")
	result := NewBlockTimeStamp(t.ToTimePoint())
	result += 1
	return result
}

func (t BlockTimeStamp) ToTimePoint() common.TimePoint {
	msec := int64(t) * int64(common.DefaultConfig.BlockIntervalMs)
	msec += int64(common.DefaultConfig.BlockTimestampEpochMs)
	return common.TimePoint(common.Milliseconds(msec))
}

func MaxBlockTime() BlockTimeStamp {
	return BlockTimeStamp(0xffff)
}

func MinBlockTime() BlockTimeStamp {
	return BlockTimeStamp(0)
}

func (t BlockTimeStamp) String() string {
	return t.ToTimePoint().String()
}

func (t BlockTimeStamp) MarshalJSON() ([]byte, error) {
	var slot int64
	if t > 0 {
		slot = int64(t)*common.DefaultConfig.BlockIntervalMs*1000000 + common.DefaultConfig.BlockTimestamoEpochNanos //为了显示0.5s
	} else {
		slot = 0 //"1970-01-01T00:00:00.000"
	}
	tm := time.Unix(0, int64(slot)).UTC()

	return []byte(fmt.Sprintf("%q", tm.Format(blockTimestampFormat))), nil
}

func (t *BlockTimeStamp) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	var temp time.Time
	temp, err = time.Parse(`"`+blockTimestampFormat+`"`, string(data))
	if err != nil {
		return
	}
	slot := (temp.UnixNano() - common.DefaultConfig.BlockTimestamoEpochNanos) / 1e6 / common.DefaultConfig.BlockIntervalMs
	if slot < 0 {
		slot = 0
	}

	*t = BlockTimeStamp(slot)
	return err
}

//func (t BlockTimeStamp) Totime() time.Time {
//	slot := int64(t)*common.DefaultConfig.BlockIntervalMs*1000000 + common.DefaultConfig.BlockTimestamoEpochNanos //为了显示0.5s
//	return time.Unix(0, int64(slot)).UTC()
//}
