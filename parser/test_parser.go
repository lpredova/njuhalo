package parser

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/html"
)

var doc *goquery.Document

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
