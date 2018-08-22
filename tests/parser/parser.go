package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParsePage(url string) []Link {
	document := getHTML(url)
	var links []Link

	domainRegex := regexp.MustCompile("^(?:https?:\\/\\/)?(?:[^@\n]+@)?(?:www\\.)?([^:\\/\n?]+)")
	baseDomainArr := domainRegex.FindStringSubmatch(url)
	fmt.Println(baseDomainArr)
	baseDomain := baseDomainArr[0]

	if baseDomain == "" {
		log.Fatal("Could not determine base domain: " + baseDomain)
	}

	document.Find("a").Each(func(i int, selection *goquery.Selection) {
		link := Link{}

		var exists bool
		link.Url, exists = selection.Attr("href")
		if exists {
			isRelative, err := regexp.MatchString("^\\/.*\\/", link.Url)
			if err != nil {
				log.Fatal(err)
			}

			//fmt.Printf("Url: %s isRelative: %s isNormal: %s\n", link.Url, isRelative, isNormalUrl)

			if isRelative || strings.Contains(link.Url, baseDomain) {
				if isRelative {
					link.Url = baseDomain + link.Url
				}

				if len(link.Url) > 40 {
					links = append(links, link)
				}
			}
		}
	})

	return links
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

func main() {
	links := ParsePage("http://spiegel.de")
	for _, a := range links {
		fmt.Println(a)
	}
	fmt.Println("eop")
}
