package common

import (
	"fmt"
	"time"
)

type BlockTimeStamp uint32

func NewBlockTimeStamp(t TimePoint) BlockTimeStamp {
	microSinceEpoch := t.TimeSinceEpoch()
	msecSinceEpoch := microSinceEpoch / 1000
	slot := (int64(msecSinceEpoch) - DefaultConfig.BlockTimestampEpochMs) / DefaultConfig.BlockIntervalMs
	return BlockTimeStamp(slot)
}

func NewBlockTimeStampSec(ts TimePointSec) BlockTimeStamp {
	secSinceEpoch := ts.SecSinceEpoch()
	slot := (int64(secSinceEpoch*1000) - DefaultConfig.BlockTimestampEpochMs) / DefaultConfig.BlockIntervalMs
	return BlockTimeStamp(slot)
}

const blockTimestampFormat = "2006-01-02T15:04:05.000"

func (t BlockTimeStamp) Next() BlockTimeStamp {
	if 0xffffffff-t < 1 {
		panic("block timestamp overflow")
	}
	result := NewBlockTimeStamp(t.ToTimePoint())
	result += 1
	return result
}

func (t BlockTimeStamp) ToTimePoint() TimePoint {
	msec := int64(t) * int64(DefaultConfig.BlockIntervalMs)
	msec += int64(DefaultConfig.BlockTimestampEpochMs)
	return TimePoint(Milliseconds(msec))
}

func MaxBlockTime() BlockTimeStamp {
	return BlockTimeStamp(0xffffffff)
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
		slot = int64(t)*DefaultConfig.BlockIntervalMs*1000000 + DefaultConfig.BlockTimestamoEpochNanos //为了显示0.5s
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
	slot := (temp.UnixNano() - DefaultConfig.BlockTimestamoEpochNanos) / 1e6 / DefaultConfig.BlockIntervalMs
	if slot < 0 {
		slot = 0
	}

	*t = BlockTimeStamp(slot)
	return err
}

func (t BlockTimeStamp) Totime() time.Time {
	slot := int64(t)*DefaultConfig.BlockIntervalMs*1000000 + DefaultConfig.BlockTimestamoEpochNanos //为了显示0.5s
	return time.Unix(0, int64(slot)).UTC()
}
