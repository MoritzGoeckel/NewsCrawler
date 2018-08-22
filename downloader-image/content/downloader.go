package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
)

type article struct {
	Headline string    `json:"headline"`
	Content  string    `json:"content"`
	Source   string    `json:"source"`
	Url      string    `json:"url"`
	Time     time.Time `json:"time"`
}

type source struct {
	Urls []string
	Name string
	Id   string
}

func main() {
	fmt.Println("Downloader version 0.02")

	sources := readSources()
	agt, pq := getRedisConnections()

	for _, e := range sources {
		downloadSource(e, agt, pq)
	}

	fmt.Println("eop")
}

func downloadSource(s source, agt *redis.Client, pq *redis.Client) {
	for _, url := range s.Urls {
		doc := getHTML(url)

		//Create fallbacks for finding the content
		doc.Find(".headline").Each(func(i int, selection *goquery.Selection) {
			entry := article{Headline: selection.Text(), Source: s.Id, Content: "TODO", Time: time.Now(), Url: "TODO"}
			h := hashArticle(entry)
			pushed := false

			data, err := json.Marshal(entry)
			if err != nil {
				log.Fatal(err)
			}

			if !alreadyGotThat(h, agt) {
				setAlreadyGotThat(h, agt)
				pushNewEntry(string(data), pq)
				pushed = true
			}

			fmt.Println("New: " + fmt.Sprint(pushed) + "\t" + fmt.Sprint(h) + "\t" + entry.Headline)
		})
	}
}

func readSources() []source {
	content, err := ioutil.ReadFile("sources.json")
	if err != nil {
		log.Fatal(err)
	}

	sources := make([]source, 0)

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

func hashArticle(a article) uint32 {
	return hashString(a.Headline + a.Content + a.Source)
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

	if agtUrl == "" || pqUrl == "" {
		log.Fatal("Enviroment variables not set")
	}

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
