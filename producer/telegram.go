package main

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
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
