package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lpredova/shnjuskhalo/builder"
	"github.com/lpredova/shnjuskhalo/model"
)

// GetListContent method gets items for sale and parses them
func GetListContent(doc *goquery.Document, selector string, offers []model.Offer) []model.Offer {

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a")

		itemID, _ := titleElement.Attr("name")
		itemTitle := titleElement.Text()

		itemLink, _ := titleElement.Attr("href")
		itemLink = fmt.Sprintf("%s%s", builder.BaseURL, itemLink)

		offers = append(offers, model.Offer{
			Name: itemTitle,
			URL:  itemLink,
			ID:   itemID,
		})
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
