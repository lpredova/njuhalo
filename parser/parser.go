package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lpredova/goquery"
	"github.com/lpredova/njuhalo/builder"
	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/helper"
	"github.com/lpredova/njuhalo/model"
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

		publishedElement := s.Find(".entity-pub-date")
		published := publishedElement.Text()

		descriptionElement := s.Find(".entity-description-main")
		description := descriptionElement.Text()
		if len(itemID) == 0 {
			return
		}
		f := strings.Fields(description)

		offers = append(offers, model.Offer{
			ID:          itemID,
			URL:         itemLink,
			Name:        itemTitle,
			Image:       image,
			Price:       price,
			Year:        helper.GetNumber(helper.GetSliceData(f, 5)),
			Location:    f[len(f)-1],
			Mileage:     helper.GetNumber(helper.GetSliceData(f, 2)),
			Published:   published,
			Description: description,
		})
	})

	return offers
}

func checkPagination(doc *goquery.Document) bool {
	hasPagination := false
	doc.Find("nav.Pagination").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("a").Text(), "Sljedeća") {
			hasPagination = true
		}
	})

	return hasPagination
}

// ParseOffer parses one entity
func ParseOffer(doc *goquery.Document) (bool, []model.Offer) {
	var index int
	var offer model.Offer
	var offers []model.Offer
	var finalOffers []model.Offer

	offers = GetListContent(doc, ".EntityList--VauVau .EntityList-item article", offers)
	offers = GetListContent(doc, ".EntityList--Standard .EntityList-item article", offers)

	if len(offers) == 0 {
		fmt.Println("No new offers found!")
	}

	for index, offer = range offers {
		if !db.GetItem(offer.ID) {
			finalOffers = append(finalOffers, offer)
		}
	}
	fmt.Println(fmt.Sprintf("Parsed %d new results", index))

	return db.InsertItem(finalOffers), finalOffers
}

// GetNextResultPage tries to determine if there are more pages?
// if there are then get them and parse
func GetNextResultPage(doc *goquery.Document, page int, filters map[string]string) (bool, int, map[string]string) {
	if !checkPagination(doc) {
		return false, 0, filters
	}

	page++
	if filters == nil {
		filters = make(map[string]string)
	}

	filters["page"] = strconv.Itoa(page)
	return true, page, filters
}
