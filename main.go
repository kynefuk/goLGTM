package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/kynefuk/goLGTM/command"
	"github.com/urfave/cli/v2"
)

// Command's ExitCode
const (
	ExitCodeOk = iota
	ExitCodeErr
)

func main() {
	app := cli.NewApp()
	app.Name = "GoLGTM"
	app.Usage = "make a LGTM picture."

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "image source",
		},
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Value:   "LGTM",
			Usage:   "message",
		},
	}

	app.Action = func(c *cli.Context) error {
		source := c.String("source")
		message := c.String("message")

		command := command.NewCommand()
		return command.Run(source, message)
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(ExitCodeErr)
	}
	os.Exit(ExitCodeOk)
}
