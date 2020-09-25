package main

import (
	"fmt"
	"log"
    "os"
	"strconv"

	"github.com/parnurzeal/gorequest"
)

type InputMedia struct {
	Type    string `json:"type"`
	Media   string `json:"media"`
	Caption string `json:"caption"`
}

func reportTg(e interface{}) {
	log.Printf("%v", e)
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", tgBotTkn)
	query := fmt.Sprintf(`{"chat_id": %s, "text":"%v"}`, tgChanErr, e)
	gorequest.New().Get(url).Send(query).End()
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}


func failOnError(err error, msg string) {
	if err != nil {
		reportTg(err)
		log.Printf("%s: %s", msg, err)
	}
}

func convStrInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
