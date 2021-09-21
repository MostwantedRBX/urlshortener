package main

import (
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mostwantedrbx/urlshortener/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//	declare strings of text/numbers to use for shortened link keys
var (
	ALPHA string  = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"
	NUM   string  = "0123456789"
	DB    *sql.DB = storage.StartDB()
)

//	fires when the page /fetchurl is requested
func fetchUrl(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	url, err := storage.FetchFromDB(DB, vars["key"])
	if err != nil {
		log.Logger.Err(err).Msg("could not fetch url from key provided")
	}
	fmt.Println(url)
	fmt.Fprintln(w, `<head><meta http-equiv="refresh" content="0; url='`+url+`'" /></head>`) //(w, "<a href='"+url+"'>"+url+"</a>")
}

//	fires when the page /puturl is requested
func putUrl(w http.ResponseWriter, req *http.Request) {
	//	just going to assume this key isnt the same as any other in the DB right now.
	key := genKey()
	err := storage.InsertToDB(DB, key, req.URL.RawQuery)
	if err != nil {
		log.Logger.Err(err).Msg("could not put url to key " + key)
		fmt.Fprintln(w, "Failed, please try again")
		return
	}
	// Fprint() writes to the page
	fmt.Fprintln(w, "<a href='http://localhost:8090/links/"+key+"'>http://localhost:8090/links/"+key+"</a>")
	log.Logger.Info().Msg("Url: " + req.URL.RawQuery + "\n	Key: " + key)
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
	url, _ := storage.FetchFromDB(DB, key)
	if url == "" {
		//	if there isnt a url then we have a winner
		return key
	}
	//	call recusively to try to get another unique key
	return genKey()
}

func main() {
	rand.Seed(time.Now().Unix())
	//	log setup
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, file))
	log.Logger.Info().Msg("Logs started")

	r := mux.NewRouter()
	r.HandleFunc("/links", putUrl).Methods("POST", "GET")
	r.HandleFunc("/links/{key}", fetchUrl).Methods("GET")

	server := &http.Server{
		Handler:      r,
		Addr:         ":8090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Logger.Info().Msg("Starting server @localhost" + server.Addr)

	log.Logger.Fatal().Err(server.ListenAndServe()).Msg("Server failed to run")
}
