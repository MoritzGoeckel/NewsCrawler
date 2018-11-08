package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	router.HandleFunc("/info", infoEndpoint)
	router.HandleFunc("/mongo_articles", mongoArticlesEndpoint)
	router.HandleFunc("/elastic_articles/{query}", elasticArticlesEndpoint)
	router.HandleFunc("/number_articles/", numberArticlesEndpoint)
	router.HandleFunc("/number_words/", numberWordsEndpoint)
	router.HandleFunc("/number_words_todate/", numberWordsTodateEndpoint)
	router.HandleFunc("/get_words/", getWordsEndpoint)
	router.HandleFunc("/get_words_todate/", getWordsToDateEndpoint)
	router.HandleFunc("/get_word_cloud/", getWordCloudEndpoint)
	router.HandleFunc("/word_statistics/{query}", getWordStatisticsEndpoint)
	router.HandleFunc("/get_article_count_since/{date}", getArticelCountSinceEndpoint)
	router.HandleFunc("/get_high_entropy_article_since/{date}", getHighEntropyArticleSinceEndpoint)
	router.HandleFunc("/get_headline_ngrams/", getHeadlineNgramsEndpoint)

	//In the end
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	fmt.Println("Starting http server")
	log.Fatal(http.ListenAndServe(":80", router))

	//TODO: Use cache (Create caching func)

	fmt.Println("eop")
}

func infoEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the index endpoint")
	type response struct {
		Status  string
		Version int
		Date    string
		Epoch   int64
	}

	resp := response{Status: "working!", Version: 0, Date: getToday(), Epoch: getEpochNow()}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func getWordStatisticsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the getWordStatistics endpoint")
	vars := mux.Vars(r)
	query := vars["query"]

	collection := mongoClient.DB("news").C("words_to_date")

	var all []WordToDate
	err := collection.Find(bson.M{"word": query}).Sort("date").Limit(100).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(all)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func numberWordsEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the numberWords endpoint")

	type response struct {
		Words int
	}

	resp := response{Words: countMongoWords(mongoClient)}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func getWordCloudEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the numberWords endpoint")

	collection := mongoClient.DB("news").C("word_cloud")

	var all []ScoredWord
	err := collection.Find(nil).Sort("-score").Limit(30).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	js, err := json.Marshal(all)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func numberWordsTodateEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request on the numberWords todate endpoint")

	type response struct {
		WordsTodate int
	}

	resp := response{WordsTodate: countMongoWordsTodate(mongoClient)}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func getWords(client *mgo.Session) []Word {
	collection := client.DB("news").C("words")

	var all []Word
	err := collection.Find(bson.M{"date": getToday()}).Sort("-count").Limit(30).All(&all)
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

func searchDocument(client *elastic.Client, query string) []JsonArticle {
	fmt.Println("Searching in elastic for: " + query)
	ctx := context.Background()
	termQuery := elastic.NewTermQuery("headline", query)
	searchResult, err := client.Search().
		Index("articles").
		Query(termQuery).
		Sort("datetime", true).
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

func countMongoWords(client *mgo.Session) int {
	fmt.Println("Counting mongodb")

	coll := client.DB("news").C("words")

	result, err := coll.Count()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found " + fmt.Sprint(result) + " documents")

	return result
}

func countMongoWordsTodate(client *mgo.Session) int {
	fmt.Println("Counting mongodb")

	coll := client.DB("news").C("words_to_date")

	//Date is today in YY.MM.dd
	//This is the word count table, not to be confused with the article table
	result, err := coll.Find(bson.M{"date": getToday()}).Count()
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
	err := coll.Find(nil).Sort("datetime").Limit(10).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found " + fmt.Sprint(len(all)) + " documents")

	return all
}

func getArticelCountSinceEndpoint(w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.DB("news").C("articles")

	vars := mux.Vars(r)
	dateq := vars["date"]
	fmt.Println("Datequeuery ", dateq)

	count, err := collection.Find(bson.M{"datetime": bson.M{"$gt": parseRFCTimestringToEpoch(dateq)}}).Count()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Count ", count)

	js, err2 := json.Marshal(count)
	if err2 != nil {
		log.Fatal(err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func getHighEntropyArticleSinceEndpoint(w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.DB("news").C("articles")

	vars := mux.Vars(r)
	dateq := vars["date"]

	var a []BsonArticle
	query := bson.M{
		"$and": []bson.M{
			bson.M{"datetime": bson.M{"$gt": parseRFCTimestringToEpoch(dateq)}},
			bson.M{"article_perplexity": bson.M{"$exists": true}},
		},
	}
	err := collection.Find(query).Sort("-article_perplexity").All(&a)
	//err := collection.Find(query).Sort("-article_perplexity").Limit(10).All(&a)
	//TODO: what should be the limit? Should the client decide?
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Highest Entropy-Article", a)

	js, err2 := json.Marshal(a)
	if err2 != nil {
		log.Fatal(err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

func getHeadlineNgramsEndpoint(w http.ResponseWriter, r *http.Request) {
	collection := mongoClient.DB("news").C("headlines")

	var h BsonNgram2Words
	err := collection.Find(bson.M{"ngram2words": bson.M{"$exists": true}}).One(&h)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Headline ngrams fetched...", h)

	js, err2 := json.Marshal(h)
	if err2 != nil {
		log.Fatal(err2)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

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
