package main

import (
	"github.com/olivere/elastic"
)

type Article struct {
	Headline    string
	Description string
	Image       string
	Content     string
	Source      string
	Url         string
	Language    string
	DateTime    int64
}

type BsonArticle struct {
	Headline    string `bson:"headline"`
	Content     string `bson:"content"`
	Description string `bson:"description"`
	Image       string `bson:"image"`
	Source      string `bson:"source"`
	Url         string `bson:"url"`
	Language    string `bson:"language"`
	DateTime    int64  `bson:"datetime"`
}

type JsonArticle struct {
	Headline    string                `json:"headline"`
	Content     string                `json:"content"`
	Description string                `json:"description"`
	Image       string                `json:"image"`
	Source      string                `json:"source"`
	Url         string                `json:"url"`
	Language    string                `json:"language"`
	DateTime    int64                 `json:"datetime"`
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
