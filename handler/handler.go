package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/helper"
	"github.com/lpredova/njuhalo/model"
	"github.com/lpredova/njuhalo/parser"
)

// IndexHandler serves up main dashboard of the app
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	offers, err := db.GetItems()
	if err != nil {
		fmt.Println(err.Error())
	}

	queries, err := db.GetQueries()
	if err != nil {
		fmt.Println(err.Error())
	}

	type viewData struct {
		Offers  *[]model.Offer
		Queries *[]model.Query
	}

	data := viewData{
		Offers:  offers,
		Queries: queries,
	}

	template.Must(
		template.New("").
			Funcs(template.FuncMap{"ts2date": helper.TimestampToDate}).
			ParseFiles("templates/index.tmpl"),
	).ExecuteTemplate(w, "index.tmpl", data)
}

// ErrorHandler serves default error page
func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(
		template.New("").
			ParseFiles("templates/error.tmpl"),
	).ExecuteTemplate(w, "error.tmpl", nil)
}

// FetchResultsHandler makes request for re-fetching data for parsing
func FetchResultsHandler(w http.ResponseWriter, r *http.Request) {
	err := parser.Run()
	if err == nil {
		http.Redirect(w, r, "/", 301)
		return
	}
	http.Redirect(w, r, "/oops", 301)
}

// SaveQueryHandler handles saving new query
func SaveQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	r.ParseForm()

	queryString := strings.Join(r.Form["query"], " ")
	if queryString == "" {
		http.Redirect(w, r, "/", 301)
		return
	}

	err := db.SaveQuery(queryString)
	if err == nil {
		parser.Run()
		http.Redirect(w, r, "/", 301)
		return
	}

	fmt.Fprint(w, err.Error())
	http.Redirect(w, r, "/", 301)
}
