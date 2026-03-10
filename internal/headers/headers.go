package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const CRLF = "\r\n"

var allowed [256]bool

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func init() {
	allowed = initValidChars()
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
	key := strings.ToLower(string(parts[0]))

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimLeft(key, " ")

	if !isValid(key) {
		return 0, false, fmt.Errorf("invalid field name: '%s'", key)
	}

	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Add(key, value string) {
	key = strings.ToLower(key)
	if v, exists := h[key]; exists {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	if _, exists := h[key]; exists {
		return h[key]
	}
	return ""
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func initValidChars() [256]bool {
	for b := byte('a'); b <= byte('z'); b++ {
		allowed[b] = true
	}
	for b := byte('0'); b <= byte('9'); b++ {
		allowed[b] = true
	}
	specials := "!#$%&'*+-.^_`|~"
	for i := 0; i < len(specials); i++ {
		allowed[specials[i]] = true
	}
	return allowed
}

func isValid(s string) bool {
	return !(hasTrailingSpaces(s) || containsInvalidChars(s))
}

func hasTrailingSpaces(s string) bool {

	return len(s) > 0 && len(strings.TrimRight(s, " ")) < len(s)
}

func containsInvalidChars(s string) bool {
	for i := 0; i < len(s); i++ {
		if !allowed[s[i]] {
			return true
		}
	}
	return false
}
