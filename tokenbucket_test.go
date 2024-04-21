package tokenbucket

import (
	"net/http"
	"testing"
)

var users = NewTokenBuckets()

func TestRateLimiter(t *testing.T) {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
}
