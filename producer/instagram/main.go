package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	initDB()
	for {
		getPersonsDB()
		checkSN()
		log.Printf("⏳ Next run is in 12 minutes...")
		time.Sleep(720 * time.Second)
		//time.Sleep(60 * time.Second)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		reportTg(err)
		log.Printf("%s: %s", msg, err)
	}
}

func checkSN() {
	// Iterate through each person
	for _, s := range persons {
		log.Printf("Checking %s's instagram...", s.Person)
		clearJSONPayload()
		jsonPayload.Person = s.Person
		jsonPayload.RepostTelegramEnabled = s.Repost_telegram_enabled
		jsonPayload.RepostTelegramChanID = s.Repost_telegram_chan_id
		jsonPayload.RepostMakabaEnabled = s.Repost_makaba_enabled
		jsonPayload.InstagramPostTimestamp = s.Instagram_post_timestamp
		jsonPayload.InstagramStoryTimestamp = s.Instagram_story_timestamp
		jsonPayload.InstagramUsername = s.Instagram_username
		jsonPayload.InstagramID = s.Instagram_id
		jsonPayload.DvachBoard = "fag"
		if s.Check_instagram_post == true {
			checkInstagramPost()
		}
		if s.Check_instagram_post == true {
			checkInstagramStory()
		}
	}
}

func sendJsonPayload() bool {
	conn, err := amqp.Dial(amqpUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/* 	q, err := ch.QueueDeclare(
	   		"vk",
	   		true,
	   		false,
	   		false,
	   		false,
	   		nil,
	   	)
	   	failOnError(err, "Failed to declare a queue") */

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
	/* 	err = ch.Publish(
	"",
	q.Name,
	false,
	false,
	amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}) */
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
	log.Printf(" [x] Sent %s", body)
	return true
}
