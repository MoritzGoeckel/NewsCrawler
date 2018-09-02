# News platform with Kubernetes
The purpose of this project is practicing and experimenting with the microservice architecture and Kubernetes. This software finds and downloads news articles from various sources and performs some lingustical analysis on them

# Architecture
This software consists of 12 microservices

![microservice architecture](https://raw.githubusercontent.com/MoritzGoeckel/KubernetesNewsService/master/media/diagram.png)

### API server
The API server is written in Go. It provides a REST API that can be used by the frontend to retrieve data from the databases and serves the frontend

### 'Already got that' Redis
There are two redis instances for storing hashes of already seen articles and links

### Article downloader
Microservice implemented in Go which receives links from the link queue and downloads the corresponding article. This article goes on the article queue

### Link downloader
Software implemented in Go that visits a set of websites and extracts potentially interesting links and pushes them on the link queue

### Link queue and article queue
There are two redis instances used as queues to hold the not yet visited links and the not yet processed articles

### Processor
Software implemented in Go. It takes articles from the article queue, extracts information like keywords etc and pushes them into the two databases

### Word cloud generator
This part of the software is also implemented in Go and is responsible to calculate the important words for today

### Elasticsearch
An elasticsearch instance to make the articles searchable

### MongoDB
A MongoDB, this is the main database of the project

### Cache
A redis instance for caching the API server

### Look at me frontend
![microservice architecture](https://raw.githubusercontent.com/MoritzGoeckel/KubernetesNewsService/master/media/news-frontend.png)
