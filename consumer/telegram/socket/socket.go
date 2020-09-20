package socket

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"go.nanomsg.org/mangos/protocol/pub"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

type clientSocket struct {
	sock mangos.Socket
	msg  []byte
	url  string
	name string
}

type serverSocket struct {
	sock mangos.Socket
	msg  []byte
	url  string
	name string
}

func NewServerSocket() mangosSocket {
	return serverSocket{}
}

func NewClientSocket() mangosSocket {
	return clientSocket{}
}

func (c *mangosSocket) ClientStart() {
	var err error
	for {
		if c.msg, err = s.sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		//fmt.Printf("CLIENT(%s): RECEIVED %s\n", name, string(msg))
		log.Printf("Client: %s: %v | Received a message: %s", c.name, reflect.TypeOf(c.msg), c.msg)

		err = json.Unmarshal(c.msg, &jsonPayload)
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

		err = json.Unmarshal(c.msg, &jsonPayload)
		if err != nil {
			log.Println(err, jsonPayload)
			reportTg(err)
		}
	}
}

func (c *mangosSocket) ClientInit(url string, name string) {
	var err error
	c.url = url
	c.name = name

	if c.sock, err = sub.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}
	if err = c.sock.Dial(url); err != nil {
		die("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = c.sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	log.Printf("[*] Telegram Consumer - Waiting for messages. To exit press CTRL+C")

}

func (s *mangosSocket) ServerInit(url string, name string) {
	var err error

	//log.Println()
	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(url); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
}

func (s *mangosSocket) ServerSend(body []byte) {
	if err := sock.Send([]byte(body)); err != nil {
		die("Failed publishing: %s", err.Error())
	}
	log.Printf(" [x] Sent via tcp socket")
}
func die(format string, v ...interface{}) {
	log.Fprintln(os.Stderr, log.Sprintf(format, v...))
	os.Exit(1)
}
