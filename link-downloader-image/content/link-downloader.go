package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("Link downloader version 0.02")

	sources := readSources()
	agt, lq := getRedisConnections()

	for _, e := range sources {
		downloadSource(e, agt, lq)
	}

	fmt.Println("eop")
}

func downloadSource(s Source, agt *redis.Client, lq *redis.Client) {
	for _, url := range s.Urls {
		links, err := GetLinks(url)
		if err != nil {
			fmt.Print("Warning: " + url + " -> ")
			fmt.Print(err)
			fmt.Print("\r\n")
		} else {
			for _, link := range links {
				link.Source = s.Id

				h := hashLink(link)
				pushed := false

				data, err := json.Marshal(link)
				if err != nil {
					log.Fatal(err)
				}

				if !alreadyGotThat(h, agt) {
					setAlreadyGotThat(h, agt)
					pushNewEntry(string(data), lq)
					pushed = true
				}

				fmt.Println("New: " + fmt.Sprint(pushed) + "\t" + fmt.Sprint(h) + "\t" + link.Url)
			}
		}
	}
}

func readSources() []Source {
	content, err := ioutil.ReadFile("sources.json")
	if err != nil {
		log.Fatal(err)
	}

	sources := make([]Source, 0)

	err = json.Unmarshal(content, &sources)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded %d sources\r\n", len(sources))

	return sources
}

func hashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func hashLink(a Link) uint32 {
	return hashString(a.Url + a.Source)
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
	agtUrl := os.Getenv("agt-link-redis-url")
	lqUrl := os.Getenv("lq-redis-url")

	if agtUrl == "" || lqUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

	fmt.Println("agt link url: " + agtUrl)
	fmt.Println("lq url: " + lqUrl)

	agtClient := redis.NewClient(&redis.Options{
		Addr:     agtUrl + ":6379",
		Password: "",
		DB:       0,
	})

	lqClient := redis.NewClient(&redis.Options{
		Addr:     lqUrl + ":6379",
		Password: "",
		DB:       0,
	})

	return agtClient, lqClient
}
