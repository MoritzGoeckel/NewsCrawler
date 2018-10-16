package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"github.com/vjeantet/jodaTime"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var replaceArray []string

func main() {
	fmt.Println("Processor version 0.02")

	ctx := context.Background()
	pq, elastic, mongo := getConnections()
	defer mongo.Close()

	fmt.Println("Ensuring index for elastic")
	ensureIndex(&ctx, elastic)

	for {
		message := getNextInQueue(pq)

		var a Article
		err := json.Unmarshal([]byte(message), &a)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Processing message: " + message)

		fmt.Println("Inserting into Mongo")
		insertIntoMongo(BsonArticle{Headline: a.Headline, Description: a.Description, Image: a.Image, Content: a.Content, Source: a.Source, DateTime: a.DateTime, Language: a.Language, Url: a.Url}, mongo)

		fmt.Println("Inserting into elastic")
		insertIntoElastic(JsonArticle{Headline: a.Headline, Description: a.Description, Image: a.Image, Content: a.Content, Source: a.Source, DateTime: a.DateTime, Language: a.Language, Url: a.Url}, &ctx, elastic)

		fmt.Println("Gettings words")
		words := getWords(a)
		fmt.Print(words)

		fmt.Println("Updating word count")
		updateWordCount(words, mongo)

		fmt.Println("Updating word count for to date")
		updateWordCountToDate(words, mongo)
	}
}

func getNextInQueue(client *redis.Client) string {
	for {
		//fmt.Println("Trying to retrieve message")
		val, err := client.BLPop(60*time.Second, "pending").Result()
		if err == redis.Nil {
			//fmt.Println("No message in queue")
			continue
		} else if err != nil {
			log.Fatal(err)
			time.Sleep(10 * time.Second)
			continue
		} else {
			return val[1]
		}
	}
}

func getWords(a Article) []string {
	content := a.Headline + " " + a.Description + " " + a.Content

	seperators := "!\"²³$%&/{()[]}=\\?´`+*~#'.:,;<>|^"
	for _, sign := range seperators {
		content = strings.Replace(content, string(sign), " ", -1)
	}

	words := strings.Fields(content)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	return words
}

func updateWordCount(words []string, mongo *mgo.Session) {
	collection := mongo.DB("news").C("words")
	bulk := collection.Bulk()

	for _, word := range words {
		count := 1

		query := bson.M{"word": word}
		change := bson.M{"$inc": bson.M{"count": count}, "$set": bson.M{"word": word}}

		bulk.Upsert(query, change)
	}

	res, err := bulk.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Result: ")
	fmt.Print(res)
	fmt.Print("\r\n")
}

func updateWordCountToDate(words []string, mongo *mgo.Session) {
	now := jodaTime.Format("YYYY.MM.dd", time.Now())

	collection := mongo.DB("news").C("words_to_date")
	bulk := collection.Bulk()

	for _, word := range words {
		count := 1

		query := bson.M{"word": word, "date": now}
		change := bson.M{"$inc": bson.M{"count": count}, "$set": bson.M{"word": word, "date": now}}

		bulk.Upsert(query, change)
	}

	res, err := bulk.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Result: ")
	fmt.Print(res)
	fmt.Print("\r\n")

	//Todo: Create some kind of cleanup for very seldom words every 24h
	//Todo: Detect the language
}

func insertIntoMongo(data BsonArticle, mongo *mgo.Session) {
	collection := mongo.DB("news").C("articles")

	err := collection.Insert(data)
	if err != nil {
		log.Fatal(err)
	}
}

func insertIntoElastic(article JsonArticle, ctx *context.Context, client *elastic.Client) {
	put1, err := client.Index().
		Index("articles").
		Type("article").
		BodyJson(article).
		Do(*ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexed Article %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

const elasticMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"article":{
			"properties":{
				"headline":{
					"type":"text",
					"store": true,
					"fielddata": true
                },
                "content":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"description":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"url":{
					"type":"text",
					"store": true
				},
				"datetime":{
					"type":"date"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func ensureIndex(ctx *context.Context, client *elastic.Client) {
	exists, err := client.IndexExists("twitter").Do(*ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		createIndex, err := client.CreateIndex("twitter").BodyString(elasticMapping).Do(*ctx)
		if err != nil {
			log.Fatal(err)
		}
		if !createIndex.Acknowledged {
			log.Fatal("Surprise: Index not acknowledged!")
		}
	}
}

func getConnections() (*redis.Client, *elastic.Client, *mgo.Session) {
	pqUrl := os.Getenv("pq-redis-url")
	elasticUrl := os.Getenv("elastic-url")
	mongoUrl := os.Getenv("mongo-url")

	if elasticUrl == "" || pqUrl == "" || mongoUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

	fmt.Println("pq url: " + pqUrl)
	fmt.Println("elastic url: " + elasticUrl)
	fmt.Println("mongo url: " + mongoUrl)

	mongoPw := os.Getenv("mongo-pw")
	mongoUser := os.Getenv("mongo-user")

	fmt.Println("mongo credentials: " + mongoUser + " " + mongoPw)

	fmt.Println("Connecting to redis...")

	pqClient := redis.NewClient(&redis.Options{
		Addr:     pqUrl + ":6379",
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

	return pqClient, elasticClient, mongoClient
}
