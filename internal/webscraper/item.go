package scraper

// Item for use in parsing JSON scraping config and passing on email
type Item struct {
	URL        string
	Name       string
	Site       string
	PriceLimit float64
}
