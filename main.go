package main

import (
	"os"

	"github.com/lpredova/njuhalo/command"
	"github.com/urfave/cli"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Version = "2.0.0"
	app.Name = "njuhalo"
	app.Author = "Lovro Predovan"
	app.Email = "lovro.predovan[at]gmail.com"
	app.Usage = "Watcher for Njuskalo.hr offers"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, con",
			Value: "storage/config.json",
			Usage: "PATH to config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "init, initialize",
			Aliases: []string{"i"},
			Usage:   "initialize configuration and database file in home dir",
			Action: func(c *cli.Context) error {
				command.CreateConfigFile()
				return nil
			},
		},
		{
			Name:    "monitor",
			Aliases: []string{"m"},
			Usage:   "Start monitoring njuskalo for items",
			Action: func(c *cli.Context) error {
				command.Monitor()
				return nil
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "Start local server",
			Action: func(c *cli.Context) error {
				command.Serve()
				return nil
			},
		},
		{
			Name:    "parse",
			Aliases: []string{"r"},
			Usage:   "Runs parser only once",
			Action: func(c *cli.Context) error {
				command.Parse()
				return nil
			},
		},
	}
}

func main() {
	app.Run(os.Args)
}
