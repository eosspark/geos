package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"testing"
	"time"
)

var KEEPTESTSEC time.Duration = 10

func Test_Timer(t *testing.T) {
	start := time.Now()
	var apply = false
	var timer = new(Timer)
	var blockNum = 1

	timerCorelationId := 0

	var scheduleProductionLoop func()

	scheduleProductionLoop = func() {
		timer.Cancel()
		base := Now()
		wakeTime := base.AddUs(Milliseconds(500))

		timer.ExpiresUntil(wakeTime)

		// test after 12 block need to apply new block to continue
		if blockNum%10 == 0 || (blockNum-1)%10 == 0 || (blockNum-2)%10 == 0 {
			apply = true
			return
		}

		timerCorelationId++
		cid := timerCorelationId
		timer.AsyncWait(func() {
			if cid == timerCorelationId {
				fmt.Println("exec async1...", time.Now())
				blockNum++
				fmt.Println("add.blockNum", blockNum)
				scheduleProductionLoop()
			}
		})
	}

	applyBlock := func() {
		for time.Now().Sub(start) <= KEEPTESTSEC*time.Second {
			if apply {
				apply = false
				blockNum++
				fmt.Println("exec apply...", time.Now(), "\n-----------apply block #.", blockNum)
				scheduleProductionLoop()
			}
		}
	}

	naughty := func() {
		for time.Now().Sub(start) <= KEEPTESTSEC*time.Second {
			time.Sleep(666 * time.Millisecond)
			scheduleProductionLoop()
		}
	}

	//go func() {
	//	sig := <-sigs
	//	fmt.Println("sig: ", sig)
	//}()

	scheduleProductionLoop()
	applyBlock()
	naughty() //try to break the schedule timer
}

func Test_Timer_Memory(t *testing.T) {
	memConsumed := func() uint64 {
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	before := memConsumed()

	var timer = new(Timer)

	var loop func()
	loop = func() {
		timer.Cancel()
		timer.ExpiresFromNow(1)
		timer.AsyncWait(func() {
			after := memConsumed()
			fmt.Printf("%.3f KB\n", float64(after-before)/1e3)
			loop()
		})
	}

	loop()

	select {}
}

func Test_TimePoint(t *testing.T) {
	fmt.Println(MaxTimePoint())
	fmt.Println(MinTimePoint())
	now := Now()
	fmt.Println(now, now.TimeSinceEpoch())

	fmt.Println(MaxTimePointSec())
	fmt.Println(MinTimePointSec())
}

func Test_FromIsoString(t *testing.T) {
	s := "2006-01-02T15:04:05.500"
	tp, e := FromIsoString(s)
	assert.NoError(t, e)

	tps, err := FromIsoStringSec(s)
	assert.NoError(t, err)

	fmt.Println(tp)
	fmt.Println(tps)
}

func Test_BlockTimestamp(t *testing.T) {}
