package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetArticle(document *goquery.Document) Article {
	document.Find(".headline").Each(func(i int, selection *goquery.Selection) {

	})

	//TODO ???

	return Article{}
}

func GetHTML(url string) *goquery.Document {
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
