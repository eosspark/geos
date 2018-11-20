package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"sort"
)

func main() {
	//done := make(chan bool)
	app := cli.NewApp()
	app.Commands = []cli.Command{
		walletCommand,
		accountCommand,
		getCommand,
		SignCommand,
		consoleCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	// app.Before = func(c *cli.Context) error {
	// 	fmt.Fprintf(c.App.Writer, "EOSIO Client \n")
	// 	return nil
	// }
	// app.After = func(c *cli.Context) error {
	// 	fmt.Fprintf(c.App.Writer, "finish EOSIO Client\n")
	// 	return nil
	// }

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//<-done

}
