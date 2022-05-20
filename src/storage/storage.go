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

//	StartDB opens the sql connection and creates a table named 'links' and two collumns
//	in that table named 'key' and 'url' with 'key' being the primary key to the table.
//	Returns an sql.DB object.
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

//	InsertUrlIntoDB takes in a db, a key and a url to link to the key and inserts it into
//	the database table named 'links' and into their respective collumns called 'key' and 'url'.
func InsertUrlIntoDB(db *sql.DB, key string, url string) error {
	log.Logger.Info().Msg("Attempting to insert the url/key combo into the db")

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

//	FetchFromDB takes an sql.DB and a string of the key you want the url for. Returns the url of the requested key.
//	Returns 'nil' in a string alongside the error if it could not find the key in the database or otherwise errors.
func FetchKeyUrlFromDB(db *sql.DB, requestedKey string) (string, error) {
	rows, err := db.Query("select key, url from links where key=$1;", requestedKey)

	if err != nil {
		return "err", err
	}

	var key string
	var url string
	//	I can probably remove this loop now that I'm selecting a specific key from the table.
	for rows.Next() {
		err = rows.Scan(&key, &url)
		if err != nil {
			return "err", err
		}
		if key == requestedKey {
			//	If we found a match, we close the rows
			rows.Close()
			return url, nil
		}
	}

	//	If there are no matches, close the rows and return an error
	rows.Close()
	return "nil", errors.New("could not find key in DB")
}
