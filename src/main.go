package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mostwantedrbx/urlshortener/storage"
)

//	declare strings of text/numbers to use for shortened link keys
//	also, start the database and keep a pointer to it to pass around
var (
	ALPHANUM string  = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	DB       *sql.DB = storage.StartDB()
)

type UrlStruct struct {
	Url string `json:"url"`
}

//	fires when the page /{key} is requested
func fetchUrl(w http.ResponseWriter, req *http.Request) {

	//	vars is from the variable {keys} in the url /links/{key}
	vars := mux.Vars(req)

	url, err := storage.FetchFromDB(DB, vars["key"])
	log.Logger.Info().Msg(url)
	if err != nil {

		log.Logger.Err(err).Msg("Could not get rows")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//	this makes it so the page refreshes with the new url from the DB
	fmt.Fprintln(w, `<head><meta http-equiv="refresh" content="0; url='`+url+`'" /></head>`)
}

//	fires when the page /links is requested
func putUrl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	key := genKey()

	var jsonData UrlStruct
	err := json.NewDecoder(req.Body).Decode(&jsonData)

	if !(len(jsonData.Url) > 0) {
		return
	} else if err != nil {
		log.Logger.Err(err).Caller().Msg("Could not decode json")
		return
	}

	log.Logger.Info().Msg("data: " + jsonData.Url)

	if !strings.Contains(jsonData.Url, "http") {
		jsonData.Url = "http://" + jsonData.Url
	}

	if err := storage.InsertToDB(DB, key, jsonData.Url); err != nil {
		log.Logger.Err(err).Msg("Could not insert url into db")
		http.Error(w, err.Error(), 500)
		return
	}

	jsonData.Url = "http://localhost:8080/" + key

	if err := json.NewEncoder(w).Encode(&jsonData); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// key := genKey()
	// u, err := url.Parse(req.URL.String())
	// if err != nil {
	// 	log.Logger.Err(err).Msg("Could not parse URL")
	// 	return
	// }
	// rawUrl := u.Query().Get("url")

	// if !strings.Contains(rawUrl, "http") {
	// 	log.Logger.Debug().Msg("Adding http to " + rawUrl)
	// 	rawUrl = "http://" + rawUrl
	// }
	// err = storage.InsertToDB(DB, key, rawUrl)
	// if err != nil {
	// 	log.Logger.Err(err).Msg("could not put url to key " + key)
	// 	fmt.Fprintln(w, "Failed, please try again")
	// 	return
	// }

	//	Fprint() writes to the page
	fmt.Fprintln(w, "<a href='http://localhost:8080/"+key+"'>http://localhost:8080/"+key+"</a>")
	log.Logger.Info().Msg("Url: " + req.URL.RawQuery + "\n	   Key: " + key)
}

//	function to generate the key for link
func genKey() string {

	var (
		length  int        = 5
		randGen *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
		byteKey            = make([]byte, length)
	)

	//	create a random string of characters using random characters from the ALPHANUM charset declared at the top
	for i := range byteKey {
		byteKey[i] = ALPHANUM[randGen.Intn(len(ALPHANUM))]
	}
	//	convert it from a byte to a string
	stringKey := string(byteKey)

	//	check if the string is in the database
	//	if there isn't a url in the database then we have a winner
	url, _ := storage.FetchFromDB(DB, stringKey)
	if url == "" {
		return stringKey
	}

	//	if we did not find a url linked to a key, call
	//	recursively to try to get another (hopefully unique) key
	return genKey()
}

func main() {

	//	log setup
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, file))
	log.Logger.Info().Msg("Logs started")

	//	set up a router for our event handlers
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./web/"))
	//	serve /static/index.htm when localhost:8080/ is requested
	r.Handle("/", fs)
	r.HandleFunc("/links/put/", putUrl).Methods("POST", "OPTIONS") //	when either /links or /links{key} gets requested, hand the data to a function
	r.HandleFunc("/{key}", fetchUrl).Methods("GET")                //	{key} is a variable that gets handed to the function fetchUrl()

	r.PathPrefix("/").Handler(fs)

	//	server settings
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Logger.Info().Msg("Starting server @http://localhost" + server.Addr)
	//	listen @ localhost:8080 for a request
	log.Logger.Fatal().Err(server.ListenAndServe()).Msg("Server failed to run")
}
