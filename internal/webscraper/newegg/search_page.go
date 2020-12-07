package newegg

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/sprucewillis/nvidia-finder/internal/util/htmlutils"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/headers"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
	"golang.org/x/net/html"
)

func checkNeweggSearchPage(client *http.Client, url string) ([]inventory.Item, error) {
	itemsInStock, err := getSearchPageItems(client, url)
	if err != nil {
		log.Println("error parsing newegg search page")
	}
	if len(itemsInStock) == 0 {
		log.Println("nothing in stock at", url)
	} else {
		itemNames := make([]string, 0)
		for _, item := range itemsInStock {
			itemNames = append(itemNames, item.Name)
		}
		log.Println("in-stock items found at", url, strings.Join(itemNames, "\n"))
	}
	return itemsInStock, err
}

// TODO consolidate these blocks into something more standardized across all clients and move to main webscraper package
func getSearchPageItems(client *http.Client, url string) ([]inventory.Item, error) {
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("unable to create HTTP request:", err)
		return nil, err
	}
	for k, v := range headers.MacHeaders {
		req.Header.Add(k, v)
	}
	// TODO figure out why this is bugged and sometimes has an unexpected EOF
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error fetching search page:", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println("error parsing response html:", err)
		return nil, err
	}
	// this has to be parsed before the defer kicks in, or we get unexpected EOF errors
	itemsInStock, err := parseSearchStatus(doc)
	if err != nil {
		log.Println("error parsing search results:", err)
		return nil, err
	}
	return itemsInStock, nil
}

func parseSearchStatus(root *html.Node) ([]inventory.Item, error) {
	rootHTMLNode := htmlutils.FindChildWithAttributeKey(root, "lang")
	if rootHTMLNode == nil {
		return nil, errors.New("unable to find root node in newegg search page")
	}
	bodyNode := htmlutils.FindChild(rootHTMLNode, func(node *html.Node) bool {
		return node.Data == "body"
	})
	if bodyNode == nil {
		return nil, errors.New("unable to find body node in newegg search page")
	}
	appNode := htmlutils.FindChildByAttribute(bodyNode, html.Attribute{Namespace: "", Key: "id", Val: "app"})
	if appNode == nil {
		return nil, errors.New("unable to find app node in newegg search page")
	}
	pageContentNode := htmlutils.FindChildByAttribute(appNode, html.Attribute{Namespace: "", Key: "class", Val: "page-content"})
	if pageContentNode == nil {
		return nil, errors.New("unable to find page contennt node in newegg search page")
	}
	pageSectionNode := htmlutils.FindChildByAttribute(pageContentNode, html.Attribute{Namespace: "", Key: "class", Val: "page-section"})
	if pageSectionNode == nil {
		return nil, errors.New("unable to find page section in newegg search page")
	}
	pageSectionInnerNode := htmlutils.FindChildByAttribute(pageSectionNode, html.Attribute{Namespace: "", Key: "class", Val: "page-section-inner"})
	if pageSectionInnerNode == nil {
		return nil, errors.New("unable to find page section inner in newegg search page")
	}
	rowLeftNode := htmlutils.FindChildByAttribute(pageSectionInnerNode, html.Attribute{Namespace: "", Key: "class", Val: "row has-side-left"})
	if rowLeftNode == nil {
		return nil, errors.New("unable to find node with class row has-side-left in newegg search poage")
	}
	rowBody := htmlutils.FindChildByAttribute(rowLeftNode, html.Attribute{Namespace: "", Key: "class", Val: "row-body"})
	if rowBody == nil {
		return nil, errors.New("unable to find node with class row-body in newegg search page")
	}
	rowBodyInner := htmlutils.FindChildByAttribute(rowBody, html.Attribute{Namespace: "", Key: "class", Val: "row-body-inner"})
	if rowBodyInner == nil {
		return nil, errors.New("unable to find node with class row-body-inner in newegg search page")
	}
	rowBodyBorder := htmlutils.FindChildWithClass(rowBodyInner, true, "row-body-border")
	if rowBodyBorder == nil {
		htmlutils.PrettyPrintHTMLNode(rowBodyInner, 2)
		return nil, errors.New("unable to find node with class row-body-border in newegg search page")
	}
	row := htmlutils.FindChildWithClass(rowBodyBorder, false, "row")
	if row == nil {
		return nil, errors.New("unable to find node with class row in newegg search page")
	}
	rowBody = htmlutils.FindChildWithClass(row, true, "row-body")
	if rowBody == nil {
		return nil, errors.New("unable to find node with class row-body in newegg search page")
	}
	rowBodyInner = htmlutils.FindChildWithClass(rowBody, true, "row-body-inner")
	if rowBodyInner == nil {
		return nil, errors.New("unable to find inner node with class row-body-inner in newegg search page")
	}
	listWrap := htmlutils.FindChildByAttribute(rowBodyInner, html.Attribute{Namespace: "", Key: "class", Val: "list-wrap"})
	if listWrap == nil {
		return nil, errors.New("unable to find node with class list-wrap in newegg search page")
	}
	// TODO iterate over all of these and find their children
	itemsInStock, err := findInStockItems(listWrap)
	if err != nil {
		return nil, err
	}
	return itemsInStock, nil
}

func findInStockItems(listWrap *html.Node) ([]inventory.Item, error) {
	itemsInStock := make([]inventory.Item, 0)
	itemCellWrappers := htmlutils.FindChildrenByClass(listWrap, false, "item-cells-wrap")
	for _, itemCellWrapperNode := range itemCellWrappers {
		itemCells := htmlutils.FindChildrenByClass(itemCellWrapperNode, true, "item-cell")
		for _, itemCell := range itemCells {
			itemContainer := htmlutils.FindChildWithClass(itemCell, true, "item-container")
			itemPtr, inStock := parseItemInformation(itemContainer)
			if inStock {
				itemsInStock = append(itemsInStock, *itemPtr)
			}
		}
	}
	return itemsInStock, nil
}

func parseItemInformation(itemContainer *html.Node) (*inventory.Item, bool) {
	if itemContainer == nil {
		// log.Println("unable to find item container")
		return nil, false
	}
	itemInfo := htmlutils.FindChildWithClass(itemContainer, true, "item-info")
	if itemInfo == nil {
		// log.Println("unable to find item info")
		return nil, false
	}
	// this is where the OUT OF STOCK part goes
	itemPromo := htmlutils.FindChildWithClass(itemInfo, true, "item-promo")
	if itemPromo != nil && itemPromo.FirstChild.NextSibling.Data == "OUT OF STOCK" {
		// log.Println("promo found, item not in stock")
		return nil, false
	}
	itemTitle := htmlutils.FindChildWithClass(itemInfo, true, "item-title")
	if itemTitle == nil {
		// log.Println("unable to find item title")
		return nil, false
	}
	itemName := itemTitle.FirstChild.Data
	itemURL, hasURL := htmlutils.GetAttribute(itemTitle, "href")
	if !hasURL {
		// log.Println(fmt.Sprintf("WARN item %v is in stock, but appears to be missing a URL."))
		return nil, false
	}
	item := inventory.Item{Name: itemName, URL: itemURL, Site: "Newegg", PriceLimit: 0, Sku: ""}
	return &item, true
}
