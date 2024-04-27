package tokenbucket

import (
	"net/http"
	"sync"
	"time"
)

type TOKENS int

const (
	MAX_TOKENS TOKENS = 2                // max 2 request per rate seconds (leetcode ftw!)
	RATE              = time.Second * 10 // 10 seconds is default rate
)

type Bucket struct {
	mu         sync.Mutex
	capacity   TOKENS // max tokens allowed
	tokens     TOKENS // tokens present in the bucket
	lastRefill time.Time
	rate       time.Duration // rate at which tokens are refilled
}

func newBucket(capacity TOKENS) *Bucket {
	return &Bucket{
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
		rate:       RATE,
	}
}

// Allow fn. return a Boolean value representing wheather user is allowed to make the request or not
func (b *Bucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.refill()

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

func (b *Bucket) refill() {
	refill := TOKENS(time.Since(b.lastRefill) / b.rate)

	if refill > 0 {
		b.tokens = min(b.capacity, refill+b.tokens)
		b.lastRefill = time.Now()
	}
}

// TokenBuckets is in memory representation of USER buckets
// TokenBuckets is a map of Bucket with IP Address as the key
type TokenBuckets struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func NewTokenBuckets() *TokenBuckets {
	return &TokenBuckets{
		buckets: make(map[string]*Bucket),
	}
}

func (db *TokenBuckets) getBucket(ip string) *Bucket {
	db.mu.Lock()
	defer db.mu.Unlock()

	bucket, exists := db.buckets[ip]

	if !exists {
		bucket = newBucket(MAX_TOKENS)
		db.buckets[ip] = bucket
	}

	return bucket
}

// Remove the stale request from the memory
func (db *TokenBuckets) ClearOldBuckets() {
	db.mu.Lock()
	defer db.mu.Unlock()

	for ip, bucket := range db.buckets {
		if time.Since(bucket.lastRefill) > time.Hour {
			delete(db.buckets, ip)
		}
	}
}

func RateLimiter(next http.Handler, users *TokenBuckets) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		bucket := users.getBucket(ip)

		if !bucket.Allow() {
			http.Error(w, "Too many requests, Try again after some time", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
