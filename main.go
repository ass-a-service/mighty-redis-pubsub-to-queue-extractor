package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func GetEnvOrFail(env_var_name string) string {
	env_var, err := os.LookupEnv((env_var_name))

	if !err {
		log.Fatalln("Could not find " + env_var_name + " environment variable. Does it exist?")
	}

	return env_var
}

func PubSub2Queue() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetEnvOrFail("REDIS_HOST") + ":" + GetEnvOrFail("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redis_channel_name := GetEnvOrFail("REDIS_CHANNEL")
	redis_queue_name := GetEnvOrFail("REDIS_QUEUE")
	pubsub := rdb.Subscribe(ctx, redis_channel_name)
	ch := pubsub.Channel()
	fmt.Fprintln(os.Stdout, "Subscribed to channel: "+redis_channel_name)
	for msg := range ch {
		rdb.RPush(ctx, redis_queue_name, msg.Payload)
		if os.Getenv("DEBUG") == "1" {
			fmt.Println(msg.Channel, msg.Payload)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when loading .env file: "+err.Error())
	} else {
		fmt.Fprintln(os.Stdout, ".env loaded successfully")
	}
	PubSub2Queue()
}
