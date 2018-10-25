package test

import (
	"os"
	"gopkg.in/urfave/cli.v1"
	"testing"
	"fmt"
	"sort"
	"log"
)



func makeArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

func TestFlag(t *testing.T) {
	makeArguments("-p","9000","-d","chian")

	app := cli.NewApp()
	app.Name = "GoTest"
	app.Usage = "hello world"
	app.Version = "1.2.3"
	app.Flags = []cli.Flag {
		cli.IntFlag{
			Name:  "port, p",
			Value: 8000,
			Usage: "listening port",
		},

		cli.StringFlag{
			Name:  "print-default-config",
			Usage: "Print default configuration template",
		},
		cli.StringFlag{
			Name:  "data-dir,d",
			Usage: "Directory containing program runtime data",
		},
		cli.StringFlag{
			Name:  "config-dir",
			Usage: "Directory containing configuration files such as config.ini",
		},
		cli.StringFlag{
			Name:  "config,c",
			Usage: "Configuration file name relative to config-dir",
		},
		cli.StringFlag{
			Name:  "logconf",
			Usage: "Logging configuration file name/path for library users",
		},
	}
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Print this help message and exit.",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print version information.",
	}


	sort.Sort(cli.FlagsByName(app.Flags))

	cli.HelpFlag = cli.BoolFlag {
		Name: "help, h",
		Usage: "Print this help message and exit.",
	}


	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		fmt.Println(c.Int("port"))
		fmt.Println(c.String("-d") == "")
		fmt.Println(c.Command)
		// if c.Int("port") == 8000 {
		// 	return cli.NewExitError("invalid port", 88)
		// }

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}






