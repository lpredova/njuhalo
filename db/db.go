package db

import (
	"database/sql"
	"fmt"

	"github.com/lpredova/shnjuskhalo/model"
	_ "github.com/mattn/go-sqlite3" // SQLlite db
)

const dbName = "./njuhalo.db"

// InsertItem method inserts new offer into database
func InsertItem(offers []model.Offer) bool {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO items(itemID, url, name, image, price, description) values(?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, offer := range offers {
		_, err := stmt.Exec(offer.ID, offer.URL, offer.Name, offer.Image, offer.Price, offer.Description)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

// GetItem method that checks if there is alreay offer with that ID
func GetItem(itemID string) bool {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM items where itemID = %s", itemID))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for rows.Next() {
		return true
	}

	fmt.Println("Item WILL BE ADDED")
	return false
}
