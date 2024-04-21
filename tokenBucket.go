package tokenbucket

import (
	"net/http"
	"sync"
	"time"
)

type TOKENS int

const (
	MAX_TOKENS TOKENS = 1                // max 1 request per rate seconds (leetcode ftw!)
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
	} else {
		return false
	}
}

func (b *Bucket) refill() {
	refill := TOKENS(time.Since(b.lastRefill) / b.rate)

	if refill > 0 {
		b.tokens = min(b.capacity, refill+b.tokens)
		b.lastRefill = time.Now()
	}
}

// TokenBuckets is in memory representation of USER buckets
// TokenBuckets is a map of Bucket with userId as key
// userId is a unique key assigned to every user when sign in to the application
// userId is sent with the request header with every request
type TokenBuckets struct {
	mu      sync.Mutex
	buckets map[string]*Bucket
}

func NewTokenBuckets() *TokenBuckets {
	return &TokenBuckets{
		buckets: make(map[string]*Bucket),
	}
}

func (db *TokenBuckets) getBucket(userId string) *Bucket {
	db.mu.Lock()
	defer db.mu.Unlock()

	bucket, exists := db.buckets[userId]

	if !exists {
		bucket = newBucket(MAX_TOKENS)
		db.buckets[userId] = bucket
	}

	return bucket
}

// Remove the users from the memory who haven't made any request in last 1 day to keep memory in check
func (db *TokenBuckets) ClearOldBuckets() {
	db.mu.Lock()
	defer db.mu.Unlock()

	for userId, bucket := range db.buckets {
		if time.Since(bucket.lastRefill) > 24*time.Hour {
			delete(db.buckets, userId)
		}
	}
}

func RateLimiter(next http.Handler, users *TokenBuckets) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("userId")
		if userId == "" {
			http.Error(w, "UserID missing from the request", http.StatusBadRequest)
			return
		}

		bucket := users.getBucket(userId)

		if !bucket.Allow() {
			http.Error(w, "Too many requests, Try again after some time", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
