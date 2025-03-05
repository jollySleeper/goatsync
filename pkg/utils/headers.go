package utils

import (
	"net/http"
	"strings"
)

// parse all headers from a request
func ParseHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for k, v := range r.Header {
		headers[k] = strings.Join(v, ",")
	}
	return headers
}

// get a specific header from a request
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}
