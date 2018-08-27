package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vjeantet/jodaTime"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	fmt.Println("Word cloud generator version 0.01")

	mongo := getConnection()
	words := getWords(mongo)
	todayWords := getWordsToDate(mongo)

	//Walk though todayWords
	//Devide by words

	//Sort
	//Return or save

	fmt.Println("eop")
}

func getWords(client *mgo.Session) []Word {
	collection := client.DB("news").C("words")

	var all []Word
	err := collection.Find(nil).Sort("-count").Limit(50 * 1000).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	return all
}

func getWordsToDate(client *mgo.Session) []WordToDate {
	now := jodaTime.Format("YYYY.MM.dd", time.Now())
	collection := client.DB("news").C("words_to_date")

	var all []WordToDate //Min is 10
	err := collection.Find(bson.M{"date": now, "count": bson.M{"$gt": 10}}).Sort("-count").Limit(50 * 1000).All(&all)
	if err != nil {
		log.Fatal(err)
	}

	return all
}

func getConnection() *mgo.Session {
	mongoUrl := os.Getenv("mongo-url")

	if mongoUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

	fmt.Println("mongo url: " + mongoUrl)

	mongoPw := os.Getenv("mongo-pw")
	mongoUser := os.Getenv("mongo-user")

	fmt.Println("mongo credentials: " + mongoUser + " " + mongoPw)
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

	return mongoClient
}
