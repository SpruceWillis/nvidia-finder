package newegg

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sprucewillis/nvidia-finder/internal/util/htmlutils"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
	"golang.org/x/net/html"
)

func checkNeweggSearchPage(client *http.Client, url string) ([]inventory.Item, error) {
	// TODO pull newegg into its own package and separate overall, invididual, and search page scraping
	rootNode, err := getSearchPageHTML(client, url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	itemsInStock, err := parseSearchStatus(rootNode)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return itemsInStock, nil
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

func parseSearchStatus(root *html.Node) ([]inventory.Item, error) {
	itemsInStock := []inventory.Item{}
	body := htmlutils.FindChild(root, func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == "lang" {
				return true
			}
		}
		return false
	})
	if body == nil {
		fmt.Println("nobody found")
		return nil, nil
	}
	htmlutils.PrettyPrintHTMLNode(body, 2)
	return itemsInStock, nil
}
