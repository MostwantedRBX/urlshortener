package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	//	Connection data is retrieved from enviroment variables
	pgHost         = os.Getenv("PG_HOST")
	pgPort, _      = strconv.Atoi(os.Getenv("PG_PORT"))
	pgPass         = os.Getenv("PG_PASS")
	pgDatabaseName = os.Getenv("PG_DATABASE_NAME")
)

func StartDB() *sql.DB {
	//	Open initial connection to database
	db, err := sql.Open("postgres", fmt.Sprintf("host= %s port= %d user= postgres password= %s dbname= %s sslmode= disable", pgHost, pgPort, pgPass, pgDatabaseName))
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("")
	}

	//	Creates the table to store urls with a string as the key
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT, url varchar(250), PRIMARY KEY (key))")

	//	TODO: Figure out a more elegant way to handle failures to connect to DB
	count := 0
	for err != nil {
		log.Logger.Warn().Err(err)
		db, _ = sql.Open("postgres", fmt.Sprintf("host= %s port= %d user= postgres password= %s dbname= %s sslmode= disable", pgHost, pgPort, pgPass, pgDatabaseName))
		statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS links (key TEXT, url varchar(250), PRIMARY KEY (key))")

		count++
		time.Sleep(5 * time.Second)
		if count > 5 {
			panic("Could not connect to db")
		}
	}

	f, err := statement.Exec()
	fmt.Println(f)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("")
	}

	log.Logger.Info().Caller().Msg("Database opened")
	return db
}

func InsertToDB(db *sql.DB, key string, url string) error {
	log.Logger.Info().Msg("Attempting to insert the url/key combo into the db")

	// 	Insert a url with the key generated in the genKey() function in main.go

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

	//	Pull in everything from the database and scan it for the key, then pull the url that the key indicates.
	rows, err := db.Query("select key, url from links where key=$1;", requestedKey)

	if err != nil {
		return "", err
	}

	var key string
	var url string
	for rows.Next() {

		rows.Scan(&key, &url)
		if key == requestedKey {
			//	If we found a match, we close the rows
			rows.Close()
			return url, nil
		}
	}

	//	If there are no matches, close the rows and return an error
	rows.Close()
	return "", errors.New("could not find key in DB")
}
