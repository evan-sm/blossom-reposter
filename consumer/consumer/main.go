package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	//"reflect"
	"time"

	"github.com/wMw9/rpdb"
	"github.com/wMw9/rpps"
)

var jsonPayload rpdb.JsonPayload

func main() {
	consume()
}

func consume() {
	c := rpps.NewClientSocket()
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
		}
		if jsonPayload.RepostMakabaEnabled {
			repost2ch()
		}
		jsonPayload = rpdb.JsonPayload{}
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
