package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("Downloader version 0.01")

	cache := getRedisConnection()

	fmt.Println("eop")
}

func getRedisConnection() *redis.Client {
	cacheUrl := os.Getenv("cache-redis-url")

	fmt.Println("cache url: " + cacheUrl)

	cacheClient := redis.NewClient(&redis.Options{
		Addr:     cacheUrl + ":6379",
		Password: "",
		DB:       0,
	})

	return cacheClient
}
