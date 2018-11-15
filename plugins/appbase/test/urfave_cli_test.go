package test

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"testing"
	main "github.com/eosspark/eos-go/plugins/producer_plugin/testing"
	"github.com/stretchr/testify/assert"
)

func TestUrfaveCli(t *testing.T) {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		name := "Nefertiti"
		if c.NArg() > 0 {
			name = c.Args().Get(0)
		}
		if c.String("lang") == "spanish" {
			fmt.Println("Hola", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil
	}

	app.Run(os.Args)
}

func Test_commend(t *testing.T) {

	main.MakeTesterArguments("-e", "-n", "eos", "-c", "18", "-p", "eosio", "-p", "yc")

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "enable, e",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Enable block production, even if the chain is stale.",
		},
		cli.IntFlag{
			Name:  "count, c",
			Usage: "Enable block production, even if the chain is stale.",
			Value: -1,
		},
		cli.StringSliceFlag{
			Name:  "producers, p",
			Usage: "Enable block production, even if the chain is stale.",
		},
	}

	app.Action = func(c *cli.Context) {
		assert.Equal(t, true, c.Bool("enable"))
		assert.Equal(t, "eos", c.String("name"))
		assert.Equal(t, 18, c.Int("count"))

	}

	app.Action = func(c *cli.Context) {
		assert.Equal(t, "eosio", c.StringSlice("producers")[0])
		assert.Equal(t, "yc", c.StringSlice("producers")[1])
	}

	err := app.Run(os.Args)
	if err != nil {
		t.Fatal(err)
	}

}
