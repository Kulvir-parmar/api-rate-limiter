package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	rate := time.Second * 10
	time.Sleep(time.Second * 11)

	refill := float64(time.Since(now) / rate)
	fmt.Println(refill)
}
