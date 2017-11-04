package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/lpredova/njuhalo/model"
	_ "github.com/mattn/go-sqlite3" // SQLlite db
)

const dbName = "./njuhalo.db"

var usr, _ = user.Current()
var homePath = usr.HomeDir + "/.njuhalo/"
var dbPath = homePath + "njuhalo.db"

// InsertItem method inserts new offer into database
func InsertItem(offers []model.Offer) bool {
	db, err := sql.Open("sqlite3", dbPath)
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

// GetItems gets all items stored in database
func GetItems() (*[]model.Offer, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,url,name,image,price,description FROM items ORDER BY id DESC;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []model.Offer{}
	for rows.Next() {
		offer := model.Offer{}
		rows.Scan(&offer.ID, &offer.URL, &offer.Name, &offer.Image, &offer.Price, &offer.Description)
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

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM items where itemID = %s", itemID))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		return true
	}

	return false
}

// CreateDatabase creates sqllite db file in user home dir
func CreateDatabase() bool {

	err := os.MkdirAll(homePath, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if _, err = os.Stat(homePath); os.IsNotExist(err) {
		os.Mkdir(homePath, 0755)
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
