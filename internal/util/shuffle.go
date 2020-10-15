package util

import "math/rand"

// ShuffleString []string implements the Fisher-Yates shuffle to do an in-place shuffle
func ShuffleString(arr []string) []string {
	sliceLength := len(arr)
	for i := sliceLength - 1; i >= 1; i-- {
		random := rand.Intn(i)
		arr[i], arr[random] = arr[random], arr[i]
	}
	return arr
}
