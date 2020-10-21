package scraper

// Card for use in parsing JSON scraping config and passing on email
type Card struct {
	URL  string
	Name string
	Site string
}
