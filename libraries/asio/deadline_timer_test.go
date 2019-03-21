package asio

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestDeadlineTimer_DuplicateTimer(t *testing.T) {
	ioc := NewIoContext()
	ta, tb, tc := NewDeadlineTimer(ioc), NewDeadlineTimer(ioc), NewDeadlineTimer(ioc)
	checks := make([]int, 0, 3)

	ta.ExpiresFromNow(time.Millisecond)
	ta.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		checks = append(checks, 4)
	})

	tb.ExpiresFromNow(time.Millisecond * 2)
	tb.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		checks = append(checks, 5)
	})

	tc.ExpiresFromNow(time.Millisecond * 3)
	tc.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		checks = append(checks, 6)
		ioc.Stop()
	})

	ioc.Run()
	assert.Equal(t, []int{4, 5, 6}, checks)
}

func TestDeadlineTimer_AsyncWaitRef(t *testing.T) {
	ioc := NewIoContext()
	timer := NewDeadlineTimer(ioc)

	timer.ExpiresAt(time.Now().Add(time.Second))

	x := 1
	timer.AsyncWait(func(err error) {
		x = 2
		timer.ExpiresAt(time.Now().Add(time.Second))
		timer.AsyncWait(func(err error) {
			assert.Equal(t, 2, x)
			x = 3
			ioc.Stop()
		})
	})

	ioc.Run()
	assert.Equal(t, 3, x)
}

func TestDeadlineTimer_Cancel(t *testing.T) {
	ioc := NewIoContext()
	timer, stop := NewDeadlineTimer(ioc), NewDeadlineTimer(ioc)
	done := false

	timer.ExpiresFromNow(time.Millisecond)
	timer.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		done = true
	})

	stop.ExpiresFromNow(time.Millisecond * 5)
	stop.AsyncWait(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
		ioc.Stop()
	})

	timer.Cancel()
	ioc.Run()

	assert.Equal(t, false, done)
}

func TestDeadlineTimer_Memory(t *testing.T) {
	//memConsumed := func() uint64 {
	//	var memStat runtime.MemStats
	//	runtime.ReadMemStats(&memStat)
	//	return memStat.Sys
	//}
	//
	//go func() {
	//	http.ListenAndServe("0.0.0.0:8080", nil)
	//}()
	//
	//before := memConsumed()
	//
	//ioc := NewIoContext()
	//timer := NewDeadlineTimer(ioc)
	//
	//var loop func()
	//loop = func() {
	//	timer.Cancel()
	//	timer.ExpiresFromNow(time.Microsecond)
	//	timer.AsyncWait(func(err error) {
	//		after := memConsumed()
	//		fmt.Printf("%.3f ram\n", float64(after-before)/1e3)
	//		loop()
	//	})
	//}
	//
	//loop()
	//
	//ioc.Run()
}
