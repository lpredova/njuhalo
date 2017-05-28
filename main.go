package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	//getItemContent()
	getCategoryContent()
}

// getCategoryContent is method that gets
func getCategoryContent() {
	//http://www.njuskalo.hr/iznajmljivanje-stanova/zagreb

	doc, err := goquery.NewDocument("http://www.njuskalo.hr/iznajmljivanje-stanova/zagreb")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".EntityList-item article .entity-title").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		title := s.Find("a").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})

}

// getItemContent is method that provides specific item info
func getItemContent() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}
