package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"time"

	"github.com/wmw9/blossom-reposter/pkg/pubsub"

	//"go.nanomsg.org/mangos/v3"
	//"go.nanomsg.org/mangos/v3/protocol/sub"
	//_ "go.nanomsg.org/mangos/v3/transport/all"
)

var c pubsub.ClientSocket

func main() {
    Start()
}

func Start() {

    c = pubsub.NewClientSocket()
    fmt.Println(reflect.TypeOf(c))
    log.Printf("%v", c)
	c.ClientInit(socketUrl, "telegram consumer")
	//c.ClientStart()

    log.Printf("%v", c)
    var err error
    var msg []byte

    for {
        if msg, err = c.ClientStart(); err != nil {
            die("Cannot recv: %s", err.Error())
        }
        log.Printf("Client: %s: %v | Received a message: %s", msg, reflect.TypeOf(msg), msg)
		err = json.Unmarshal(msg, &jsonPayload)

        if err != nil {
			log.Println(err, jsonPayload)
			reportTg(err)
			die("Can't unmarshal msg: %s", err.Error())
		}

		if jsonPayload.RepostTelegramEnabled == true {
			repostTg()
			jsonPayload = JsonPayload{}
		} else {
			log.Printf("Reposting to Telegram is disabled: %v", jsonPayload.RepostTelegramEnabled)
		}

		log.Println("Sleep for 20 sec...")
		time.Sleep(2 * time.Second)
    }
}


//func (s *mangosSocket) Start() {
//	var err error
//	for {
//		if s.msg, err = s.sock.Recv(); err != nil {
//			die("Cannot recv: %s", err.Error())
//		}
//		//fmt.Printf("CLIENT(%s): RECEIVED %s\n", name, string(msg))
//		log.Printf("Client: %s: %v | Received a message: %s", s.name, reflect.TypeOf(s.msg), s.msg)
//
//		err = json.Unmarshal(s.msg, &jsonPayload)
//		if err != nil {
//			log.Println(err, jsonPayload)
//			reportTg(err)
//			die("Can't unmarshal msg: %s", err.Error())
//		}
//
//		if jsonPayload.RepostTelegramEnabled == true {
//			repostTg()
//			jsonPayload = JsonPayload{}
//		} else {
//			log.Printf("Reposting to Telegram is disabled: %v", jsonPayload.RepostTelegramEnabled)
//		}
//
//		log.Println("Sleep for 20 sec...")
//		time.Sleep(2 * time.Second)
//
//		err = json.Unmarshal(s.msg, &jsonPayload)
//		if err != nil {
//			log.Println(err, jsonPayload)
//			reportTg(err)
//		}
//	}
//}



func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		reportTg(err)
		die("Something went wrong")
	}
}

/* func client(url string, name string) {
	var sock mangos.Socket
	var err error
	var msg []byte

	if sock, err = sub.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		die("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	log.Printf("[*] Telegram Consumer - Waiting for messages. To exit press CTRL+C")

	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		//fmt.Printf("CLIENT(%s): RECEIVED %s\n", name, string(msg))
		log.Printf("Client: %s: %v | Received a message: %s", name, reflect.TypeOf(msg), msg)

		err = json.Unmarshal(msg, &jsonPayload)
		if err != nil {
			log.Println(err, jsonPayload)
			reportTg(err)
			die("Can't unmarshal msg: %s", err.Error())
		}

		if jsonPayload.RepostTelegramEnabled == true {
			repostTg()
			jsonPayload = JsonPayload{}
		} else {
			log.Printf("Reposting to Telegram is disabled: %v", jsonPayload.RepostTelegramEnabled)
		}

		log.Println("Sleep for 20 sec...")
		time.Sleep(2 * time.Second)

		err = json.Unmarshal(msg, &jsonPayload)
		if err != nil {
			log.Println(err, jsonPayload)
			reportTg(err)
		}
	}
} */
