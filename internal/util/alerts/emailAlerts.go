package alerts

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sprucewillis/nvidia-finder/internal/email"
	"github.com/sprucewillis/nvidia-finder/internal/email/auth"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

// SetupAlerts read email configuration and send email alerts to the folks it is configured for
func SetupEmailAlerts(c chan inventory.Item) error {
	// TODO do the auth stuff for the email, exit
	creds := auth.GetSmtpCreds()
	plainAuth := creds.GetPlainAuth()
	from := creds.Username
	to := GetEmailRecipients()
	log.Println("email recipients:", strings.Join(to, ","))
	for itemInStock := range c {
		if len(to) > 0 {
			err := email.SendInStockEmail(itemInStock.Site, itemInStock.URL, creds.GetURL(), from, to, plainAuth)
			if err != nil {
				fmt.Println("WARN: unable to send email", err)
			}
		}
	}
	return nil
}

// GetEmailRecipients read the recipient list from a file name
func GetEmailRecipients() []string {
	var result []string
	fileName := "recipients.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("WARN:", err, "Emails will not be sent")
		return []string{}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("WARN: error reading emails:", err, "Emails will not be sent")
		return []string{}
	}
	return result
}
