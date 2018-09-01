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

/*document.Find("h1").Each(func(i int, selection *goquery.Selection) {

})*/

func GetArticle(document *goquery.Document, url string, source string) (Article, error) {
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

	//Source

	article.Source = source
	datapointsStr += "Source "
	datapoints++

	//Title

	if og.Title != "" {
		article.Headline = og.Title
		datapointsStr += "Title "
		datapoints++
	} else {
		article.Headline = document.Find("h1").First().Text()
		if article.Headline != "" {
			datapointsStr += "Title "
			datapoints++
		}
	}

	//URL

	if og.URL != "" {
		article.Url = og.URL
	} else {
		article.Url = url
	}
	datapointsStr += "URL "
	datapoints++

	//Description

	if og.Description != "" {
		article.Description = og.Description
		datapointsStr += "Description "
		datapoints++
	} else {
		description, exists := document.Find("meta[name=description]").First().Attr("content")
		if exists {
			article.Description = description
			datapointsStr += "Description "
			datapoints++
		}
	}

	//Image

	if len(og.Images) != 0 && og.Images[0].URL != "" {
		article.Image = og.Images[0].URL
		datapointsStr += "Image "
		datapoints++
	} else {
		image, exists := document.Find("meta[name=twitter:image]").First().Attr("content")
		if exists {
			article.Image = image
			datapointsStr += "Image "
			datapoints++
		}
		//search for biggest image on site
	}

	//Search for content
	//Should we add keywords?

	if datapoints <= 3 {
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
