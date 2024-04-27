package tokenbucket

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	users := NewTokenBuckets()
	handler := RateLimiter(next, users)

	for i := 0; i < 4; i++ {
		handler.ServeHTTP(rr, req)

		if status := rr.Code; i < 2 && status == http.StatusOK {
			t.Logf("Got HTTP Status: %v", status)
		} else if i >= 2 && status != http.StatusTooManyRequests {
			t.Errorf("HTTP request should return %v but got %v insted", http.StatusTooManyRequests, status)
		} else {
			t.Logf("Got HTTP Status: %v", status)
		}
	}
}
