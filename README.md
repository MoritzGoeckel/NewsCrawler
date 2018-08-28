# News platform with Kubernetes
The purpose of this project is practicing and experimenting with the microservice architecture and Kubernetes. This software finds and downloads news articles from various sources and performs some lingustical analysis on them.

# Microservices
This software consists of 14 microservices:

## API server
The API server is written in Go and provides a REST API that can be used by the frontend to retrieve data from the databases.

## Agt article
A redis instance for storing hashes of already seen articles

## Agt link
A redis instance for storing hashes of already seen links

## Article downloader
Microservice implemented in Go which receives links from the link queue and downloads the corresponding article. This article goes on the article queue

## Link downloader

## Link queue
A redis queue to hold the not yet visited links for the downloader

## Article queue
A redis queue to hold the not yet processed articles

## Processor
Software implemented in Go. It takes articles from the article queue, extracts information like keywords etc and pushes them into the two databases

## Static content server
A httpd server that serves the frontend

## Word cloud generator

## Cleaner (TODO)

## Elastic search
To make the articles searchable

## MongoDB
A MongoDB, this is the main database of the project

## Cache
A redis instance for caching
