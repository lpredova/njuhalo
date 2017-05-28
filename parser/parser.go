package parser

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lpredova/shnjuskhalo/builder"
)

var numOfItems = 0

// GetListContent method gets items for sale and parses them
func GetListContent(doc *goquery.Document, selector string) {
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a")

		itemID, _ := titleElement.Attr("name")
		itemTitle := titleElement.Text()

		itemLink, _ := titleElement.Attr("href")
		itemLink = fmt.Sprintf("%s%s", builder.BaseURL, itemLink)

		numOfItems++

		fmt.Printf("\n\nReview %d:\nID:%s\n%s\n%s", numOfItems, itemID, itemTitle, itemLink)
	})
}

// CheckPagination method checks if there is pagination element on html
func CheckPagination(doc *goquery.Document) bool {
	hasPagination := false
	doc.Find(" div.entity-list-pagination").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a").Text()
		if strings.Contains(titleElement, "Sljedeća") {
			hasPagination = true
		}
	})

	return hasPagination
}
