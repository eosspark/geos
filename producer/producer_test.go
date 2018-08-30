package producer_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"testing"
	"time"
)

const KEEPTESTSEC = 99910 /*seconds*/

func Test_Timer(t *testing.T) {
	start := time.Now()
	var apply = false
	var timer = new(scheduleTimer)
	var blockNum = 1

	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGINT)

	var scheduleProductionLoop func()

	scheduleProductionLoop = func() {
		timer.cancel()
		base := time.Now()
		minTimeToNextBlock := int64(common.DefaultConfig.BlockIntervalUs) - base.UnixNano()/1e3%int64(common.DefaultConfig.BlockIntervalUs)
		wakeTime := base.Add(time.Microsecond * time.Duration(minTimeToNextBlock))

		timer.expiresUntil(wakeTime)

		// test after 12 block need to apply new block to continue
		if blockNum%10 == 0 || (blockNum-1)%10 == 0 || (blockNum-2)%10 == 0 {
			apply = true
			return
		}

		timerCorelationId++
		cid := timerCorelationId
		timer.asyncWait(func() bool { return cid == timerCorelationId }, func() {
			fmt.Println("exec async1...", time.Now())
			blockNum++
			fmt.Println("add.blockNum", blockNum)
			scheduleProductionLoop()
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

func Test_producer_start(t *testing.T) {
	start := time.Now()
	os.Args = []string{"--enable-stale-production", "-p", "eosio", "-p", "yuanc"}

	//for i:=0; i< len(os.Args); i++ {
	//	fmt.Println(os.Args[i])
	//}

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	//app.Commands = []cli.Command{
	//	{
	//		Name:     "add",
	//		Aliases:  []string{"a"},
	//		Usage:    "calc 1+1",
	//		Category: "arithmetic",
	//		Action: func(c *cli.Context) error {
	//			fmt.Println("1 + 1 = ", 1+1)
	//			return nil
	//		},
	//	},
	//	{
	//		Name:     "sub",
	//		Aliases:  []string{"s"},
	//		Usage:    "calc 5-3",
	//		Category: "arithmetic",
	//		Action: func(c *cli.Context) error {
	//			fmt.Println("5 - 3 = ", 5-3)
	//			return nil
	//		},
	//	},
	//	{
	//		Name:     "db",
	//		Usage:    "database operations",
	//		Category: "database",
	//		Subcommands: []cli.Command{
	//			{
	//				Name:  "insert",
	//				Usage: "insert data",
	//				Action: func(c *cli.Context) error {
	//					fmt.Println("insert subcommand")
	//					return nil
	//				},
	//			},
	//			{
	//				Name:  "delete",
	//				Usage: "delete data",
	//				Action: func(c *cli.Context) error {
	//					fmt.Println("delete subcommand")
	//					return nil
	//				},
	//			},
	//		},
	//	},
	//}

	produce := new(ProducerPlugin)
	produce.Initialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	produce.Startup()

	for {
		if time.Now().Sub(start) > KEEPTESTSEC*time.Second {
			produce.Shutdown()
			break
		}
	}
}
