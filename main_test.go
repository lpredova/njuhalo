package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

var doc *goquery.Document

func TestMain(t *testing.T) {
	main()
}

func TestGetLocation(t *testing.T) {
	getLocation("iznajmljivanje-stanova", "zagreb")
}

func TestGetNotExistingLocation(t *testing.T) {
	getLocation("", "zagreb123")
}

func TestGetListContent(t *testing.T) {
	doc = loadDoc("testListPage.html")
	getListContent(doc, ".EntityList--VauVau .EntityList-item article .entity-title")
}

func TestGetContent(t *testing.T) {
	getItemContent()
}

func loadDoc(page string) *goquery.Document {
	var f *os.File
	var e error

	if f, e = os.Open(fmt.Sprintf("./testData/%s", page)); e != nil {
		panic(e.Error())
	}
	defer f.Close()

	var node *html.Node
	if node, e = html.Parse(f); e != nil {
		panic(e.Error())
	}
	return goquery.NewDocumentFromNode(node)
}
