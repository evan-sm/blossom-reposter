package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	//"reflect"
	"time"

	"github.com/wmw9/blossom-reposter/pkg/pubsub"

)


func main() {
    consume()
}

func consume() {
    c := pubsub.NewClientSocket()
	c.ClientInit(socketUrl, "Telegram Consumer")

    var err error
    var msg []byte

    for {
        if msg, err = c.ClientStart(); err != nil {
            die("Cannot recv: %s", err.Error())
        }

		if err = json.Unmarshal(msg, &jsonPayload); err != nil {
    		log.Println(err, jsonPayload)
			reportTg(err)
			die("Can't unmarshal msg: %s", err.Error())
        }

		if jsonPayload.RepostTelegramEnabled {
			repostTg()
			jsonPayload = JsonPayload{}
		} else {
			log.Printf("Reposting to Telegram is disabled: %v", jsonPayload.RepostTelegramEnabled)
		}

		log.Println("Sleep for a sec...")
		time.Sleep(time.Second)
    }
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		reportTg(err)
		die("Something went wrong")
	}
}

