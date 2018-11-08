package main

import (
	"github.com/olivere/elastic"
)

type Article struct {
	Headline          string
	Description       string
	Image             string
	Content           string
	Source            string
	Url               string
	Language          string
	DateTime          int64
	ArticlePerplexity float64
}

type BsonArticle struct {
	Headline          string  `bson:"headline"`
	Content           string  `bson:"content"`
	Description       string  `bson:"description"`
	Image             string  `bson:"image"`
	Source            string  `bson:"source"`
	Url               string  `bson:"url"`
	Language          string  `bson:"language"`
	DateTime          int64   `bson:"datetime"`
	ArticlePerplexity float64 `bson:"article_perplexity"`
}

type JsonArticle struct {
	Headline          string                `json:"headline"`
	Content           string                `json:"content"`
	Description       string                `json:"description"`
	Image             string                `json:"image"`
	Source            string                `json:"source"`
	Url               string                `json:"url"`
	Language          string                `json:"language"`
	DateTime          int64                 `json:"datetime"`
	ArticlePerplexity float64               `json:"article_perplexity"`
	Suggest           *elastic.SuggestField `json:"suggest_field,omitempty"`
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

type Ngram2Words struct {
	Ngram2Words []Headline
}

type BsonNgram2Words struct {
	Ngram2Words []BsonHeadline `bson:"ngram2words"`
}

type JsonNgram2Words struct {
	Ngram2Words []JsonHeadline `json:"ngram2words"`
}

type Headline struct {
	Ngram             string
	Count             int
	FrequentNeighbors map[string]int
}

type BsonHeadline struct {
	Ngram             string         `bson:"ngram"`
	Count             int            `bson:"count"`
	FrequentNeighbors map[string]int `bson:"frequent_neighbors"`
}

type JsonHeadline struct {
	Ngram             string         `json:"ngram"`
	Count             int            `json:"count"`
	FrequentNeighbors map[string]int `json:"frequent_neighbors"`
}
