package email

import (
	"fmt"

	"github.com/sprucewillis/nvidia-finder/internal/email/auth"
)

func inStockMessage(site, url string) string {
	return fmt.Sprint("product in stock at", site, "url:", url)
}

func SendMail(site, url string, creds *auth.Credentials) err {
	message := inStockMessage(site, url)

}
