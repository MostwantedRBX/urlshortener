package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	pgPass = os.Getenv("PG_PASS")
	pgPort = 5432 //	TODO: move to env var
)

func StartDB() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("host= localhost port= %d user= postgres password= %s dbname= urlshortener sslmode= disable", pgPort, pgPass))
	if err != nil {
		log.Logger.Fatal().Err(err)
	}
	//	creates the table to store urls with a string as the key
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT, url varchar(250), PRIMARY KEY (key))")
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("")
	}

	f, err := statement.Exec() //	Line 33
	fmt.Println(f)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("")
	}

	log.Logger.Info().Caller().Msg("Database opened")
	return db
}

func InsertToDB(db *sql.DB, key string, url string) error {
	log.Logger.Info().Msg("Attempting to insert the url/key combo into the db")

	// 	insert a url with the key generated in the genKey() function in main.go

	statement, err := db.Prepare(`INSERT INTO links(key, url) VALUES ($1, $2);`)
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
	rows, err := db.Query("SELECT key, url FROM links;")

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
