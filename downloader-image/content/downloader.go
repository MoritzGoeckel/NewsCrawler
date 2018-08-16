// go get github.com/PuerkitoBio/goquery
// go get -u github.com/go-redis/redis

package main

import (
	"strconv"
    "os"
    "fmt"
	"log"
	"net/http"
    "github.com/go-redis/redis"
	"github.com/PuerkitoBio/goquery"
)

func main() {
    //agt, pq := getRedisConnections()

    doc := getHTML("http://spiegel.de")

	doc.Find(".headline").Each(func(i int, s *goquery.Selection) {
		//band := s.Find("a").Text()
		//title := s.Find("i").Text()
		//fmt.Printf("Review %d: %s - %s\n", i, band, title)
        fmt.Println(s.Text())
	})

    fmt.Println("eop")
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

func pushNewEntry(data string, client *redis.Client){
    err := client.LPush("pending", data).Err()
    if err != nil {
        log.Fatal(err)
    }
}

func getAlreadyGotThat(hash int, client *redis.Client) bool{
    _, err := client.Get(strconv.Itoa(hash)).Result()
    if err == redis.Nil {
        return false
    } else if err != nil {
        log.Fatal(err)
        return true
    } else {
        return true
        //fmt.Println("key2", val)
    }
}

func setAlreadyGotThat(hash int, client *redis.Client) {
    err := client.Set(strconv.Itoa(hash), "seen",  60 * 60 * 72).Err()
    if err != nil {
        log.Fatal(err)
    }
}

func getRedisConnections() (*redis.Client, *redis.Client){
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
