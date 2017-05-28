package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"github.com/lpredova/shnjuskhalo/builder"
	"github.com/lpredova/shnjuskhalo/configuration"
	"github.com/lpredova/shnjuskhalo/model"
	"github.com/lpredova/shnjuskhalo/parser"
)

var page = 0
var doc *goquery.Document
var conf model.Configuration
var filters map[string]string

// StartMonitoring starts watcher that monitors items
func StartMonitoring() {
	conf = configuration.ParseConfig()

	gocron.Every(uint64(conf.RunIntervalMin)).Minute().Do(checkItems)
	<-gocron.Start()
	fmt.Println("Started monitoring offers...")
}

// CreateConfigFile method crates boilerplate config file
func CreateConfigFile() {
	if configuration.CreateFileConfig() {
		fmt.Println("Config file created")
	} else {
		fmt.Println("Error creating config file")
	}
}

func checkItems() {
	builder.SetMainLocation("iznajmljivanje-stanova", "zagreb")

	filters = make(map[string]string)
	/*filters["locationId"] = "2619"
	filters["price[max]"] = "260"
	filters["mainArea[max]"] = "50"*/

	builder.SetFilters(filters)
	doc := builder.GetDoc()
	parseOffer(doc)

	for {
		if !checkForMore(doc) {
			break
		}

		parseOffer(doc)
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

	offers = parser.GetListContent(doc, ".EntityList--VauVau .EntityList-item article .entity-title", offers)
	offers = parser.GetListContent(doc, ".EntityList--Standard .EntityList-item article .entity-title", offers)

	for _, offer := range offers {
		fmt.Println(fmt.Sprintf("ID:%s\nURL:%s\nTitle:%s", offer.ID, offer.URL, offer.Name))
	}
}
