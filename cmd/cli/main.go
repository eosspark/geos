package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		walletCommand,
		accountCommand,
		getCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
