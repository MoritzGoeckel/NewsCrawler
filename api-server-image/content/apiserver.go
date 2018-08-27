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
	"github.com/vjeantet/jodaTime"
	mgo "gopkg.in/mgo.v2"
)

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
	router.HandleFunc("/number_articles/", numberArticlesEndpoint)
	router.HandleFunc("/get_words/", getWordsEndpoint)
	router.HandleFunc("/get_words_todate/", getWordsToDateEndpoint)
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
		Date    string
	}

	resp := response{Status: "working", Version: 0, Date: jodaTime.Format("YYYY.MM.dd", time.Now())}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func numberArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the numberArticles endpoint")

	type response struct {
		Documents int
	}

	resp := response{Documents: countMongo(mongoClient)}

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

	js, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getWordsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the words endpoint")
	result := getWords(mongoClient)

	js, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getWordsToDateEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the words to date endpoint")
	result := getWordsToDate(mongoClient)

	js, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getWords(client *mgo.Session) []Word {
	collection := client.DB("news").C("words")

	var all []Word
	err := collection.Find(nil).Sort("-count").Limit(30).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	return all
}

func getWordsToDate(client *mgo.Session) []WordToDate {
	collection := client.DB("news").C("words_to_date")

	var all []WordToDate
	err := collection.Find(nil).Sort("-date", "-count").Limit(30).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	return all
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

func searchDocument(client *elastic.Client, query string) []JsonArticle {
	fmt.Println("Searching in elastic for: " + query)
	ctx := context.Background()
	termQuery := elastic.NewTermQuery("headline", query)
	searchResult, err := client.Search().
		Index("articles").
		Query(termQuery).
		Sort("time", true).
		From(0).Size(10).
		Pretty(false).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Query took %d milliseconds and resulted in %d\n",
		searchResult.TookInMillis,
		searchResult.TotalHits())

	var articles []JsonArticle

	var ttyp JsonArticle
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(JsonArticle)
		fmt.Println("Found: " + t.Headline)
		articles = append(articles, t)
	}

	return articles
}

func countMongo(client *mgo.Session) int {
	fmt.Println("Counting mongodb")

	coll := client.DB("news").C("articles")

	result, err := coll.Count()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found " + fmt.Sprint(result) + " documents")

	return result
}

func searchMongo(client *mgo.Session) []BsonArticle {
	fmt.Println("Listing mongodb")

	coll := client.DB("news").C("articles")

	var all []BsonArticle
	err := coll.Find(nil).Sort("time").Limit(10).All(&all)
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

	if elasticUrl == "" || cacheUrl == "" || mongoUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

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
