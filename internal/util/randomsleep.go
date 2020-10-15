package util

import (
  "math/rand"
  "time"
)

func RandomSleep(lowerBound, upperBound int){
  diff := upperBound - lowerBound
  sleepTime := rand.Intn(diff + 1) + lowerBound
  time.Sleep(time.Duration(sleepTime) * time.Second)
}
