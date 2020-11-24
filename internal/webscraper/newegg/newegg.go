package newegg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html"

	"github.com/sprucewillis/nvidia-finder/internal/util"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

func parseConfig() (NeweggConfig, error) {
	var config NeweggConfig
	rawConfig, err := ioutil.ReadFile("./src/github.com/sprucewillis/nvidia-finder/internal/webscraper/newegg_config.json")
	if err != nil {
		log.Println("error: unable to read Newegg card config file")
		return NeweggConfig{}, err
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Println("error: unable to parse Newegg config from JSON")
		return NeweggConfig{}, err
	}
	return config, nil
}

// CheckNewegg check newegg stock by individual card pages
func CheckNewegg(client *http.Client, config NeweggConfig, c chan inventory.Item) {
	items := config.items
	searchPageUrls := config.searchPages
	log.Printf("found %v cards to check at newegg", len(items))
	if config.scrapeItemPages {
		go checkIndividualItems(client, items, c)
	}
	if config.scrapeSearchPages {
		go checkSearchPages(client, searchPageUrls, c)
	}
}

func checkSearchPages(client *http.Client, searchPageUrls []string, c chan inventory.Item) {
	for {
		scrambledUrls := util.ShuffleString(searchPageUrls)
		for _, url := range scrambledUrls {
			itemsInStock, err := checkNeweggSearchPage(client, url)
			if err != nil {
				continue
			}
			for _, item := range itemsInStock {
				c <- item
			}
		}
	}
}

func checkNeweggSearchPage(client *http.Client, url string) ([]inventory.Item, error) {
	// TODO pull newegg into its own package and separate overall, invididual, and search page scraping

	return nil, nil
}

func checkIndividualItems(client *http.Client, items []inventory.Item, c chan inventory.Item) {
	numRetries := 1
	for {
		scrambledItems := util.ShuffleItems(items)
		for _, item := range scrambledItems {
			status, err := checkItemStatusWithRetries(client, item.URL, numRetries)
			if err != nil {
				log.Println("error: unable to parse data for newegg item", item.Name)
			}
			if status {
				logInStock(item)
				item.Site = "newegg"
				c <- item
			} else {
				log.Println(item.Name, "not in stock at Newegg")
			}
			util.RandomSleep(1, 4)
		}
	}
}

func logInStock(item inventory.Item) {
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	log.Println(string(colorGreen), "Found in stock", item.Name, "at Newegg, url:", item.URL, string(colorReset))
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
}
