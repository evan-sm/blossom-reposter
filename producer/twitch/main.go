package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"github.com/parnurzeal/gorequest"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Twitch WebHooks started")
	initDB()
	go serve()
	for {
		getPersonsDB()
		getTwitchOAuth()
		subscribeWebHooks()
		log.Printf("⏳ Next webhook resubscribe is in 2 hours")
		time.Sleep(7200 * time.Second)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		reportTg(err)
		log.Printf("%s: %s", msg, err)
	}
}

func getTwitchOAuth() {
	body := OAuthTokenBody{}
	body.ClientID = TwitchClientID
	body.ClientSecret = TwitchClientSecret
	body.GrantType = "client_credentials"

	resp, rbody, errs := gorequest.New().Post("https://id.twitch.tv/oauth2/token").
		Query(body).EndStruct(&appAccessTkn)
	if errs != nil {
		log.Printf("resp: %v\nbody: %v\nerrs: %v\n", resp, rbody, errs)
		reportTg(errs)
	}
	pp.Println(appAccessTkn)
	TwitchOAuthTkn = fmt.Sprintf("Bearer %v", appAccessTkn.AccessToken)
	pp.Println(TwitchOAuthTkn)
}

func serve() {
	router := gin.Default()

	// Subscription verify
	router.GET("/webhook/userfollows", func(c *gin.Context) {
		challenge := c.Query("hub.challenge")
		c.String(http.StatusOK, challenge)
		log.Printf("%v", challenge)

	})

	router.GET("/webhook/streamchanged", func(c *gin.Context) {
		challenge := c.Query("hub.challenge")
		c.String(http.StatusOK, challenge)
		log.Printf("%v", challenge)

	})

	// API endpoint
	router.POST("/webhook/userfollows", func(c *gin.Context) {
		var json userFollowsPayload
		if err := c.BindJSON(&json); err != nil {
			log.Printf("ShouldBindJSON error: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		}
		//log.Printf("new follower\n")
		pp.Printf("[%v]: new follower %v\n", json.Data[0].ToName, json.Data[0].FromName)
		pp.Printf("%v", json.Data[0])
		//log.Printf("%v %v", json.Data[0].FollowedAt, json.Data[0].FromName)
		c.String(http.StatusOK, "ok")
	})

	router.POST("/webhook/streamchanged", func(c *gin.Context) {
		var json *streamChangedPayload

		if err := c.BindJSON(&json); err != nil {
			log.Printf("BindJSON error: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			reportTg(err.Error())
			return
		}
		//var j streamChangedPayload
		log.Printf("%v\n", len(json.Data))
		if len(json.Data) == 0 {
			log.Printf("Empty webhook, skip.")
			c.String(http.StatusOK, "ok")
			return
		}
		pp.Println(json.Data[0])

		log.Printf("StreamID: %v\njson.Data[0].ID: %v\n", streamID, json.Data[0].ID)
		if twitchNotificationID == c.GetHeader("Twitch-Notification-Id") {
			log.Printf("Same notification ID: %v. Skip!", twitchNotificationID)
			c.String(http.StatusOK, "ok")
			//reportTg("Same notification ID. Skip!")
			return
		}
		if streamID == json.Data[0].ID {
			c.String(http.StatusOK, "ok")
			pp.Println("Stream changed info")
			pp.Println(json.Data[0])
			return
		}
		if twitchNotificationID != c.GetHeader("Twitch-Notification-Id") {
			c.String(http.StatusOK, "ok")
			twitchNotificationID = c.GetHeader("Twitch-Notification-Id")
			log.Printf("%v", twitchNotificationID)
			log.Printf("%v goes live! Playing %v. '%v'\n", json.Data[0].UserName, json.Data[0].GameID, json.Data[0].Title)
			json.sendAnnounce()
		}
		streamID = json.Data[0].ID
		twitchNotificationID = c.GetHeader("Twitch-Notification-Id")
	})

	router.Run(":9696") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func subscribe(id int64, topic string) {
	var CallbackURL string

	switch topic {
	case userFollowsTopic:
		topic = fmt.Sprintf("%v%v", topic, id)
		CallbackURL = fmt.Sprintf("%v/webhook/userfollows", CallbackHost)
		pp.Printf("Topic: %v\nCallback URL: %v\n", topic, CallbackURL)
	case streamChangedTopic:
		topic = fmt.Sprintf("%v%v", topic, id)
		CallbackURL = fmt.Sprintf("%v/webhook/streamchanged", CallbackHost)
		pp.Printf("Topic: %v\nCallback URL: %v\n", topic, CallbackURL)
	}

	pp.Println(topic)
	body := subscribeBody{}
	body.Mode = "subscribe"
	body.Callback = CallbackURL
	body.LeaseSeconds = 86400
	body.Secret = "pook"
	body.Topic = topic

	resp, respBody, errs := gorequest.New().Post(twitchWebhookURL).Set("Authorization", TwitchOAuthTkn).
		Set("Client-ID", TwitchClientID).
		Send(body).End()
	if errs != nil {
		log.Printf("resp: %v\nbody: %v\nerrs: %v", resp, respBody, errs)
	}
	log.Printf("resp: %v\nbody: %v\nerrs: %v", resp, respBody, errs)
}

func subscribeWebHooks() {
	// Iterate through each person
	for _, s := range persons {
		log.Printf("Trying to subscribe to %s's twitch webhooks...", s.Person)
		if s.Announce_twitch_live {
			subscribe(s.Twitch_id, streamChangedTopic)
		}
	}
}

func (j *streamChangedPayload) composeJSONPayload(s *Person) {
	jsonPayload = JsonPayload{}
	jsonPayload.Person = s.Person
	jsonPayload.RepostTelegramEnabled = s.Repost_telegram_enabled
	jsonPayload.RepostTelegramChanID = s.Repost_telegram_chan_id
	jsonPayload.RepostMakabaEnabled = s.Repost_makaba_enabled
	jsonPayload.RepostVkStatusEnabled = s.Repost_vk_status_enabled
	jsonPayload.RepostVkPageEnabled = s.Repost_vk_page_enabled
	jsonPayload.RepostVkPublicEnabled = s.Repost_vk_public_enabled
	jsonPayload.InstagramPostTimestamp = s.Instagram_post_timestamp
	jsonPayload.InstagramStoryTimestamp = s.Instagram_story_timestamp
	jsonPayload.VkPageTimestamp = s.Vk_page_timestamp
	jsonPayload.VkPublicTimestamp = s.Vk_public_timestamp
	jsonPayload.VkStatusTimestamp = s.Vk_status_timestamp
	jsonPayload.TwitchUsername = s.Twitch_username
	jsonPayload.TwitchID = s.Twitch_id
	jsonPayload.InstagramUsername = s.Instagram_username
	jsonPayload.InstagramID = s.Instagram_id
	jsonPayload.VkPageID = s.Vk_page_id
	jsonPayload.DvachBoard = "fag"
	jsonPayload.From = "twitch"
	jsonPayload.Type = "live"
	jsonPayload.Source = fmt.Sprintf("https://twitch.tv/%v", s.Twitch_username)
	jsonPayload.Language = j.Data[0].Language
	jsonPayload.Caption = j.Data[0].Title

	files = nil
	url := j.Data[0].ThumbnailURL
	url = strings.Replace(url, "{width}", "1920", -1)
	url = strings.Replace(url, "{height}", "1080", -1)
	files = append(files, url) // add .jpg to slice

	jsonPayload.Files = files

}

func (j *streamChangedPayload) sendAnnounce() {
	findPersonDB(j.Data[0].UserID)
	for _, s := range persons {
		if s.Announce_twitch_live {
			j.composeJSONPayload(s)
			pp.Print(jsonPayload)
			sent := sendJSONPayload()
			if sent {
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func sendJSONPayload() bool {
	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body, _ := json.Marshal(jsonPayload)
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Printf("%s", err)
		return false
	}
	log.Printf(" [x] Sent to rMQ")
	return true
}
