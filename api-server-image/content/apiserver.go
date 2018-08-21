package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
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

type jsonArticle struct {
	Headline string                `json:"headline"`
	Content  string                `json:"content"`
	Source   string                `json:"source"`
	Time     time.Time             `json:"time"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

var cacheClient *redis.Client
var elasticClient *elastic.Client
var mongoClient *mgo.Session

func main() {
	fmt.Println("API-Server version 0.02")

	fmt.Println("Connecting to services...")
	cacheClient, elasticClient, mongoClient = getConnections()

	fmt.Println("Defining http endpoints...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/mongo_articles", mongoArticlesEndpoint)
	router.HandleFunc("/elastic_articles/{query}", elasticArticlesEndpoint)
    //Set Endpoint without query to get everything

	fmt.Println("Starting http server")
	log.Fatal(http.ListenAndServe(":80", router))

	//TODO: Use cache (Create caching func)

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

func mongoArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the mongo endpoint")
	result := searchMongo(mongoClient)

	js, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func elasticArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the elastic endpoint")
	vars := mux.Vars(r)
	query := vars["query"]
	result := searchDocument(elasticClient, query)

	//TODO: Does not work, returns null

	js, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*func getCache(client *redis.Client) {
	val, err := client.Get(fmt.Sprint(hash)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Fatal(err)
		return true
	} else {
		return true
	}
}

func setCache(client *redis.Client) {
	expiration := time.Duration(1) * time.Hour
	err := client.Set(fmt.Sprint(hash), "seen", expiration).Err()
	if err != nil {
		log.Fatal(err)
	}
}*/

func searchDocument(client *elastic.Client, query string) []jsonArticle {
	fmt.Println("Searching in elastic for: " + query)
	ctx := context.Background()
	termQuery := elastic.NewTermQuery("headline", query)
	searchResult, err := client.Search().
		Index("articles").
		Query(termQuery).
		//Sort("time", true).
		//From(0).Size(10). // take documents 0-9
		//Pretty(false).    // pretty print request and response JSON
		Do(ctx) // execute
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query took %d milliseconds and resulted in %d\n",
		searchResult.TookInMillis,
		searchResult.TotalHits())

	var articles []jsonArticle

	var ttyp jsonArticle
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(jsonArticle)
		fmt.Println("Found: " + t.Headline)
		articles = append(articles, t)
	}

	return articles
}

func searchMongo(client *mgo.Session) []bsonArticle {
	fmt.Println("Listing mongodb")

	coll := client.DB("news").C("articles")

	var all []bsonArticle
	err := coll.Find(nil).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found " + fmt.Sprint(len(all)) + " documents")

	return all
}

func getConnections() (*redis.Client, *elastic.Client, *mgo.Session) {
	cacheUrl := os.Getenv("cache-redis-url")
	elasticUrl := os.Getenv("elastic-url")
	mongoUrl := os.Getenv("mongo-url")

	fmt.Println("cache url: " + cacheUrl)
	fmt.Println("elastic url: " + elasticUrl)
	fmt.Println("mongo url: " + mongoUrl)

	mongoPw := os.Getenv("mongo-pw")
	mongoUser := os.Getenv("mongo-user")

	fmt.Println("mongo credentials: " + mongoUser + " " + mongoPw)

	fmt.Println("Connecting to redis...")

	cacheClient := redis.NewClient(&redis.Options{
		Addr:     cacheUrl + ":6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("Redis connection established")
	fmt.Println("Connecting to elastic...")

	elasticClient, err := elastic.NewClient(
		elastic.SetURL("http://"+elasticUrl+":9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Elastic connection established")
	fmt.Println("Connecting to mongo...")

	mongoClient, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{mongoUrl + ":27017"},
		Username: mongoUser,
		Password: mongoPw,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo connection established")

	return cacheClient, elasticClient, mongoClient
}
