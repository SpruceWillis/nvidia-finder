package email

import (
	"fmt"
	"log"
	"net/smtp"
)

func inStockMessage(from, site, url, to string) []byte {
	message := fmt.Sprintf("To: %v\r\n", to) +
		"Subject: GPU in stock\r\n" +
		"\r\n" +
		fmt.Sprintf("GPU in stock at %v, url: %v\r\n", site, url)
	return []byte(message)
}

// Send wrap smtp SendMail
func Send(site, url, from string, to []string, auth *smtp.Auth) error {
	for _, recipient := range to {
		message := inStockMessage(from, site, url, recipient)
		err := smtp.SendMail(url, *auth, from, to, message)
		if err != nil {
			log.Println("unable to send email to", recipient, err)
			return err
		}
	}
	return nil
}
