package newegg

import "github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"

// NeweggConfig specify search pages and specific item pages to scrape
type NeweggConfig struct {
	items             []inventory.Item
	searchPages       []string
	scrapeSearchPages bool
	scrapeItemPages   bool
}
