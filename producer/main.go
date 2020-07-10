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
		time.Sleep(2 * time.Second)
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
		if s.Check_instagram == true {
			log.Printf("Checking %s's instagram...", s.Person)
			checkInstagramPost(s.Person, s.Instagram_username, s.Instagram_post_timestamp)
			checkInstagramStory(s.Person, s.Instagram_username,
				s.Instagram_id, s.Instagram_story_timestamp)
		}
		if s.Check_telegram == true {
			log.Printf("Checking %s's telegram, channel ID: %v...", s.Person, s.Telegram_chan_id)
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

	q, err := ch.QueueDeclare(
		"add",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body, _ := json.Marshal(jsonPayload)
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
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
