package main

import (
	"os"
	//    "fmt"
)

const (
	makabaUrl  = "https://2ch.hk/makaba/makaba.fcgi"
	postingUrl = "https://2ch.hk/makaba/posting.fcgi?json=1"
)

var (
	passcode  = os.Getenv("PASSCODE") // https://2ch.hk/2ch/
	tripcode  = os.Getenv("TRIPCODE") // https://2ch.hk/2ch/
	amqpUri   = os.Getenv("AMQP_URL") // amqp://guest:guest@localhost:5672/
	amqpUrl   = os.Getenv("AMQP_URL") // amqp://guest:guest@localhost:5672/
	tgBotTkn  = os.Getenv("TG_BOT_TKN")
	tgChanErr = os.Getenv("TG_CHAN_ERR")
)
