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
    for {
    
    
        log.Println(amqpUrl)
        conn, err := amqp.Dial(amqpUrl)
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()
        notify := conn.NotifyClose(make(chan *amqp.Error)) //error channel

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
        for {
            select {
            case err = <-notify:
                break
            case d = <- msgs:
			    log.Printf("%v | Received a message: %s", reflect.TypeOf(d.Body), d.Body)
            }
        }
    }

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}

func IsJSON(jsByte []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(jsByte, &js) == nil
}
