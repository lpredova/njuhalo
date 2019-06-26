package command

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/goquery"
	"github.com/lpredova/njuhalo/configuration"
	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/handler"
	"github.com/lpredova/njuhalo/model"
	"github.com/lpredova/njuhalo/parser"
)

var doc *goquery.Document
var conf model.Configuration
var filters map[string]string

// CreateConfigFile method crates boilerplate config file
func CreateConfigFile() {

	if db.CreateDatabase() {
		fmt.Println("Database created")
	} else {
		fmt.Println("Error creating database")
	}

	config := model.Configuration{}
	if configuration.CreateFileConfig(config) {
		fmt.Println("Config file created")
	} else {
		fmt.Println("Error creating config file")
	}
}

// Parse runs parser only once :)
func Parse() {
	parser.Run()
}

// Monitor starts watcher that monitors items
func Monitor() {
	if conf.RunIntervalMin <= 0 {
		fmt.Println("Please provide valid watcher run interval (larger than 0)")
		return
	}

	gocron.Every(uint64(conf.RunIntervalMin)).Minutes().Do(parser.Run())
	<-gocron.Start()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		gocron.Clear()
		os.Exit(1)
	}()
}

// Serve for listing results in browser
func Serve() {
	fmt.Println("Serving results: http://localhost:8080")
	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/oops", handler.ErrorHandler)
	http.HandleFunc("/save-query", handler.SaveQueryHandler)
	http.HandleFunc("/fetch-results", handler.FetchResultsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
