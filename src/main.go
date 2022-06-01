package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"

	"github.com/mostwantedrbx/urlshortener/storage"
)

//	Docker info {
//	Download the postgres image for the database and run it with: 'docker run -e POSTGRES_PASSWORD=dbpasswordhere postgres:latest'
//	Build command: 'docker build --tag urlshortener:latest .'
//	Run command: 'docker run -e PG_PASS=dbpasswordhere -e PG_HOST=dbIPhere -e PG_PORT=5432 -e PG_DATABASE_NAME=dbnamehere --name nameofcontainer -p 8080:8080 -d urlshortener:latest'
//	}

var (
	//	Declare strings of text/numbers to use for shortened link keys
	//	Also, start the database and keep a pointer to it to pass around
	ALPHANUM string  = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"
	DB       *sql.DB = storage.StartDB()
	PROD     bool    = os.Getenv("PROD") == "true"
)

type UrlStruct struct {
	Url string `json:"url"`
}

//	Fires when the page /{key} is requested
func fetchUrl(w http.ResponseWriter, req *http.Request) {

	//	Vars is from the variable {keys} in the url /links/{key}
	vars := mux.Vars(req)
	// fmt.Println(vars["key"])
	if vars["key"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url, err := storage.FetchKeyUrlFromDB(DB, vars["key"])
	// log.Logger.Info().Msg(url)
	if err != nil {

		log.Logger.Err(err).Msg("Could not get rows")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//	This refreshes and redirects with the new url from the DB
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

//	Fires when the page /put/ is requested
func putUrl(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	key, err := genKey()
	if err != nil {
		log.Logger.Err(err).Caller().Msg("Could not decode json")
		return
	}

	var jsonData UrlStruct
	err = json.NewDecoder(req.Body).Decode(&jsonData)

	if !(len(jsonData.Url) > 0) {
		return
	} else if err != nil {
		log.Logger.Err(err).Caller().Msg("Could not decode json")
		return
	}

	if !strings.Contains(jsonData.Url, "http") {
		jsonData.Url = "http://" + jsonData.Url
	}

	if err := storage.InsertUrlIntoDB(DB, key, jsonData.Url); err != nil {
		log.Logger.Err(err).Msg("Could not insert url into db")
		http.Error(w, err.Error(), 500)
		return
	}

	jsonData.Url = "srtlink.net/" + key

	if err := json.NewEncoder(w).Encode(&jsonData); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// log.Logger.Info().Msg("Url: " + req.URL.RawQuery + "\n	   Key: " + key)
}

//	Function to generate the key for the link
func genKey() (string, error) {
	var (
		length  int = 5
		byteKey     = make([]byte, length)
	)

	//	Create a random string of characters using random characters from the ALPHANUM charset declared at the top
	for i := range byteKey {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(ALPHANUM))))
		if err != nil {
			return "", err
		}
		byteKey[i] = ALPHANUM[randNum.Int64()]
	}
	//	Convert it from a byte to a string
	stringKey := string(byteKey)

	//	Check if the string is in the database
	//	If there isn't a url in the database then we have a winner
	url, _ := storage.FetchKeyUrlFromDB(DB, stringKey)
	if url == "nil" {
		return stringKey, nil
	}

	//	If the key generated is already in use, call recursively
	//	to try to get another (hopefully unique) key
	return genKey()
}

func main() {

	//	Log setup
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, file))
	log.Logger.Info().Msg("Logs started")

	//	Set up a router for our event handlers
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./web/"))
	//	Serve /web/index.htm when localhost:8080/ is requested
	r.Handle("/", fs)
	r.HandleFunc("/put/", putUrl).Methods("POST")   //, "OPTIONS"	when either /links or /links{key} gets requested, hand the data to a function
	r.HandleFunc("/{key}", fetchUrl).Methods("GET") //	{key} is a variable that gets handed to the function fetchUrl()

	r.PathPrefix("/").Handler(fs)

	//	Server settings
	m := &autocert.Manager{
		Cache:      autocert.DirCache("/root/secrets"),
		Prompt:     autocert.AcceptTOS,
		Email:      "mostwantedrbxsteam@gmail.com",
		HostPolicy: autocert.HostWhitelist("srtlink.net"),
	}

	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//	Listen @ localhost:80 for a request
	log.Logger.Info().Msg("Starting server @localhost" + server.Addr)
	if PROD {
		server.Addr = ":443"
		server.TLSConfig = m.TLSConfig()
		go func() {
			log.Logger.Fatal().Err(http.ListenAndServe(":http", m.HTTPHandler(nil))).Msg("Server failed to run")
		}()
		log.Logger.Fatal().Err(server.ListenAndServeTLS("", "")).Msg("Server failed to run")
	} else {
		log.Logger.Fatal().Err(server.ListenAndServe()).Msg("Server failed to run")
	}
}
