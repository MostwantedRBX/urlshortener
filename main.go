package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/mostwantedrbx/urlshortener/storage"
)

//	declare strings of text/numbers to use for shortened link keys
var (
	ALPHA string  = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	NUM   string  = "0123456789"
	DB    *sql.DB = storage.StartDB() // dunno if this is fine
)

//	fires when the page /fetchurl is requested
func fetchUrl(w http.ResponseWriter, req *http.Request) {
	url, err := storage.FetchFromDB(DB, req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, url)
}

//	fires when the page /puturl is requested
func putUrl(w http.ResponseWriter, req *http.Request) {
	// just going to assume this key isnt the same as any other in the DB right now.
	key := genKey()
	err := storage.InsertToDB(DB, key, req.URL.RawQuery)
	if err != nil {
		fmt.Println("failed to put the url in DB")
		fmt.Fprintln(w, "Failed, please try again")
		return
	}
	// Fprint() writes to the page
	fmt.Fprintln(w, "http://localhost:8090/fetchurl?"+key)
}

//	function to generate the key for link
func genKey() string {
	//	length of key
	var (
		LENGTH int    = 5
		key    string = ""
	)

	//	I'd like to come up with a better randomizer
	for i := 0; i < LENGTH; i++ {
		if rand.Intn(2) == 1 {
			key += string(ALPHA[rand.Intn(len(ALPHA))])
			continue
		}
		key += string(NUM[rand.Intn(len(NUM))])
	}
	return key
}

func main() {
	rand.Seed(time.Now().Unix())
	genKey()

	//	what functions handle what requests
	http.HandleFunc("/fetchurl", fetchUrl)
	http.HandleFunc("/puturl", putUrl)
	//	listen for reqeusts on port 8090. I'm unsure if I need to make my own handler type yet
	http.ListenAndServe(":8090", nil)
}
