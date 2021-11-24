package storage

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func StartDB() *sql.DB {

	//	open/create(if it isn't made) the database
	db, err := sql.Open("sqlite3", "./shortener.db")
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
	//	creates the table to store urls with a string as the key
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT PRIMARY KEY, url TEXT)")
	if err != nil {
		log.Logger.Fatal().Err(err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Logger.Fatal().Err(err)
	}

	log.Logger.Info().Caller().Msg("Database opened")
	return db
}

func InsertToDB(db *sql.DB, key string, url string) error {
	log.Logger.Info().Msg("Attempting to insert the url/key combo into the db")

	// 	insert a url with the key generated in the genKey() function in main.go
	statement, err := db.Prepare("INSERT INTO links (key, url) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(key, url)
	if err != nil {
		return err
	}

	return nil
}

func FetchFromDB(db *sql.DB, requestedKey string) (string, error) {

	//	pull in everything from the database and scan it for the key, then pull the url that the key indicates.
	rows, err := db.Query("SELECT key, url FROM links")

	if err != nil {
		return "", err
	}

	var key string
	var url string
	for rows.Next() {
		rows.Scan(&key, &url)
		if key == requestedKey {
			//	if we found a match, we close the rows
			rows.Close()
			return url, nil
		}
	}

	//	if there are no matches, close the rows and return an error
	rows.Close()
	return "", errors.New("could not find key in DB")
}
