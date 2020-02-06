package command

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/goquery"
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

// StartMonitoring starts watcher that monitors items
func StartMonitoring() {
	conf = configuration.ParseConfig()
	if conf.RunIntervalMin > 0 {

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

			/*for {
				if checkForMore(doc) {
					parseOffer(doc)
				}
			}*/
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

		if filters == nil {
			filters = make(map[string]string)
		}

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

	if len(offers) == 0 {
		fmt.Println("No new offers found!")
	}

	for index, offer := range offers {
		if !db.GetItem(offer.ID) {
			finalOffers = append(finalOffers, offer)
			fmt.Println(fmt.Sprintf("%d. %s - (%s) %s ", index, offer.Name, offer.Price, offer.URL))
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

// ClearQueries removes all queries from config
func ClearQueries() {
	if configuration.ClearQueries() {
		fmt.Println("Queries cleared")
	} else {
		fmt.Println("Error clearing queries")
	}
}

// SaveQuery method saves query url to config
func SaveQuery(query string) {
	if len(query) > 0 {

		nj := randomString()
		client := &http.Client{}
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			fmt.Println("Unable to create request")
		}

		req.Header.Set("User-Agent", nj)
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
	} else {
		fmt.Println("Please provide valid njuskalo.hr URL")
	}
}

func randomString() string {
	n := rand.Intn(20)
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}
