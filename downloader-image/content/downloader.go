// go get github.com/PuerkitoBio/goquery
// go get -u github.com/go-redis/redis

package main

import (
	"fmt"
	"log"
	"net/http"
    "github.com/go-redis/redis"
	"github.com/PuerkitoBio/goquery"
)

// This example scrapes the reviews shown on the home page of metalsucks.net.
func main() {
    url := "agt-redis.default.svc.cluster.local"

    client := redis.NewClient(&redis.Options{
        Addr:     url + ":6379",
        Password: "",
        DB:       0,
    })

    err := client.Set("thekey", "thevalue", 0).Err()
    if err != nil {
        log.Fatal(err)
    }

    // Request the HTML page.
	res, err := http.Get("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})

    fmt.Println("eop")

	// To see the output of the Example while running the test suite (go test), simply
	// remove the leading "x" before Output on the next line. This will cause the
	// example to fail (all the "real" tests should pass).

	// xOutput: voluntarily fail the Example output.
}
