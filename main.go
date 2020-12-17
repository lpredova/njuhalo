package main

import (
	"os"

	"github.com/lpredova/njuhalo/command"
	"github.com/urfave/cli"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Version = "1.0.4"
	app.Name = "njuhalo"
	app.Author = "Lovro Predovan"
	app.Email = "lovro.predovan[at]gmail.com"
	app.Usage = "Watcher for Njuskalo.hr items"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, con",
			Value: "$HOME/.njuhalo/config.json",
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
			Name:    "start, serve",
			Aliases: []string{"s"},
			Usage:   "start monitoring njuskalo for items",
			Action: func(c *cli.Context) error {
				command.StartMonitoring()
				return nil
			},
		},
		{
			Name:    "add, query",
			Aliases: []string{"a"},
			Usage:   "adds query to watch to config",
			Action: func(c *cli.Context) error {
				command.SaveQuery(c.Args().Get(0))
				return nil
			},
		},
		{
			Name:    "clean, clear",
			Aliases: []string{"c"},
			Usage:   "clears all query from config",
			Action: func(c *cli.Context) error {
				command.ClearQueries()
				return nil
			},
		},
		{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "Prints currently active config file",
			Action: func(c *cli.Context) error {
				command.PrintConfigFile()
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "lists currently saved items",
			Action: func(c *cli.Context) error {
				command.ListItems()
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
