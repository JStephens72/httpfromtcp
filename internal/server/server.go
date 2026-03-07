package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/JStephen72/httpfromtcp/internal/request"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func New(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error listening on %s: %w", addr, err)
	}

	log.Printf("Server listening on %s\n", addr)

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %w", err)
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr().String())

		go s.handle(conn)
	}
}

func (s *Server) Serve(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s = &Server{
		listener: listener,
	}

	go s.listen()
	return s, nil

}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	_, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Invalid request: %s", err.Error())
		return
	}

	// build a simple response
	body := []byte("Hello World!\n")

	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/plain\r\n"+
			"Content-Length: %d\r\n"+
			"\r\n",
		len(body),
	)

	if _, err := conn.Write([]byte(response)); err != nil {
		log.Printf("write error: %v\n", err)
		return
	}

	if _, err := conn.Write(body); err != nil {
		log.Printf("write error: %v\n", err)
		return
	}
}
