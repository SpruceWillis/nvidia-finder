package util

import (
	"math/rand"

	inventory "github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

// ShuffleString []string implements the Fisher-Yates shuffle to do an in-place shuffle
func ShuffleString(arr []string) []string {
	sliceLength := len(arr)
	for i := sliceLength - 1; i >= 1; i-- {
		random := rand.Intn(i)
		arr[i], arr[random] = arr[random], arr[i]
	}
	return arr
}

// ShuffleItems []scraper.Item implement the Fisher-Yates shuffle to do an in-place shuffle
func ShuffleItems(arr []inventory.Item) []inventory.Item {
	sliceLength := len(arr)
	for i := sliceLength - 1; i >= 1; i-- {
		random := rand.Intn(i)
		arr[i], arr[random] = arr[random], arr[i]
	}
	return arr
}
