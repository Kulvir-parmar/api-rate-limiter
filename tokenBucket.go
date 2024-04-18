package main

import (
	"errors"
	"time"
)

const (
	MAX_TOKENS  = 10
	REFILL_RATE = 2
)

type Bucket struct {
	Tokens     int
	LastRefill time.Time
}

func (b *Bucket) NewBucket(userId string) *Bucket {
	return &Bucket{
		Tokens:     MAX_TOKENS,
		LastRefill: time.Now(),
	}
}

func RefillBucket(userId string) {
	lastRefill := DB[userId].LastRefill
	elapsed := time.Now().Sub(lastRefill).Minutes()

	tokensToAdd := int(elapsed / REFILL_RATE)

	if tokensToAdd > 0 {
		presentTokens := DB[userId].Tokens
		DB[userId].Tokens = min(presentTokens+tokensToAdd, MAX_TOKENS)

		DB[userId].LastRefill = time.Now()
	}
}

func isTokensAvailable(userId string) bool {
	return DB[userId].Tokens > 0
}

func consumeToken(userId string) error {
	if isTokensAvailable(userId) {
		DB[userId].Tokens--
		return nil
	} else {
		return errors.New("API Limit reached")
	}
}
