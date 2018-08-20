package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	mgo "gopkg.in/mgo.v2"
)

type bsonArticle struct {
	Headline string    `bson:"headline"`
	Content  string    `bson:"content"`
	Source   string    `bson:"source"`
	Time     time.Time `bson:"time"`
}

var cacheClient *redis.Client
var elasticClient *elastic.Client
var mongoClient *mgo.Session

func main() {
	fmt.Println("API-test-Server version 0.02")

	fmt.Println("Defining http endpoints...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)

	fmt.Println("Starting http server")
	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("eop")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the index endpoint")
	type response struct {
		Status  string
		Version int
	}

	resp := response{Status: "working", Version: 0}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
