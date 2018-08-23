package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetLinks(url string) ([]Link, error) {
	document, err := getHTML(url)
	if err != nil {
		return nil, err
	}

	var links []Link

	domainRegex := regexp.MustCompile("^(?:https?:\\/\\/)?(?:[^@\n]+@)?(?:www\\.)?([^:\\/\n?]+)")
	baseDomainArr := domainRegex.FindStringSubmatch(url)
	fmt.Println(baseDomainArr)
	baseDomain := baseDomainArr[0]

	if baseDomain == "" {
		return nil, errors.New("Warning: Could not determine base domain: " + baseDomain)
	}

	document.Find("a").Each(func(i int, selection *goquery.Selection) {
		link := Link{}

		var exists bool
		link.Url, exists = selection.Attr("href")
		if exists {
			isRelative, err := regexp.MatchString("^\\/.*\\/", link.Url)
			if err != nil {
				fmt.Print("Warning: ")
				fmt.Print(err)
				fmt.Print("\r\n")
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

	return links, nil
}

func getHTML(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Status code error: " + fmt.Sprint(res.StatusCode) + " " + fmt.Sprint(res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
