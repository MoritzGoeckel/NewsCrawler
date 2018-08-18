package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("Downloader version 0.01")

	agt, pq := getRedisConnections()

	doc := getHTML("http://spiegel.de")

	doc.Find(".headline").Each(func(i int, s *goquery.Selection) {
		headline := s.Text()
		h := hashString(headline)
		pushed := false

		if !alreadyGotThat(h, agt) {
			setAlreadyGotThat(h, agt)
			pushNewEntry(headline, pq)
			pushed = true
		}

		fmt.Println("New: " + fmt.Sprint(pushed) + "\t" + fmt.Sprint(h) + "\t" + headline)
	})

	fmt.Println("eop")
}

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func getHTML(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
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

func getRedisConnections() (*redis.Client, *redis.Client) {
	agtUrl := os.Getenv("agt-redis-url")
	pqUrl := os.Getenv("pq-redis-url")

	fmt.Println("agt url: " + agtUrl)
	fmt.Println("pq url: " + pqUrl)

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

	return agtClient, pqClient
}
