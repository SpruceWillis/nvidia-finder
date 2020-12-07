package newegg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sprucewillis/nvidia-finder/internal/util"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

// ReadNeweggConfig read newegg config from file
func ReadNeweggConfig() (Config, error) {
	var config Config
	filePath := "./src/github.com/sprucewillis/nvidia-finder/internal/webscraper/newegg/config.json"
	rawConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("error: unable to read Newegg card config file")
		return Config{}, err
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Println("error: unable to parse Newegg config from JSON")
		return Config{}, err
	}
	return config, nil
}

// CheckNewegg check newegg stock by individual card pages
func CheckNewegg(client *http.Client, config Config, c chan inventory.Item) {
	if config.ScrapeItemPages {
		items := config.Items
		log.Printf("found %v items to check at newegg", len(items))
		go checkIndividualItems(client, items, c)
	} else {
		log.Println("skipping newegg item check")
	}
	if config.ScrapeSearchPages {
		searchPageUrls := config.SearchPages
		log.Printf("found %v searches to check at newegg", len(searchPageUrls))
		go checkSearchPages(client, searchPageUrls, c)
	} else {
		log.Println("skipping newegg search check")
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
			util.RandomSleep(3, 5)
		}
		util.RandomSleep(25, 30)
	}
}

func checkIndividualItems(client *http.Client, items []inventory.Item, c chan inventory.Item) {
	alertCooldown, _ := time.ParseDuration("3min")
	numRetries := 0
	itemFoundCache := make(map[string]bool)
	for {
		scrambledItems := util.ShuffleItems(items)
		for _, item := range scrambledItems {
			status, err := checkIndividualItemWithRetries(client, item, numRetries)
			if err != nil {
				log.Println("error: unable to parse data for newegg item", item.Name)
			}
			if status {
				logInStock(item)
				_, foundInCache := itemFoundCache[item.URL]
				if !foundInCache {
					item.Site = "newegg"
					c <- item
					itemFoundCache[item.URL] = true
					time.AfterFunc(alertCooldown, func() {
						delete(itemFoundCache, item.URL)
					})
				}
			} else {
				log.Println(item.Name, "not in stock at Newegg")
			}
			util.RandomSleep(8, 10)
		}
	}
}

func logInStock(item inventory.Item) {
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	log.Println(string(colorGreen), "Found in stock", item.Name, "at Newegg, url:", item.URL, string(colorReset))
}
