package main

import (
	"log"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (p *HelloService) Hello(request string, response *string) error {
	*response = "Hello " + request
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	for {
		connection, err2 = listener.Accept()
		if err2 != nil {
			log.Fatal(err2)
		}

		go rpc.ServeConn()
	}
}
