package main

import (
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/util/alerts"
	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

func main() {
	client := &http.Client{}
	emailAlertChannel := make(chan inventory.Item)
	audioAlertChannel := make(chan inventory.Item)
	go alerts.SetupEmailAlerts(emailAlertChannel)
	go alerts.SetupAudioAlerts(audioAlertChannel)
	alertChannel := alerts.SetUpAlertChannel([]chan inventory.Item{emailAlertChannel, audioAlertChannel})
	go scraper.CheckBestBuy(client, false, alertChannel)
	go scraper.CheckNewegg(client, alertChannel)
	select {
	// keep program alive
	}
}
