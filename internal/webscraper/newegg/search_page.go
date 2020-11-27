package newegg

import (
	"log"
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
	"golang.org/x/net/html"
)

func checkNeweggSearchPage(client *http.Client, url string) ([]inventory.Item, error) {
	// TODO pull newegg into its own package and separate overall, invididual, and search page scraping

	return nil, nil
}

// TODO consolidate these blocks into something more standardized across all clients and move to main webscraper package
func getSearchPageHTML(client *http.Client, url string) (*html.Node, error) {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return doc, nil
}

func parseSearchStatus(root *html.Node) []inventory.Item {

	return nil
}
