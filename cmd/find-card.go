package main

import (
	"fmt"
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/email/auth"
)

func main() {
	client := &http.Client{}
	// async check all the things until they don't work no more
	creds := auth.GetSmtpCreds()
	plainAuth := creds.GetPlainAuth()
	fmt.Println(plainAuth)
	// go scraper.CheckBestBuy(client, false)
	// go scraper.CheckNewegg(client)
	_ = client
	// gmailClient, err := email.GetGmailClient()
	// if err != nil {
	// log.Println("found client %v", gmailClient)
	// }
	for {
		// keep program alive
	}
}
