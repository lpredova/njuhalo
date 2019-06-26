package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/lpredova/njuhalo/db"
	"github.com/lpredova/njuhalo/helper"
	"github.com/lpredova/njuhalo/model"
	"github.com/lpredova/njuhalo/parser"
)

// IndexHandler serves up main dashboard of the app
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	var offers *[]model.Offer
	var err error

	offerID, hasParam := r.URL.Query()["offerId"]
	if hasParam {
		offerID, _ := strconv.ParseInt(offerID[0], 10, 64)
		offers, err = db.GetQueryItems(offerID)
	} else {
		offers, err = db.GetDashboardItems()
	}

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

// DeleteQueryHandler deletes existing query
func DeleteQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}

	r.ParseForm()
	query := strings.Join(r.Form["queryId"], " ")
	if query == "" {
		http.Redirect(w, r, "/", 301)
		return
	}

	queryID, _ := strconv.ParseInt(query, 10, 64)
	err := db.DeleteQuery(queryID)
	if err == nil {
		parser.Run()
		http.Redirect(w, r, "/", 301)
		return
	}

	fmt.Fprint(w, err.Error())
	http.Redirect(w, r, "/", 301)
}
