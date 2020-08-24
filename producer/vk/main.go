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
		log.Printf("⏳ Next run is in 2 minutes...")
		time.Sleep(120 * time.Second)
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
		log.Printf("Checking %s's VK...", s.Person)
		//UsersGet()
		if s.Repost_vk_page_enabled {
			composeJSONPayload(s)
			if s.Check_vk_page == true {
				WallGet(s.Vk_page_id)
			}
		}
		if s.Repost_vk_public_enabled {
			composeJSONPayload(s)
			if s.Check_vk_public == true {
				WallGet(s.Vk_public_id)
			}
		}
		if s.Repost_vk_status_enabled {
			composeJSONPayload(s)
			if s.Check_vk_status == true {
				UsersGet(s.Vk_page_id, s.Vk_status_text)
			}
		}
	}
}

func composeJSONPayload(s *Person) {
	jsonPayload = JsonPayload{}
	files = nil
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
	jsonPayload.InstagramUsername = s.Instagram_username
	jsonPayload.InstagramID = s.Instagram_id
	jsonPayload.VkPageID = s.Vk_page_id
	jsonPayload.DvachBoard = "fag"
	jsonPayload.From = "vk"
}

func sendJSONPayload() bool {
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
