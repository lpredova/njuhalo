package builder

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

// BaseURL is constant part of url bound to domain
const BaseURL string = "http://www.njuskalo.hr"

// PathURL changable path
var pathURL string

// QueryURL query params of GET request
var queryURL string

// SetMainLocation merhod sets path file of url
func SetMainLocation(searchPath string) {
	pathURL = searchPath
}

// SetFilters method adds user defined filters to url as GET param
func SetFilters(filters map[string]string) {
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

// GetDoc method that makes HTTP requests and saves response as NewDocument
func GetDoc() *goquery.Document {
	url := BaseURL + pathURL + queryURL
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic("Cannot get doc")
	}

	return doc
}
