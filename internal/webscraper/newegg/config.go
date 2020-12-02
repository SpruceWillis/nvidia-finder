package newegg

import "github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"

// Config specify search pages and specific item pages to scrape
type Config struct {
	Items             []inventory.Item
	SearchPages       []string
	ScrapeSearchPages bool
	ScrapeItemPages   bool
}
