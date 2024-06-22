package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := redis.NewClient(&redis.Options{Addr: "172.28.18.24:6379"})
	cmd := client.Set("foo", "bar", time.Second*10)
	log.Println(cmd.Result())
}
