package main

import (
	"os"

	"github.com/lpredova/njuhalo/command"
	"github.com/urfave/cli"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "njuhalo"
	app.Usage = "Monitor Njuskalo as a PRO"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, con",
			Value: "$HOME/njhalo.json",
			Usage: "PATH to config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "initialize config file",
			Action: func(c *cli.Context) error {
				command.CreateConfigFile()
				return nil
			},
		},
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start monitoring njuskalo for items",
			Action: func(c *cli.Context) error {
				command.StartMonitoring()
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "adds query to watch to config",
			Action: func(c *cli.Context) error {
				command.SaveQuery(c.Args().Get(0))
				return nil
			},
		},
	}
}

func main() {
	app.Run(os.Args)
}
