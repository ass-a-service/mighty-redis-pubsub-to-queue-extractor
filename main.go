package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Event struct {
	Version        int                    `json: "v"`
	Name           string                 `json:"event_name"`
	MessageID      string                 `json:"message_id"`
	GuildID        string                 `json:"guild_id"`
	AuthorID       string                 `json:author_id"`
	AdditionalData map[string]interface{} `json:"additional_data"`
}

type XReadGroupArgs struct {
	Group    string
	Consumer string
	Streams  []string // list of streams and ids, e.g. stream1 stream2 id1 id2
	Count    int64
	Block    time.Duration
	NoAck    bool
}

var ctx = context.Background()

func GetEnvOrFail(env_var_name string) string {
	env_var, err := os.LookupEnv((env_var_name))

	if !err {
		log.Fatalln("Could not find " + env_var_name + " environment variable. Does it exist?")
	}

	return env_var
}

func parseMessage(redis_message redis.XMessage) Event {
	var event Event
	for field_name, field_value := range redis_message.Values {
		if field_name == "event_name" {
			event.Name = fmt.Sprintf("%v", field_value)
		} else if field_name == "author_id" {
			event.AuthorID = fmt.Sprintf("%v", field_value)
		}
	}
	return event
}

func PubSub2Queue() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     GetEnvOrFail("REDIS_HOST") + ":" + GetEnvOrFail("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//rdb.XGroupCreate(ctx, "bus:ajo", "go", "0")
	args := redis.XReadGroupArgs{
		Group:    "go",
		Consumer: "si",
		Streams:  []string{"bus:ajo", ">"},
		Count:    10000,
		Block:    0 * time.Millisecond,
	}

	for {
		res, _ := rdb.XReadGroup(ctx, &args).Result()
		fmt.Println("Estamos leyendo")

		for _, s := range res {
			for _, k := range s.Messages {
				message := parseMessage(k)
				if message.Name == "farm" {
					rdb.ZIncrBy(ctx, "lb", 1, message.AuthorID)
				}
				rdb.XAck(ctx, "bus:ajo", "go", k.ID)
			}
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
	fmt.Fprintf(os.Stdout, "hasta luego cowboy del espacio")
}
