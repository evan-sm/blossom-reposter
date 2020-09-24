package pubsub

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

type ClientSocket struct {
	sock mangos.Socket
	msg  []byte
	url  string
	name string
}

type ServerSocket struct {
	sock mangos.Socket
	msg  []byte
	url  string
	name string
}

func NewServerSocket() ServerSocket {
	return ServerSocket{}
}

func NewClientSocket() ClientSocket {
	return ClientSocket{}
}

func (c *ClientSocket) ClientStart() ([]byte, error) {
	var err error
	log.Printf("%v [*] Waiting for messages. To exit press CTRL+C", c.name)
	if c.msg, err = c.sock.Recv(); err != nil {
		die("Cannot recv: %s", err.Error())
	}
	log.Printf("Client: %s: %v | Received a message: %s", c.name, reflect.TypeOf(c.msg), c.msg)

	
    return c.msg, err
}

func (c *ClientSocket) ClientInit(url string, name string) {
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

}

func (s *ServerSocket) ServerInit(url string, name string) {
	var err error

	if s.sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = s.sock.Listen(url); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
}

func (s *ServerSocket) ServerSend(body []byte) {
	if err := s.sock.Send([]byte(body)); err != nil {
		die("Failed publishing: %s", err.Error())
	}
	log.Println(" [x] Sent via TCP")
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

