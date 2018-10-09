package json

import (
	"testing"

	"github.com/eosspark/eos-go/database/storm/codec/internal"
)

func TestJSON(t *testing.T) {
	internal.RoundtripTester(t, Codec)
}
