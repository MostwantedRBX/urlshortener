package storage

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func StartDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./shortener.db")
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT PRIMARY KEY, url TEXT)")
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
	statement.Exec()
	log.Logger.Info().Caller().Msg("Database opened")
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
		if key == requestedKey {
			return url, nil
		}
	}

	return "", errors.New("could not key in DB")
}
