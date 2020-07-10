package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"
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

	q, err := ch.QueueDeclare(
		"add", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(1, 0, false)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
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
				repost2ch()
				repostTg()
			}
			time.Sleep(3 * time.Second)
		}
	}()

	forever := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
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
				d.Ack(true) // remove anyway
				time.Sleep(20 * time.Second)
				log.Println("Sleep for 20 sec...")
			} else {
				log.Println("This is not Json message or not validated")
				d.Ack(true)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func repost2ch() bool {
	board, thread := findThread()
	log.Printf("https://2ch.hk/%v/res/%v.html", board, thread)
	valuesBase := prepareBase(board, thread)
	valuesFiles := prepareFiles()

	client, ok := customClient()
	if ok == false {
		return false
	}
	//fmt.Println("valuesFiles type is:", reflect.TypeOf(valuesFiles))
	err, success, num := makabaPost(client, postingUrl, valuesBase, valuesFiles)
	if err != nil {
		log.Println(err)
		reportTg(err)
	}
	if success {
		log.Printf("%v", num)
	}
	return success
}

func makabaPost(client *http.Client, url string, valuesBase map[string]io.Reader, valuesFiles map[string]io.Reader) (err error, success bool, num float64) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range valuesBase {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if fw, err = w.CreateFormField(key); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false, num
		}

	}
	for key, r := range valuesFiles {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if fw, err = w.CreateFormFile(key, ""); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false, num
		}

	}
	w.Close()

	// Prepare handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Высрать в тред
	res, err := client.Do(req)
	if err != nil {
		log.Println("client.Do(req) error:", err)
		reportTg(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error:", err)
		reportTg(err)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if result["Error"] != nil {
		log.Println("Makaba post error:", result)
		reportTg(result)
	}
	log.Println(result)
	if result["Error"] == nil {
		log.Println("Successfully made post 👌🏻")
		success = true
		num = result["Num"].(float64)
		log.Printf("%v", result["Num"])
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		reportTg(err)
	}
	return err, success, num
}

func IsJSON(jsByte []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(jsByte, &js) == nil
}
