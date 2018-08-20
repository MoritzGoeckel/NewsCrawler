package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2"
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
				"time":{
					"type":"date"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func main() {
	fmt.Println("Processor version 0.02")

	ctx := context.Background()
	pq, elastic, mongo := getConnections()
	//defer mongo.Close()

	fmt.Println("Ensuring index for elastic")
	ensureIndex(&ctx, elastic)
	fmt.Println("Index done")

	for {
		message := getNextInQueue(pq)
		fmt.Println("Processing message: " + message)
		insertIntoMongo(bsonArticle{Headline: message, Content: "NOTSET", Source: "NOTSET", Time: time.Now()}, mongo)
		insertIntoElastic(jsonArticle{Headline: message, Content: "NOTSET", Source: "NOTSET", Time: time.Now()}, &ctx, elastic)
	}
}

func getNextInQueue(client *redis.Client) string {
	for {
		fmt.Println("Trying to retrieve message")
		val, err := client.BLPop(30*time.Second, "pending").Result()
		if err == redis.Nil {
			fmt.Println("No message in queue")
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

func insertIntoMongo(data bsonArticle, mongo *mgo.Session) {
	collection := mongo.DB("news").C("articles")

	err := collection.Insert(data)
	if err != nil {
		log.Fatal(err)
	}
}

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

// https://olivere.github.io/elastic/
func insertIntoElastic(article jsonArticle, ctx *context.Context, client *elastic.Client) {
	put1, err := client.Index().
		Index("articles").
		Type("article").
		Id("1"). //How to assign an id automatically?
		BodyJson(article).
		Do(*ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Indexed Article %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func getConnections() (*redis.Client, *elastic.Client, *mgo.Session) {
	pqUrl := os.Getenv("pq-redis-url")
	elasticUrl := os.Getenv("elastic-url")
	mongoUrl := os.Getenv("mongo-url")

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
