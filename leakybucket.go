package main

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	mu sync.Mutex

	Capacity  int
	Tokens    int
	LastCheck time.Time
}

func (lb *LeakyBucket) Add() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
}

func (lb *LeakyBucket) Allow() bool {
	return true
}
