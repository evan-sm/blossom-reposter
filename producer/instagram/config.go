package main

import (
	"os"
)

var (
	amqpUrl           = os.Getenv("AMQP_URL") // amqp://guest:guest@localhost:5672/
	IGSessionID       = os.Getenv("IG_SESSIONID")
	IGQueryHash       = os.Getenv("IG_QUERYHASH")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_USER     = os.Getenv("POSTGRES_USER")
	POSTGRES_DB       = os.Getenv("POSTGRES_DB")
	POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
	tgBotTkn          = os.Getenv("TG_BOT_TKN")
	tgChanErr         = os.Getenv("TG_CHAN_ERR")
	socketURL           = os.Getenv("SOCKET_URL")
)
