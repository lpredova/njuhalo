package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "http://www.njuskalo.hr/"

var url string

func main() {
	setMainLocation("iznajmljivanje-stanova", "zagreb")

	filters := make(map[string]string)
	filters["locationId"] = "2619"
	filters["price[max]"] = "260"
	filters["mainArea[max]"] = "50"
	setFilters(filters)

	doc := getDoc()
	parseOffer(doc)
}

// appendFilters method adds user defined filters to url as GET param
func setFilters(filters map[string]string) {
	filtersNumber := len(filters)
	if filtersNumber == 0 {
		return
	}

	var fil string
	var i int
	for key, value := range filters {
		fil += fmt.Sprintf("%s=%s", key, value)
		i++

		if i < filtersNumber {
			fil += "&"
		}
	}

	url = url + "?" + fil
}

func getDoc() *goquery.Document {
	url = baseURL + url
	fmt.Println(url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic("Cannot get doc")
	}

	return doc
}

func setMainLocation(category string, param string) {
	url = fmt.Sprintf("%s/%s", category, param)
}

func parseOffer(doc *goquery.Document) {
	fmt.Println("Vau Vau offer")
	getListContent(doc, ".EntityList--VauVau .EntityList-item article .entity-title")

	fmt.Println("Regular offer")
	getListContent(doc, ".EntityList--Standard .EntityList-item article .entity-title")
}

func getListContent(doc *goquery.Document, selector string) {

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a")

		itemID, _ := titleElement.Attr("name")
		itemTitle := titleElement.Text()

		itemLink, _ := titleElement.Attr("href")
		itemLink = fmt.Sprintf("%s%s", baseURL, itemLink)

		fmt.Printf("Review %d: ID:%s - %s - %s\n", i+1, itemID, itemTitle, itemLink)
	})
}

// getItemContent is method that provides specific item info
func getItemContent() {

}
