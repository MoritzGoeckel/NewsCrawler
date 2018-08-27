// go get -u github.com/go-redis/redis

package main

import "github.com/go-redis/redis"

func main(){
    url := "agt-redis.default.svc.cluster.local"

  client := redis.NewClient(&redis.Options{
    Addr:     url + ":6379",
    Password: "",
    DB:       0,
  })

  err := client.Set("thekey", "thevalue", 0).Err()
  if err != nil {
      panic(err)
  }

  /*
  val, err := client.Get("key").Result()
  if err == redis.Nil {
    fmt.Println("key does not exist")
  } else if err != nil {
    panic(err)
  } else {
    fmt.Println("key", val)
  }
  */
}
