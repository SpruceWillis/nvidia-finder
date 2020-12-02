package inventory

// Item for use in parsing JSON scraping config and passing on email
type Item struct {
	URL        string
	Name       string
	Site       string
	PriceLimit float64
	Sku        string // currently only used in Best Buy
}

// IsBelowPriceLimit (price float64) determine whether the available price is ok or not
func (i Item) IsBelowPriceLimit(price float64) bool {
	return i.PriceLimit <= 0 || price <= i.PriceLimit
}
