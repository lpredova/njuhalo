package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lpredova/njuhalo/helper"
	"github.com/lpredova/njuhalo/model"
	_ "github.com/mattn/go-sqlite3" // SQLlite db
)

const dbPath = "./storage/njuhalo.db"

// InsertItem method inserts new offer into database
func InsertItem(offers []model.Offer) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO items(queryID, itemID, isNew, url, name, image, price, description, location, year, mileage, published, createdAt) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, offer := range offers {
		_, err := stmt.Exec(offer.QueryID, offer.ItemID, 1, offer.URL, offer.Name, offer.Image, offer.Price, offer.Description, offer.Location, offer.Year, offer.Mileage, offer.Published, int32(time.Now().Unix()))
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

// GetDashboardItems gets all items stored in database
func GetDashboardItems() (*[]model.Offer, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT id, isNew, url, name, image, price, description, location, year, mileage, published, createdAt FROM items WHERE isNew = %d ORDER BY isNew DESC, id ASC;", 1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []model.Offer{}
	for rows.Next() {
		offer := model.Offer{}
		rows.Scan(&offer.ID, &offer.IsNew, &offer.URL, &offer.Name, &offer.Image, &offer.Price, &offer.Description, &offer.Location, &offer.Year, &offer.Mileage, &offer.Published, &offer.CreatedAt)
		offers = append(offers, offer)
	}

	return &offers, nil
}

// GetQueryItems gets all query related items stored in the database
func GetQueryItems(queryID int64) (*[]model.Offer, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT id, isNew, url, name, image, price, description, location, year, mileage, published, createdAt FROM items WHERE queryId= %d ORDER BY isNew DESC, id ASC;", queryID)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []model.Offer{}
	for rows.Next() {
		offer := model.Offer{}
		rows.Scan(&offer.ID, &offer.IsNew, &offer.URL, &offer.Name, &offer.Image, &offer.Price, &offer.Description, &offer.Location, &offer.Year, &offer.Mileage, &offer.Published, &offer.CreatedAt)
		offers = append(offers, offer)
	}

	return &offers, nil
}

// GetItem method that checks if there is alreay offer with that ID
func GetItem(itemID string) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	var count int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(id) FROM items WHERE isNew=%d AND itemID=%s", 1, itemID)).Scan(&count)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if count != 0 {
		_, err = db.Exec(fmt.Sprintf("UPDATE items SET isNew=%d WHERE itemID = %s", 0, itemID))
		return true
	}

	return false
}

// InsertQuery saves new query
func InsertQuery(query model.Query) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO queries(name, isActive, url, filters, createdAt) values(?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(query.Name, "Y", query.URL, query.Filters, int32(time.Now().Unix()))
	if err != nil {
		return err
	}

	return nil
}

// GetQueries returns all queries saved in db
func GetQueries() (*[]model.Query, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, url, isActive, filters, createdAt FROM queries")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	queries := []model.Query{}
	for rows.Next() {
		query := model.Query{}
		rows.Scan(&query.ID, &query.Name, &query.URL, &query.IsActive, &query.Filters, &query.CreatedAt)
		queries = append(queries, query)
	}

	return &queries, nil
}

// CreateDatabase creates sqllite db file in user home dir
func CreateDatabase() bool {

	err := os.MkdirAll("./storage", 0755)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if _, err = os.Stat("./storage"); os.IsNotExist(err) {
		os.Mkdir("./storage", 0755)
	}
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	f, err := os.Create(dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer f.Close()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE items (id integer PRIMARY KEY AUTOINCREMENT, itemID integer, url text, name text, image text, price text, description text, createdAt integer)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	stmt, err = db.Prepare("CREATE INDEX index_itemID ON items (itemID)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// SaveQuery method saves query url to config
func SaveQuery(query string) error {
	if len(query) == 0 {
		return errors.New("Please provide valid njuskalo.hr URL")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return err
	}

	random := helper.RandomString()
	req.Header.Set("User-Agent", random)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		u, err := url.Parse(query)
		if err != nil {
			return errors.New("Error parsing URL")
		}

		if u.Host == "www.njuskalo.hr" {
			parsed, _ := url.ParseQuery(u.RawQuery)
			rawFilters := make(map[string]string)
			for k, v := range parsed {
				rawFilters[k] = strings.Join(v, "")
			}

			filters, err := json.Marshal(rawFilters)

			query := model.Query{
				Name:    u.Path,
				URL:     u.Path,
				Filters: string(filters),
			}

			err = InsertQuery(query)
			if err == nil {
				return nil
			}

			return err
		}
		return errors.New("Given url is not from njuskalo")
	}
	return errors.New("This URL is not alive")
}

// DeleteQuery method saves query url to config
func DeleteQuery(queryID int64) error {
	db, err := sql.Open("sqlite3", dbPath)
	db.Exec("PRAGMA foreign_keys = ON;")

	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DELETE FROM queries where id = %d", queryID))
	if err != nil {
		return err
	}

	return nil
}
