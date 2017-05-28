package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "http://www.njuskalo.hr/"

var pathURL string
var queryURL string

var filters map[string]string
var page = 0
var numOfItems = 0

func main() {
	setMainLocation("iznajmljivanje-stanova", "zagreb")

	filters = make(map[string]string)
	filters["locationId"] = "2619"
	filters["price[max]"] = "260"
	filters["mainArea[max]"] = "50"
	setFilters(filters)

	doc := getDoc()
	parseOffer(doc)

	for {
		if !checkForMore(doc) {
			break
		}

		parseOffer(doc)
	}
}

// appendFilters method adds user defined filters to url as GET param
func setFilters(filters map[string]string) {
	filtersNumber := len(filters)
	if filtersNumber == 0 {
		return
	}

	var i int
	queryURL = ""
	for key, value := range filters {
		queryURL += fmt.Sprintf("%s=%s", string(key), string(value))
		i++

		if i < filtersNumber {
			queryURL += "&"
		}
	}

	queryURL = "?" + queryURL
}

func getDoc() *goquery.Document {
	url := baseURL + pathURL + queryURL
	fmt.Println(url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic("Cannot get doc")
	}

	return doc
}

func setMainLocation(category string, param string) {
	pathURL = fmt.Sprintf("%s/%s", category, param)
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

		numOfItems++

		fmt.Printf("\n\nReview %d:\nID:%s\n%s\n%s", numOfItems, itemID, itemTitle, itemLink)
	})
}

func checkForMore(doc *goquery.Document) bool {
	// try to see if there are more pages?
	// if there are then get them and parse
	if checkPagination(doc) {
		page++

		fmt.Println("\nThere is more...")
		time.Sleep(time.Second * 3)
		fmt.Println(fmt.Sprintf("\nGetting page %d", page))

		filters["page"] = strconv.Itoa(page)
		setFilters(filters)

		getDoc()
		return true
	}

	return false
}

func checkPagination(doc *goquery.Document) bool {

	hasPagination := false
	doc.Find(" div.entity-list-pagination").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a").Text()
		if strings.Contains(titleElement, "Sljedeća") {
			hasPagination = true
		}
	})

	return hasPagination
}

// getItemContent is method that provides specific item info
func getItemContent() {

}
