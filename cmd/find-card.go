package main

import (
  "net/http"
  // "log"
  "github.com/sprucewillis/nvidia-finder/internal/webscraper"
  // "github.com/sprucewillis/nvidia-finder/internal/email"
)

func main() {
  client := &http.Client{}
  // async check all the things until they don't work no more
  go scraper.CheckBestBuy(client, false)
  go scraper.CheckNewegg(client)
  _ = client
  // gmailClient, err := email.GetGmailClient()
  // if err != nil {
    // log.Println("found client %v", gmailClient)
  // }
  for {
    // keep program alive
  }
}
