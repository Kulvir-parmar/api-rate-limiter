package main

import (
	"errors"
	"time"
)

const (
	MAX_TOKENS    = 10
	REFILL_AMOUNT = 2
	REFILL_TIME   = int(time.Minute)
)

type Bucket struct {
	CurrentTokens int
	LastAccessed  time.Time
}

func (b *Bucket) NewBucket(userId string) *Bucket {
	return &Bucket{
		CurrentTokens: MAX_TOKENS,
		LastAccessed:  time.Now(),
	}
}

// NOTE: user redis to store this DB
var DB map[string]*Bucket

func RefillBucket(userId string) {
	lastAccessTime := DB[userId].LastAccessed
	elapsedTime := int(time.Now().Sub(lastAccessTime)) / REFILL_TIME
	refillCount := elapsedTime * REFILL_AMOUNT

	presentTokens := DB[userId].CurrentTokens
	DB[userId].CurrentTokens = min(presentTokens+refillCount, MAX_TOKENS)
}

func isTokensAvailable(userId string) bool {
	return DB[userId].CurrentTokens > 0
}

func consumeToken(userId string) error {
	if isTokensAvailable(userId) {
		DB[userId].CurrentTokens--
		DB[userId].LastAccessed = time.Now()
		return nil
	} else {
		return errors.New("API Limit reached")
	}
}

func main() {
}
