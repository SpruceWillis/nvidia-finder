package scraper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html"

	"github.com/sprucewillis/nvidia-finder/internal/util"
)

var items []Item

func init() {
	rawConfig, err := ioutil.ReadFile("./src/github.com/sprucewillis/nvidia-finder/internal/webscraper/newegg_config.json")
	if err != nil {
		log.Println("error: unable to read Newegg card config file")
	}
	err = json.Unmarshal(rawConfig, &items)
	if err != nil {
		log.Println("error: unable to parse Newegg config from JSON")
	}
}

// CheckNewegg check newegg stock by individual card pages
func CheckNewegg(client *http.Client, c chan Item) {
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	numRetries := 1
	// TODO: randomize time
	log.Printf("found %v cards to check at newegg", len(items))
	for {
		foundMatch := false
		for _, card := range items {
			status, err := checkItemStatusWithRetries(client, card.URL, numRetries)
			if err != nil {
				log.Println("error: unable to parse data for newegg card", card.Name)
			}
			if status {
				log.Println(string(colorGreen), "Found in stock", card.Name, "at Newegg, url:", card.URL, string(colorReset))
				card.Site = "newegg"
				c <- card
				foundMatch = true
			} else {
				log.Println(card.Name, "not in stock at Newegg")
			}
			util.RandomSleep(5, 10)
		}
		if !foundMatch {
			log.Println("nothing in stock at Newegg")
		}
		util.RandomSleep(35, 60)
	}
}

// retry so we can reduce the number of false positives
func checkItemStatusWithRetries(client *http.Client, url string, numRetries int) (bool, error) {
	status, err := checkItemStatus(client, url)
	if err != nil {
		return false, err
	}
	if status {
		if numRetries == 0 {
			return status, nil
		} else {
			log.Println("card at", url, "possibly in stock at newegg, retrying to confirm")
			return checkItemStatusWithRetries(client, url, numRetries-1)
		}
	} else {
		return status, nil
	}
}

func checkItemStatus(client *http.Client, url string) (bool, error) {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return parseItemStatus(resp, url), nil
}

func parseItemStatus(resp *http.Response, url string) bool {
	log.Println("checking", url)
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("WARN: unable to parse HTML body of ", url)
		return false
	} else if resp.StatusCode != 200 {
		log.Printf("WARN: bad status code for %v: %v", url, resp.StatusCode)
	}
	bodyNode := doc.LastChild.LastChild
	var appNode *html.Node
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key == "id" && attr.Val == "app" {
				appNode = c
				break
			}
		}
	}
	if appNode == nil {
		log.Println("warn: unable to find app node")
		return false
	}
	var pageContentNode *html.Node
	for c := appNode.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key == "class" && attr.Val == "page-content" {
				pageContentNode = c
				break
			}
		}
	}
	if pageContentNode == nil {
		log.Println("warn: unable to find page content node")
		return false
	}
	productBuyBoxNode := pageContentNode.FirstChild.FirstChild.FirstChild.FirstChild.FirstChild
	var productBuyNode *html.Node
	for c := productBuyBoxNode.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key == "class" && attr.Val == "product-buy" {
				productBuyNode = c
				break
			}
		}
	}
	if productBuyNode == nil {
		log.Println("warn: unable to find product buy node")
		return false
	}
	navRow := productBuyNode.FirstChild
	for c := navRow.FirstChild; c != nil; c = c.NextSibling {
		for _, attr := range c.Attr {
			if attr.Key == "class" && attr.Val == "nav-col has-qty-box" {
				return true
			}
		}
	}
	return false

	// roughly the sequence on the page is:
	// div class=page-content
	// div class = page-section
	// page-section-inner
	// row is-product has-side-right has-side-items
	// row-side
	// product-buy-box
	// product-buy
	// nav-row
	// nav-col - if there's a quantity box, things are for sure in stock
	// otherwise we can recursively check for an add-to-cart button
}
