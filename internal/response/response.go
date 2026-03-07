package response

import (
	"fmt"
	"io"

	"github.com/JStephen72/httpfromtcp/internal/headers"
)

type Response struct {
	RequestLine StatusLine
	Headers     headers.Headers
	Body        []byte
	state       responseState
}

type StatusLine struct {
	HttpVersion  string
	StatusCode   string
	ReasonPhrase string
}

type responseState int

const (
	responseStateInitialized responseState = iota
	responseStateDone
)

type StatusCode int

const (
	STATUS_OK                    StatusCode = 200
	STATUS_BAD_REQUEST           StatusCode = 400
	STATUS_INTERNAL_SERVER_ERROR StatusCode = 500
)

var StatusMap = map[StatusCode]string{
	STATUS_OK:                    "OK",
	STATUS_BAD_REQUEST:           "Bad Request",
	STATUS_INTERNAL_SERVER_ERROR: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, StatusMap[statusCode])
	if _, err := w.Write([]byte(statusLine)); err != nil {
		return fmt.Errorf("error writing response status")
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": fmt.Sprintf("%d", contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerList := ""
	for key, value := range headers {
		headerList += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	headerList += "\r\n"
	if _, err := w.Write([]byte(headerList)); err != nil {
		return fmt.Errorf("error writing response headers")
	}
	return nil
}
