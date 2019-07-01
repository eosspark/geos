package unittests

import (
	"testing"

	. "github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"

	"github.com/stretchr/testify/assert"
)

func TestTimePointToBlockTimeStamp(t *testing.T) {
	tp := common.TimePoint(common.Seconds(978307200))
	bt := NewBlockTimeStamp(tp)
	assert.Equal(t, uint32(978307200-946684800)*2, uint32(bt),
		"Time point constructor gives wrong value")
}

func TestBlockTimeStamp_ToTimePoint(t *testing.T) {
	bt := BlockTimeStamp(0)
	tp := bt.ToTimePoint()
	assert.Equal(t, int64(946684800), tp.TimeSinceEpoch().ToSeconds(), "Time point conversion failed")

	bt1 := BlockTimeStamp(200)
	tp = bt1.ToTimePoint()
	assert.Equal(t, int64(946684900), tp.TimeSinceEpoch().ToSeconds(), "Time point conversion failed")

}

func TestBlockTimeStamp_MarshalJSON(t *testing.T) {
	tp := common.TimePoint(common.Seconds(978307200))
	bt := NewBlockTimeStamp(tp)

	json, err := bt.MarshalJSON()
	assert.NoError(t, err, "BlockTimeStamp MarshalJSON failed")

	valid := NewBlockTimeStamp(0)
	err = valid.UnmarshalJSON(json)
	assert.NoError(t, err, "BlockTimeStamp UnmarshalJSON failed")

	assert.Equal(t, valid, bt)
}
