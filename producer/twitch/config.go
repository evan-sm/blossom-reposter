package main

import (
	"os"
)

var (
	amqpUrl            = os.Getenv("AMQP_URL") // amqp://guest:guest@localhost:5672/
	IGSessionID        = os.Getenv("IG_SESSIONID")
	IGQueryHash        = os.Getenv("IG_QUERYHASH")
	POSTGRES_PASSWORD  = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_USER      = os.Getenv("POSTGRES_USER")
	POSTGRES_DB        = os.Getenv("POSTGRES_DB")
	POSTGRES_HOST      = os.Getenv("POSTGRES_HOST")
	tgBotTkn           = os.Getenv("TG_BOT_TKN")
	tgChanErr          = os.Getenv("TG_CHAN_ERR")
	vkAccessTkn        = os.Getenv("VK_ACCESS_TKN")
	TwitchOAuthTkn     = os.Getenv("TWITCH_OAUTH_TKN")     // "Bearer b57agl4irthna2uhgpbimm028kwgadv"
	TwitchClientID     = os.Getenv("TWITCH_CLIENT_ID")     // "sdc9vfkrfzmtqaodkdipni1k1ccb4m"
	TwitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET") // "sdc9vfkrfzmtqaodkdipni1k1ccb4m"
	CallbackHost       = os.Getenv("CALLBACK_HOST")        //"https://9323c95d112e.ngrok.io"
)

const (
	twitchWebhookURL   = "https://api.twitch.tv/helix/webhooks/hub"
	userFollowsTopic   = "https://api.twitch.tv/helix/users/follows?first=1&to_id=" // "https://api.twitch.tv/helix/users/follows?first=1&to_id=172230472"
	streamChangedTopic = "https://api.twitch.tv/helix/streams?user_id="             // "https://api.twitch.tv/helix/streams?user_id=172230472"
    charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
