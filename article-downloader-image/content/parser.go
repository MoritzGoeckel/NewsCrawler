package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dyatlov/go-opengraph/opengraph"
)

func GetArticle(document *goquery.Document) (Article, bool) {
	document.Find(".headline").Each(func(i int, selection *goquery.Selection) {

	})

	html, err := document.Html()
	if err != nil {
		log.Fatal(err)
	}

	og := opengraph.NewOpenGraph()

	err = og.ProcessHTML(strings.NewReader(html))
	if err != nil {
		log.Fatal(err)
	}

	article := Article{Time: time.Now()}

	datapoints := 0
	if og.Title != "" {
		article.Headline = og.Title
		datapoints++
	}

	if og.URL != "" {
		article.Url = og.URL
		datapoints++
	}

	if og.Description != "" {
		article.Description = og.Description
		datapoints++
	}

	if len(og.Images) != 0 && og.Images[0].URL != "" {
		article.Image = og.Images[0].URL
		datapoints++
	}

	return article, datapoints >= 2
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
