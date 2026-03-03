package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const remoteAddr = "localhost:42069"

func main() {
	endPoint, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		log.Fatalf("error resolving %s: %v", remoteAddr, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, endPoint)
	if err != nil {
		log.Fatalf("error creating connection to %v: %v", endPoint.Port, err.Error())
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error reading input from StdIO: %v", err.Error())
		}

		n, err := conn.Write([]byte(input))
		if err != nil {
			log.Printf("error writing to connection: %v", err.Error())
		}
		log.Printf("%v bytes written", n)

	}
}
