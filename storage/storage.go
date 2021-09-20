package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func StartDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./shortener.db")
	if err != nil {
		panic(err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT PRIMARY KEY, url TEXT)")
	if err != nil {
		panic(err)
	}
	statement.Exec()
	return db
}

func InsertToDB(db *sql.DB, key string, url string) error {
	statement, err := db.Prepare("INSERT INTO links (key, url) VALUES (?, ?)")
	if err != nil {
		return err
	}

	fmt.Println("inserting " + url + "with key of: " + key)
	_, err = statement.Exec(key, url)
	if err != nil {
		return err
	}

	return nil
}

func FetchFromDB(db *sql.DB, requestedKey string) (string, error) {

	rows, err := db.Query("SELECT key, url FROM links")
	if err != nil {
		return "", err
	}

	var key string
	var url string
	for rows.Next() {
		rows.Scan(&key, &url)
		fmt.Println("Key: " + key + "\nUrl: " + url)
		if key == requestedKey {
			return url, nil
		}
	}

	return "err", errors.New("could not key in DB")
}
