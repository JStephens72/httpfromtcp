package response

import (
	"fmt"

	"github.com/JStephen72/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
