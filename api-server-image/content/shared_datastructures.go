package main

import (
	"time"

	"github.com/olivere/elastic"
)

type Article struct {
	Headline    string
	Description string
	Image       string
	Content     string
	Source      string
	Url         string
	Time        time.Time
}

type BsonArticle struct {
	Headline    string    `bson:"headline"`
	Content     string    `bson:"content"`
	Description string    `bson:"description"`
	Image       string    `bson:"image"`
	Source      string    `bson:"source"`
	Url         string    `bson:"url"`
	Time        time.Time `bson:"time"`
}

type JsonArticle struct {
	Headline    string                `json:"headline"`
	Content     string                `json:"content"`
	Description string                `json:"description"`
	Image       string                `json:"image"`
	Source      string                `json:"source"`
	Url         string                `json:"url"`
	Time        time.Time             `json:"time"`
	Suggest     *elastic.SuggestField `json:"suggest_field,omitempty"`
}

type Source struct {
	Urls []string
	Name string
	Id   string
}

type Link struct {
	Url    string
	Source string
}

type Word struct {
	Word  string
	Count int
}

type WordToDate struct {
	Word  string
	Count int
	Date  string
}

type ScoredWord struct {
	Word  string
	Score float64
}
