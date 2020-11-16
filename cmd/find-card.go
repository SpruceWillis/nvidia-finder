package main

import (
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/email"
	scraper "github.com/sprucewillis/nvidia-finder/internal/webscraper"
)

func main() {
	client := &http.Client{}
	emailChannel := make(chan scraper.Item)
	go email.SetupAlerts(emailChannel)
	go scraper.CheckBestBuy(client, false, emailChannel)
	go scraper.CheckNewegg(client, emailChannel)
	for {
		// keep program alive
	}
}
