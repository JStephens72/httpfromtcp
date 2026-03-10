package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/JStephen72/httpfromtcp/internal/headers"
	"github.com/JStephen72/httpfromtcp/internal/request"
	"github.com/JStephen72/httpfromtcp/internal/response"
	"github.com/JStephen72/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	srv, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		videoHandler(w, req)
		return
	}
	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kind sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	url := ""
	if strings.HasPrefix(target, "/httpbin/") {
		trimmed := strings.TrimPrefix(target, "/httpbin")
		url = "https://httpbin.org" + trimmed
	}

	fmt.Println("Proxying to ", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeSuccess)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Add("Trailer", "X-Content-SHA256")
	h.Add("Trailer", "X-Content-Length")
	w.WriteHeaders(h)

	fullBody := make([]byte, 0)

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Printf("Read %d bytes\n", n)
		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Printf("Error writing chunked body: %v\n", err)
				break
			}
			fullBody = append(fullBody, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error reading response body: %v\n", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("Error writing chunked body done: %v\n", err)
	}

	trailer := headers.NewHeaders()
	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailer.Set("X-Content-SHA256", sha256)
	trailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	err = w.WriteTrailers(trailer)
	if err != nil {
		fmt.Printf("Error writing trailers: %v\n", err)
	}
	fmt.Println("Wrote trailers")
}

func videoHandler(w *response.Writer, _ *request.Request) {
	if err := w.WriteStatusLine(response.StatusCodeSuccess); err != nil {
		fmt.Printf("error writing status line: %v\n", err)
	}

	const filepath = "assets/vim.mp4"
	dat, err := os.ReadFile(filepath)
	if err != nil {
		handler500(w, nil)
		return
	}

	h := response.GetDefaultHeaders(len(dat))
	h.Set("Content-Type", "video/mp4")
	if err := w.WriteHeaders(h); err != nil {
		fmt.Printf("error writing headers: %v\n", err)
	}
	if _, err := w.WriteBody(dat); err != nil {
		fmt.Printf("error writing body: %v\n", err)
	}
}
