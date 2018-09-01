package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("Article downloader version 0.02")

	agt, pq, lq := getRedisConnections()

	for {
		message := getNextInQueue(lq)

		var link Link
		err := json.Unmarshal([]byte(message), &link)
		if err != nil {
			log.Fatal(err)
		}

		downloadArticle(link, agt, pq)
	}
}

func downloadArticle(link Link, agt *redis.Client, pq *redis.Client) {
	doc, err := GetHTML(link.Url)
	if err != nil {
		fmt.Print("Warning: " + link.Url + " -> ")
		fmt.Print(err)
		fmt.Print("\r\n")
		return
	}

	article, err := GetArticle(doc, link.Url, link.Source)
	if err != nil {
		fmt.Print("Warning: " + link.Url + " -> ")
		fmt.Print(err)
		fmt.Print("\r\n")
		return
	}

	h := hashArticle(article)
	pushed := false

	data, err := json.Marshal(article)
	if err != nil {
		log.Fatal(err)
	}

	if !alreadyGotThat(h, agt) {
		setAlreadyGotThat(h, agt)
		pushNewEntry(string(data), pq)
		pushed = true
	}

	pushedStr := "agt"
	if pushed {
		pushedStr = "new"
	}

	fmt.Println(pushedStr + "\t" + article.Headline)
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

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func hashArticle(a Article) uint32 {
	return hashString(a.Headline + a.Content + a.Source)
}

func pushNewEntry(data string, client *redis.Client) {
	err := client.LPush("pending", data).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func alreadyGotThat(hash uint32, client *redis.Client) bool {
	_, err := client.Get(fmt.Sprint(hash)).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Fatal(err)
		return true
	} else {
		return true
	}
}

func setAlreadyGotThat(hash uint32, client *redis.Client) {
	expiration := time.Duration(72) * time.Hour
	err := client.Set(fmt.Sprint(hash), "seen", expiration).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func getRedisConnections() (*redis.Client, *redis.Client, *redis.Client) {
	agtUrl := os.Getenv("agt-article-redis-url")
	pqUrl := os.Getenv("pq-redis-url")
	lqUrl := os.Getenv("lq-redis-url")

	if agtUrl == "" || pqUrl == "" || lqUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

	fmt.Println("agt url: " + agtUrl)
	fmt.Println("pq url: " + pqUrl)
	fmt.Println("lq url: " + lqUrl)

	agtClient := redis.NewClient(&redis.Options{
		Addr:     agtUrl + ":6379",
		Password: "",
		DB:       0,
	})

	pqClient := redis.NewClient(&redis.Options{
		Addr:     pqUrl + ":6379",
		Password: "",
		DB:       0,
	})

	lqClient := redis.NewClient(&redis.Options{
		Addr:     lqUrl + ":6379",
		Password: "",
		DB:       0,
	})

	return agtClient, pqClient, lqClient
}
