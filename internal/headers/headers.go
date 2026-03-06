package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// headers are done, consume the CRLF
		return 2, true, nil
	}
	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	fieldName := string(parts[0])

	if fieldName != strings.TrimRight(fieldName, " ") {
		return 0, false, fmt.Errorf("invalid header name: '%s'", fieldName)
	}

	fieldValue := bytes.TrimSpace(parts[1])
	fieldName = strings.TrimSpace(fieldName)

	h.Set(fieldName, string(fieldValue))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
