package db

import (
	"database/sql"
	"fmt"
	"time"

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

	stmt, err := db.Prepare("INSERT INTO items(itemID, url, name, image, price, description, createdAt) values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	for _, offer := range offers {
		_, err := stmt.Exec(offer.ID, offer.URL, offer.Name, offer.Image, offer.Price, offer.Description, int32(time.Now().Unix()))
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

	return false
}

// ClearOldItems checks db for items older than 90 days and does not track them anymore
func ClearOldItems() {

	time := int32(time.Now().AddDate(0, 0, -90).Unix())

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM items where createdAt < ?")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec(time)
	if err != nil {
		fmt.Println(err.Error())
	}

}
