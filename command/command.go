package command

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/njuhalo/alert"
	"github.com/lpredova/njuhalo/builder"
	"github.com/lpredova/njuhalo/configuration"
	"github.com/lpredova/njuhalo/db"
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

// StartMonitoring starts watcher that monitors items
func StartMonitoring() {
	conf = configuration.ParseConfig()

	if conf.RunIntervalMin > 0 {
		runParser()
		gocron.Every(uint64(conf.RunIntervalMin)).Minute().Do(runParser)
		<-gocron.Start()
	} else {
		fmt.Println("Please provide valid watcher run interval (larger than 0)")
	}
}

func runParser() {
	if len(conf.Queries) > 0 {
		for _, query := range conf.Queries {
			builder.SetMainLocation(query.BaseQueryPath)
			builder.SetFilters(query.Filters)

			doc := builder.GetDoc()
			parseOffer(doc)

			for {
				if !checkForMore(doc) {
					break
				}

				parseOffer(doc)
			}
		}
	} else {
		fmt.Println("There are no filters in your config, please check help")
	}
}

// try to see if there are more pages?
// if there are then get them and parse
func checkForMore(doc *goquery.Document) bool {
	if parser.CheckPagination(doc) {
		page++
		time.Sleep(time.Second * time.Duration(int(conf.SleepIntervalSec)))
		filters["page"] = strconv.Itoa(page)
		builder.SetFilters(filters)
		builder.GetDoc()
		return true
	}

	return false
}

func parseOffer(doc *goquery.Document) {
	var offers []model.Offer
	var finalOffers []model.Offer

	offers = parser.GetListContent(doc, ".EntityList--VauVau .EntityList-item article", offers)
	offers = parser.GetListContent(doc, ".EntityList--Standard .EntityList-item article", offers)

	for _, offer := range offers {
		if !db.GetItem(offer.ID) {
			finalOffers = append(finalOffers, offer)
		}
	}

	if db.InsertItem(finalOffers) {
		if conf.Slack {
			alert.SendItemsToSlack(finalOffers)
		}

		if conf.Mail {
			alert.SendItemsToMail(finalOffers)
		}
	}
}

// SaveQuery method saves query url to config
func SaveQuery(query string) {
	if len(query) > 0 {
		resp, err := http.Get(query)
		if err != nil {
			fmt.Println("Error checking URL")
		}

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
	} else {
		fmt.Println("Please provide valid njuskalo.hr URL")
	}
}
