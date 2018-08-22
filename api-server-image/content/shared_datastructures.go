package main

import (
	"time"

	"github.com/olivere/elastic"
)

type Article struct {
	Headline string
	Content  string
	Source   string
	Url      string
	Time     time.Time
}

type BsonArticle struct {
	Headline string    `bson:"headline"`
	Content  string    `bson:"content"`
	Source   string    `bson:"source"`
	Url      string    `bson:"url"`
	Time     time.Time `bson:"time"`
}

type JsonArticle struct {
	Headline string                `json:"headline"`
	Content  string                `json:"content"`
	Source   string                `json:"source"`
	Url      string                `json:"url"`
	Time     time.Time             `json:"time"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

type Source struct {
	Urls []string
	Name string
	Id   string
}
