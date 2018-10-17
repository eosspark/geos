package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	type Action struct {
		ActionAccount uint64
		Data          []byte
	}
	type Transaction struct {
		Expiration         uint32
		NetUsageWords      uint
		MaxCPUUsageMs      uint8
		DelaySec           uint
		ContextFreeActions []*Action
	}

	action1 := &Action{
		ActionAccount: 9876543,
		Data:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
	action2 := &Action{
		ActionAccount: 987654321,
		Data:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
	test := &Transaction{Expiration: 100,
		NetUsageWords:      9,
		MaxCPUUsageMs:      199,
		DelaySec:           99999,
		ContextFreeActions: []*Action{action1, action2},
	}
	assert.Equal(t, false, Empty(test))
	test = nil
	assert.Equal(t, true, Empty(test))

	test1 := Transaction{}
	assert.Equal(t, true, Empty(test1))

}
