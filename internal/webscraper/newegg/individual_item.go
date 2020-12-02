package newegg

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sprucewillis/nvidia-finder/internal/util/htmlutils"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
	"golang.org/x/net/html"
)

const class string = "class"

// retry so we can reduce the number of false positives
func checkIndividualItemWithRetries(client *http.Client, item inventory.Item, numRetries int) (bool, error) {
	url := item.URL
	status, err := checkIndividualItemStatus(client, item)
	if err != nil {
		return false, err
	}
	if status {
		if numRetries == 0 {
			return status, nil
		} else {
			log.Println("card at", url, "possibly in stock at newegg, retrying to confirm")
			return checkIndividualItemWithRetries(client, item, numRetries-1)
		}
	} else {
		return status, nil
	}
}

func checkIndividualItemStatus(client *http.Client, item inventory.Item) (bool, error) {
	method := "GET"
	req, err := http.NewRequest(method, item.URL, nil)
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
	return parseItemStatus(resp, item), nil
}

func parseItemStatus(resp *http.Response, item inventory.Item) bool {
	url := item.URL
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
			if attr.Key == class && attr.Val == "page-content" {
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
			if attr.Key == class && attr.Val == "product-buy" {
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
				price, err := getItemPrice(productBuyBoxNode)
				if err != nil {
					log.Println("unable to determine price of", item.Name, err)
					return true // it's still in stock right?
				}
				return item.IsBelowPriceLimit(price)
			}
		}
	}
	return false
}

func getItemPrice(productBuyBoxNode *html.Node) (float64, error) {
	productPane := htmlutils.FindChild(productBuyBoxNode, func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "product-pane" {
				return true
			}
		}
		return false
	})
	if productPane == nil {
		return 0, errors.New("cannot find product pane html")
	}
	productPrice := htmlutils.FindChild(productPane, func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == class && attr.Val == "product-price" {
				return true
			}
		}
		return false
	})
	if productPrice == nil {
		return 0, errors.New("cannot find product price html")
	}
	priceListElement := htmlutils.FindChild(productPrice, func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == class && attr.Val == "price" {
				return true
			}
		}
		return false
	})
	if priceListElement == nil {
		return 0, errors.New("unable to find current price element")
	}
	currentPrice := htmlutils.FindChild(priceListElement, func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if attr.Key == class && attr.Val == "price-current" {
				return true
			}
		}
		return false
	})
	dollarAmountHeader := htmlutils.FindChild(currentPrice, func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == "strong"
	})
	if dollarAmountHeader == nil {
		return 0, errors.New("unable to find dollar amount")
	}
	dollarAmount := dollarAmountHeader.FirstChild.Data
	centAmountHeader := dollarAmountHeader.NextSibling
	if centAmountHeader == nil {
		return 0, errors.New("unable to find cents")
	}
	centAmount := centAmountHeader.FirstChild.Data
	totalAmount := fmt.Sprintf("%v%v", dollarAmount, centAmount)
	return strconv.ParseFloat(totalAmount, 64)
}
