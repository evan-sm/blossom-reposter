package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		reportTg(err)
	}
}

func main() {
	//log.Println(amqpUrl)
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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	i := false
	go func() {
		for {
			if i {
				log.Println("Cycle")
				i = false
				if jsonPayload.RepostTelegramEnabled == true {
					repostTg()
					jsonPayload = JsonPayload{}
				} else {
					log.Printf("Reposting to Telegram is disabled: %v", jsonPayload.RepostTelegramEnabled)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()

	forever := make(chan bool)

	go func() {
		log.Printf("Telegram consumer ready, PID: %d", os.Getpid())
		for d := range msgs {
			log.Printf("%v | Received a message: %s", reflect.TypeOf(d.Body), d.Body)
			if IsJSON(d.Body) == true {
				err := json.Unmarshal(d.Body, &jsonPayload)
				if err != nil {
					log.Println(err, jsonPayload)
					reportTg(err)
				}
				log.Println("Json validatd")
				i = true
				//d.Ack(true) // remove anyway
				time.Sleep(20 * time.Second)
				log.Println("Sleep for 20 sec...")
			} else {
				log.Println("This is not Json message or not validated")
				//d.Ack(true)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func IsJSON(jsByte []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(jsByte, &js) == nil
}
