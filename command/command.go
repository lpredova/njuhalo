package command

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/goquery"
	"github.com/lpredova/njuhalo/builder"
	"github.com/lpredova/njuhalo/configuration"
	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/helper"
	"github.com/lpredova/njuhalo/model"
	"github.com/lpredova/njuhalo/parser"
)

var page = 0
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

// PrintConfigFile prints currently active config file for user to see
func PrintConfigFile() {
	configuration.PrintConfig()
}

// ListItems lists all items in database
func ListItems() {
	offers, err := db.GetItems()
	if err != nil {
		fmt.Println(err.Error())
	}

	if len(*offers) == 0 {
		fmt.Println("No offers saved yet :)")
	}

	for index, offer := range *offers {
		fmt.Println(fmt.Sprintf("%d. %s - (%s) %s ", index, offer.Name, offer.Price, offer.URL))
	}
}

// Parse runs parser only once :)
func Parse() {
	conf = configuration.ParseConfig()
	runParser()
}

// Monitor starts watcher that monitors items
func Monitor() {
	conf = configuration.ParseConfig()
	if conf.RunIntervalMin <= 0 {
		fmt.Println("Please provide valid watcher run interval (larger than 0)")
		return
	}

	runParser()
	gocron.Every(uint64(conf.RunIntervalMin)).Minutes().Do(runParser)
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

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save-query", saveQueryHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func saveQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("Error")
		http.Redirect(w, r, "/", 301)
	}

	r.ParseForm()
	SaveQuery(r.Form["query"])

	http.Redirect(w, r, "/", 301)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	offers, err := db.GetItems()
	if err != nil {
		fmt.Println(err.Error())
	}

	type viewData struct {
		Offers *[]model.Offer
	}

	data := viewData{
		Offers: offers,
	}

	template.Must(
		template.New("").
			Funcs(template.FuncMap{"ts2date": helper.TimestampToDate}).
			ParseFiles("templates/index.tmpl"),
	).ExecuteTemplate(w, "index.tmpl", data)
}

func runParser() {
	if len(conf.Queries) <= 0 {
		fmt.Println("There are no filters in your config, please check help")
		return
	}

	for _, query := range conf.Queries {
		builder.SetMainLocation(query.BaseQueryPath)
		builder.SetFilters(query.Filters)

		doc := builder.GetDoc()
		parser.ParseOffer(doc)

		for {
			if parser.CheckForMore(doc) {
				parser.ParseOffer(doc)
			}
		}
	}
}

// ClearQueries removes all queries from config
func ClearQueries() {
	if configuration.ClearQueries() {
		fmt.Println("Queries cleared")
		return
	}

	fmt.Println("Error clearing queries")
}

// SaveQuery method saves query url to config
func SaveQuery(query string) {
	if len(query) == 0 {
		fmt.Println("Please provide valid njuskalo.hr URL")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		fmt.Println("Unable to create request")
	}

	req.Header.Set("User-Agent", helper.RandomString())
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error checking URL")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		u, err := url.Parse(query)
		if err != nil {
			fmt.Println("Error parsing URL")
		}

		if u.Host == "www.njuskalo.hr" {
			parsed, _ := url.ParseQuery(u.RawQuery)
			rawFilters := make(map[string]string)
			for k, v := range parsed {
				rawFilters[k] = strings.Join(v, "")
			}

			query := model.Query{
				BaseQueryPath: u.Path,
				Filters:       rawFilters,
			}

			if configuration.AppendFilterToConfig(query) {
				fmt.Println("URL added to filters")
			} else {
				fmt.Println("Error adding URL to filters")
			}
		} else {
			fmt.Println("Given url is not from njuskalo")
		}
	} else {
		fmt.Println("This URL is not alive")
	}
}
