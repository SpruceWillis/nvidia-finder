package newegg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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
		}
	}
}

func checkIndividualItems(client *http.Client, items []inventory.Item, c chan inventory.Item) {
	numRetries := 1
	for {
		scrambledItems := util.ShuffleItems(items)
		for _, item := range scrambledItems {
			status, err := checkIndividualItemWithRetries(client, item.URL, numRetries)
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
