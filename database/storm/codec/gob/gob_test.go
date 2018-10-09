package gob

import (
	"testing"

	"github.com/eosspark/eos-go/database/storm/codec/internal"
)

func TestGob(t *testing.T) {
	internal.RoundtripTester(t, Codec)
}
