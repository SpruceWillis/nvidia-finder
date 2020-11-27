package newegg

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// retry so we can reduce the number of false positives
func checkIndividualItemWithRetries(client *http.Client, url string, numRetries int) (bool, error) {
	status, err := checkIndividualItemStatus(client, url)
	if err != nil {
		return false, err
	}
	if status {
		if numRetries == 0 {
			return status, nil
		} else {
			log.Println("card at", url, "possibly in stock at newegg, retrying to confirm")
			return checkIndividualItemWithRetries(client, url, numRetries-1)
		}
	} else {
		return status, nil
	}
}

func checkIndividualItemStatus(client *http.Client, url string) (bool, error) {
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
