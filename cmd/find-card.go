package main

import (
	"net/http"

	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
)

func main() {
	client := &http.Client{}
	// TODO send these results to a channel for message sending
	go scraper.CheckBestBuy(client, false)
	go scraper.CheckNewegg(client)
	for {
		// keep program alive
	}
}
