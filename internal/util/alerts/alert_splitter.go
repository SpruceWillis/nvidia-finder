package alerts

import (
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

// SetUpAlertChannelSplitter(channels []chan inventory.Item) send all items to listening alerting channels
func SetUpAlertChannel(channels []chan inventory.Item) chan inventory.Item {
	itemChannelSplitter := make(chan inventory.Item)
	go func() {
		for item := range itemChannelSplitter {
			for _, c := range channels {
				c <- item
			}
		}
	}()
	return itemChannelSplitter
}
