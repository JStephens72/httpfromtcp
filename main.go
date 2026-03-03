package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const listenerPort = ":42069"

func main() {
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		log.Fatalf("err listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s", listenerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection accept error: %v", err.Error())
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr().String(), " closed")
	}
}

// getLinesChannel reads from f in a goroutine and returns a channel of lines.
func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		currentLineContents := ""

		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				// Emit any remaining partial line before exiting
				if currentLineContents != "" {
					out <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					return
				}
				// For non-EOF errors, we still stop reading
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")

			// Emit all complete lines
			for i := 0; i < len(parts)-1; i++ {
				out <- currentLineContents + parts[i]
				currentLineContents = ""
			}

			// Last part is either a full line fragment or empty
			currentLineContents += parts[len(parts)-1]
		}
	}()

	return out
}
