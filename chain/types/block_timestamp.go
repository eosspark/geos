package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/math"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
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
	EosAssert(math.MaxUint32-t >= 1, &OverflowException{}, "block timestamp overflow")
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
	return []byte(t.String()), nil
}

func (t *BlockTimeStamp) UnmarshalJSON(data []byte) (err error) {
	tp, err := common.FromIsoString(string(data))
	if err != nil {
		return err
	}
	*t = BlockTimeStamp((int64(tp.TimeSinceEpoch()/1000) - common.DefaultConfig.BlockTimestampEpochMs) / common.DefaultConfig.BlockIntervalMs)
	return nil
}

//func (t BlockTimeStamp) Totime() time.Time {
//	slot := int64(t)*common.DefaultConfig.BlockIntervalMs*1000000 + common.DefaultConfig.BlockTimestamoEpochNanos //为了显示0.5s
//	return time.Unix(0, int64(slot)).UTC()
//}
