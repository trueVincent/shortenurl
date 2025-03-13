package main

import (
	"math/rand"
	"time"
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}