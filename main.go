package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

const baseURL string = "http://www.njuskalo.hr/"

func main() {
	getLocation("iznajmljivanje-stanova", "zagreb")
}

func getLocation(category string, param string) {

	doc, err := goquery.NewDocument(fmt.Sprintf("%s%s/%s", baseURL, category, param))
	if err != nil {
		log.Fatal(err)
	}

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
