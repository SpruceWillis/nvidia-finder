package main

import (
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/util/alerts"
	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

func main() {
	client := &http.Client{}
	alertChannel := make(chan inventory.Item)
	go alerts.SetupEmailAlerts(alertChannel)
	go alerts.SetupAudioAlerts(alertChannel)
	go scraper.CheckBestBuy(client, false, alertChannel)
	go scraper.CheckNewegg(client, alertChannel)
	for {
		// keep program alive
	}
}
