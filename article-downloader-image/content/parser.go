package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dyatlov/go-opengraph/opengraph"
)

func GetArticle(document *goquery.Document) (Article, error) {
	/*document.Find(".headline").Each(func(i int, selection *goquery.Selection) {

	})*/

	html, err := document.Html()
	if err != nil {
		return Article{}, err
	}

	og := opengraph.NewOpenGraph()

	err = og.ProcessHTML(strings.NewReader(html))
	if err != nil {
		return Article{}, err
	}

	article := Article{Time: time.Now()}

	datapoints := 0
	datapointsStr := ""

	if og.Title != "" {
		article.Headline = og.Title
		datapointsStr += "Title "
		datapoints++
	}

	if og.URL != "" {
		article.Url = og.URL
		datapointsStr += "URL "
		datapoints++
	}

	if og.Description != "" {
		article.Description = og.Description
		datapointsStr += "Description "
		datapoints++
	}

	if len(og.Images) != 0 && og.Images[0].URL != "" {
		article.Image = og.Images[0].URL
		datapointsStr += "Image "
		datapoints++
	}

	if datapoints <= 2 {
		return Article{}, errors.New("Not enough datapoints: " + fmt.Sprint(datapoints) + " -> " + datapointsStr)
	}

	return article, nil
}

func GetHTML(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Status code error: " + fmt.Sprint(res.StatusCode) + " " + fmt.Sprint(res.Status) + "%s")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
