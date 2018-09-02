package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/vjeantet/jodaTime"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	fmt.Println("Word cloud generator version 0.02")

	mongo := getConnection()

	fmt.Println("Retrieving words")
	words := getWords(mongo)

	fmt.Println("Retrieving words to date")
	todayWords := getWordsToDate(mongo)

	fmt.Println("Making words searchable")
	wordsBaselineMap := make(map[string]int)

	for _, word := range words {
		wordsBaselineMap[word.Word] = word.Count
		if word.Count == 0 {
			log.Fatal("Assertion: Word count has been 0!")
		}
	}

	//Number of spots on the leaderboard
	num := 100
	var leaderboard []ScoredWord

	addToLeaderboard := func(word string, score float64) {
		leaderboard = append(leaderboard, ScoredWord{word, score})
		sort.Slice(leaderboard, func(i, j int) bool {
			return leaderboard[i].Score > leaderboard[j].Score
		})

		if len(leaderboard) > num {
			leaderboard = leaderboard[:num] // + 1?
		}
	}

	fmt.Println("Calculating scores")
	for step, word := range todayWords {
		score := float64(word.Count) / float64(wordsBaselineMap[word.Word])

		if len(leaderboard) < num || leaderboard[len(leaderboard)-1].Score < score {
			if len(leaderboard) < num {
				fmt.Println("Filling empty slots")
			}
			addToLeaderboard(word.Word, score)
		}

		if step%(len(todayWords)/10) == 0 {
			fmt.Printf("%f percent done\n", float64(step)/float64(len(todayWords))*100.0)
		}
	}

	fmt.Println("The leaderboard")
	fmt.Print(leaderboard)
	fmt.Print("\n")

	collection := mongo.DB("news").C("word_cloud")

	fmt.Println("Deleting old from mongo")
	_, err := collection.RemoveAll(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Writing to mongo")
	bulk := collection.Bulk()
	for _, entry := range leaderboard {
		bulk.Insert(entry)
	}
	_, err = bulk.Run()
	if err != nil {
		log.Fatal(err)
	}

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

	var all []WordToDate
	err := collection.Find(bson.M{"date": now, "count": bson.M{"$gt": 5}}).Sort("-count").Limit(50 * 1000).All(&all)
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
