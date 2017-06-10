package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lpredova/njuhalo/builder"
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
