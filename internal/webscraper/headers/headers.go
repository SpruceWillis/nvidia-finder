package headers

// MacHeaders headers that make it look like the request is from Firefox on a Mac
var MacHeaders = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0",
	"Accept":          "*/*",
	"Accept-Language": "en-Us,en;q=0.5",
	"DNT":             "1",
}
