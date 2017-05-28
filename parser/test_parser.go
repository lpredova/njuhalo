package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

var doc *goquery.Document

func TestGetListContent(t *testing.T) {
	doc = loadStaticDoc("testListPage.html")
	GetListContent(doc, ".EntityList--VauVau .EntityList-item article .entity-title")
}

func loadStaticDoc(page string) *goquery.Document {
	var f *os.File
	var e error

	if f, e = os.Open(fmt.Sprintf("./test_data/%s", page)); e != nil {
		panic(e.Error())
	}
	defer f.Close()

	var node *html.Node
	if node, e = html.Parse(f); e != nil {
		panic(e.Error())
	}
	return goquery.NewDocumentFromNode(node)
}
