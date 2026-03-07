package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JStephen72/httpfromtcp/internal/request"
)

const listenerPort = ":42069"

func main() {
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		log.Fatalf("err listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", listenerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection accept error: %v", err.Error())
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

		req, err := request.RequestFromReader(conn)

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %v: %v\n", key, value)
		}
	}
}
