package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lpredova/goquery"
	"github.com/lpredova/njuhalo/alert"
	"github.com/lpredova/njuhalo/builder"
	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/model"
	"github.com/lpredova/njuhalo/parser"
)

// GetListContent method gets items for sale and parses them
func GetListContent(doc *goquery.Document, selector string, offers []model.Offer) []model.Offer {

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find(".entity-title a")

		itemID, _ := titleElement.Attr("name")
		itemTitle := titleElement.Text()

		itemLink, _ := titleElement.Attr("href")
		itemLink = fmt.Sprintf("%s%s", strings.TrimSuffix(builder.BaseURL, "/"), itemLink)

		imageElement := s.Find(".entity-thumbnail img")
		image, _ := imageElement.Attr("data-src")
		image = fmt.Sprintf("%s%s", "http:", image)

		priceElement := s.Find(".entity-prices .price-item .price--eur")
		price := priceElement.Text()

		descriptionElement := s.Find(".entity-description-main")
		description := descriptionElement.Text()

		if len(itemID) > 0 {
			offers = append(offers, model.Offer{
				ID:          itemID,
				URL:         itemLink,
				Name:        itemTitle,
				Image:       image,
				Price:       price,
				Description: description,
			})
		}
	})

	return offers
}

// CheckPagination method checks if there is pagination element on html
func CheckPagination(doc *goquery.Document) bool {
	hasPagination := false
	doc.Find("div.entity-list-pagination").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("a").Text(), "Sljedeća") {
			hasPagination = true
		}
	})

	return hasPagination
}

// ParseOffer parses one entity
func ParseOffer(doc *goquery.Document) {
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

// CheckForMore tries to determine if there are more pages?
// if there are then get them and parse
func CheckForMore(doc *goquery.Document) bool {
	if !parser.CheckPagination(doc) {
		return false
	}

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
