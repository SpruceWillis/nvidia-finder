package main

import (
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/util/alerts"
	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/newegg"
)

func main() {
	client := &http.Client{}
	// email
	emailAlertChannel := make(chan inventory.Item)
	go alerts.SetupEmailAlerts(emailAlertChannel)

	// audio
	audioAlertChannel := make(chan inventory.Item)
	go alerts.SetupAudioAlerts(audioAlertChannel)

	// overall alerts
	alertChannel := alerts.SetUpAlertChannel([]chan inventory.Item{emailAlertChannel, audioAlertChannel})
	neweggConfig, err := newegg.ReadNeweggConfig()
	if err == nil {
		go newegg.CheckNewegg(client, neweggConfig, alertChannel)
	}
	go scraper.CheckBestBuy(client, false, alertChannel)
	select {
	// keep program alive
	}
}
