package email

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func inStockMessage(from, site, url string, to string) []byte {
	message := fmt.Sprintf("To: %v\r\n", to) +
		"Subject: Item in stock\r\n" +
		"\r\n" +
		fmt.Sprintf("Item in stock at %v, url: %v\r\n", site, url)
	return []byte(message)
}

// SendInStockEmail email the selected folks about things being in stock
func SendInStockEmail(site, productURL, smtpURL, from string, to []string, auth *smtp.Auth) error {
	toCSV := strings.Join(to, ",")
	message := inStockMessage(from, site, productURL, toCSV)
	err := smtp.SendMail(smtpURL, *auth, from, to, message)
	if err != nil {
		log.Println("unable to send email", err)
		return err
	}
	log.Println("instock email(s) sent to indicated recipients")
	return nil
}
