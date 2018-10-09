package msgpack

import (
	"testing"

	"github.com/eosspark/eos-go/database/storm/codec/internal"
)

func TestMsgpack(t *testing.T) {
	internal.RoundtripTester(t, Codec)
}
