package asio

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	)

func TestIoContext_Run(t *testing.T) {
	ioc := NewIoContext()
	equal := 0
	go func() {
		time.Sleep(time.Second)
		ioc.Stop()
		equal = 1
	}()
	ioc.Run()
	assert.Equal(t, 1, equal)
}

func TestIoContext_Post(t *testing.T) {
	ioc := NewIoContext()
	ioc.Post(func() {
		ioc.Stop()
	})
	ioc.Run()
}

func TestIoContext_pRun(t *testing.T) {
	const COUNT = 1000000
	ioc := NewIoContext()
	j := 0

	//delay := NewDeadlineTimer(ioc)
	//delay.ExpiresFromNow(time.Second)
	//delay.AsyncWait(func(err error) {
	//	ioc.Stop()
	//})

	for i:=0; i<COUNT; i++ {
		ioc.Post(func() {
			j++
			if j == COUNT {
				ioc.Stop()
			}
		})
	}

	ioc.Run()
	assert.Equal(t, COUNT, j)
}
