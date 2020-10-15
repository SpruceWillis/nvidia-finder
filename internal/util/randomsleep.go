package util

import (
	"math/rand"
	"time"
)

// RandomSleep for a randomly determined number of seconds between the lower and upper bounds
func RandomSleep(lowerBound, upperBound int) {
	diff := upperBound - lowerBound
	sleepTime := rand.Intn(diff+1) + lowerBound
	time.Sleep(time.Duration(sleepTime) * time.Second)
}
