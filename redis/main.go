package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var PORT = ":4000"

var rds = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})

func Healthcheck(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte("Server is healthy"))
}

func SavetoRedis(w http.ResponseWriter, req *http.Request) {

	value := req.URL.Query().Get("name")

	if value == "" {
		w.Write([]byte("query name is empty"))
		return
	}

	err := rds.Set(ctx, "name", value, time.Minute).Err()
	if err != nil {
		w.Write([]byte("Error saving to redis"))
		return
	}

	username, err := rds.Get(ctx, "name").Result()
	if err != nil {
		w.Write([]byte("Failed to get the key"))
	}
	w.Write([]byte("Your name " + username + " is set in the redis DB"))
}

func main(){

	fmt.Println("Connected ", rds)
	fmt.Printf("Go server running on Port %v", PORT)

	http.HandleFunc("/", Healthcheck)
	http.HandleFunc("/save-to-redis", SavetoRedis)


	http.ListenAndServe(PORT, nil)
}