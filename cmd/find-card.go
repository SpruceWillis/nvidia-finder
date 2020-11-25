package main

import (
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/util/alerts"
	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

func main() {
	client := &http.Client{}
	emailChannel := make(chan inventory.Item)
	go alerts.SetupEmailAlerts(emailChannel)
	go scraper.CheckBestBuy(client, false, emailChannel)
	go scraper.CheckNewegg(client, emailChannel)
	for {
		// keep program alive
	}
}
