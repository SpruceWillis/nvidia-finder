package scraper

// TODO reconfigure so that this is based on SKUs only, search is blocked LOL
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/sprucewillis/nvidia-finder/internal/util"
)

var bestBuySkus []string

// CheckBestBuy(*http.Client, bool) check best buy client
func CheckBestBuy(client *http.Client, findSkusFromWeb bool) {

	for {
		// checkInterval := randomNumber * time.Second
		bestBuyStatuses := getBestBuyStatus(client, findSkusFromWeb)
		foundMatch := false
		for _, v := range bestBuyStatuses {
			if v {
				foundMatch = true
			}
		}
		if !foundMatch {
			log.Println("nothing in stock at Best Buy")
		}
		// try not to get caught by the bot police
		util.RandomSleep(15, 20)
	}
}

func getBestBuyStatus(client *http.Client, findSkusFromWeb bool) map[string]bool {
	var skus []string
	if bestBuySkus != nil {
		skus = util.ShuffleString(bestBuySkus)
	} else if findSkusFromWeb {
		skus = getSkusFromWeb(client)
	} else {
		skus = getSkusFromFile()
	}
	fmt.Println("best buy skus:", skus)
	bestBuySkus = skus
	return getSkuStatuses(skus, client)
}

func getSkusFromFile() []string {
	var skus []string
	rawConfig, err := ioutil.ReadFile("./src/github.com/sprucewillis/nvidia-finder/internal/webscraper/bestbuy_config.json")
	if err != nil {
		log.Println("error: unable to read Best Buy card config file", err)
		return nil
	}
	err = json.Unmarshal(rawConfig, &skus)
	if err != nil {
		log.Println("error: unable to parse Best Buy config from file", err)
		return nil
	}
	return skus
}

// curl 'https://www.bestbuy.com/api/tcfb/model.json?paths=%5B%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%5B%22alternativeDataSetPresent%22%2C%22positionToSlotCarousel%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22redirect%22%2C%7B%22from%22%3A0%2C%22to%22%3A5%7D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22request_info%22%2C%22query%22%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sponsoredProducts%22%2C%5B%22onLoadBeacon%22%2C%22uuid%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22suggestQueryInfo%22%2C%22correctedQuery%22%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22documents%22%2C%7B%22from%22%3A0%2C%22to%22%3A24%7D%2C%5B%22id%22%2C%22onClickBeacon%22%2C%22onViewBeacon%22%2C%22sponsoredproduct%22%2C%22type%22%2C%22variationid%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sortOptions%22%2C%7B%22from%22%3A0%2C%22to%22%3A25%7D%2C%5B%22displayName%22%2C%22value%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sponsoredProducts%22%2C%22appliedAds%22%2C%22all%22%5D%5D&method=get' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' --compressed -H 'Referer: https://www.bestbuy.com/site/computer-cards-components/video-graphics-cards/abcat0507002.c?id=abcat0507002&qp=gpusv_facet%3DGraphics%20Processing%20Unit%20(GPU)~NVIDIA%20GeForce%20RTX%203080' -H 'DNT: 1' -H 'Connection: keep-alive' -H 'Pragma: no-cache' -H 'Cache-Control: no-cache' -H 'TE: Trailers'

// return a slice containing the SKUs that are in stock for Best Buy
func getSkusFromWeb(client *http.Client) []string {
	url := "https://www.bestbuy.com/api/tcfb/model.json?paths=%5B%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%5B%22alternativeDataSetPresent%22%2C%22positionToSlotCarousel%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22redirect%22%2C%7B%22from%22%3A0%2C%22to%22%3A5%7D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22request_info%22%2C%22query%22%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sponsoredProducts%22%2C%5B%22onLoadBeacon%22%2C%22uuid%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22suggestQueryInfo%22%2C%22correctedQuery%22%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22documents%22%2C%7B%22from%22%3A0%2C%22to%22%3A24%7D%2C%5B%22id%22%2C%22onClickBeacon%22%2C%22onViewBeacon%22%2C%22sponsoredproduct%22%2C%22type%22%2C%22variationid%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sortOptions%22%2C%7B%22from%22%3A0%2C%22to%22%3A25%7D%2C%5B%22displayName%22%2C%22value%22%5D%5D%2C%5B%22shop%22%2C%22ssic%22%2C%22optional%22%2C%22requestTypes%22%2C%22search%22%2C%22queries%22%2C%22categoryid%2524abcat0507002%22%2C%22numberOfRows%22%2C24%2C%22pageNumber%22%2C1%2C%22facets%22%2C%22gpusv_facet%3DNVIDIA%2520GeForce%2520RTX%25203080%22%2C%22options%22%2C%22autoFacet%3Dtrue%26availableStoresList%3D102%26browsedCategory%3Dabcat0507002%26enableSponsoredProduct%3Dfalse%26fieldList%3Dskuid%2Cid%2Ctype%26inventoryType%3Dstorepickup%26preferredStore%3D102%26rdp%3Dl%26searchTest2%3Dbestselling_bigquery%26searchTest3%3Dsearchmetric_bigquery%26sort%3DBest-Selling%22%2C%22sponsoredProducts%22%2C%22appliedAds%22%2C%22all%22%5D%5D&method=get"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	for k, v := range bestBuyHeaders() {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	documentSection, err := findDocumentSection(body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return parseDocumentsForSkus(documentSection)
}

func bestBuyHeaders() map[string]string {
	return map[string]string{
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0",
		"Accept":          "*/*",
		"Accept-Language": "en-Us,en;q=0.5",
		"DNT":             "1",
	}
}

/**
  recursively parse the json body
*/
func findDocumentSection(body []byte) (map[string]interface{}, error) {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		log.Fatal(err)
	}
	for k, v := range jsonBody {
		switch vv := v.(type) {
		case map[string]interface{}:
			if k == "documents" {
				return vv, nil
			} else {
				recurseResult, err := findDocumentSectionMap(vv)
				if err == nil {
					return recurseResult, nil
				}
			}
		case []interface{}:
			// TODO iterate over the map elements and see if they contain documents, otherwise carry on
			// for now it's ok to skip as we know it's in a map for sure as of this writing
		}
	}
	return nil, errors.New("unable to find document section")
}

// pretty sloppy, bue hopefully this works
func findDocumentSectionMap(section map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range section {
		switch vv := v.(type) {
		case map[string]interface{}:
			if k == "documents" {
				return vv, nil
			} else {
				recurseResult, err := findDocumentSectionMap(vv)
				if err == nil {
					return recurseResult, nil
				}
			}
		case []interface{}:
			// TODO iterate over the entries and see if the documents section is there
		}
	}
	return nil, errors.New("unable to find document section")
}

func parseDocumentsForSkus(documents map[string]interface{}) []string {
	var result []string
	for _, docValues := range documents {
		switch docValueMap := docValues.(type) {
		case map[string]interface{}:
			sku, hasSku := parseDocumentSku(docValueMap)
			if hasSku {
				result = append(result, sku)
			}
		default:
			// do nothing, the SKUs cannot be located here
		}
	}
	return result
}

func parseDocumentSku(document map[string]interface{}) (string, bool) {
	// if type exists, is a map[string]interface{} with "value": "sku" as an entry, and the id has a value, return ${id.sku}, true
	// else nil, false
	rawType, typeExists := document["type"]
	typeMap, typeConversionOk := rawType.(map[string]interface{})
	if typeExists && typeConversionOk && isSku(typeMap) {
		return getSkuFromDocument(document), true
	}
	return "", false
}

func isSku(typeMap map[string]interface{}) bool {
	value, ok := typeMap["value"]
	return (ok && value == "sku")
}

func getSkuFromDocument(document map[string]interface{}) string {
	rawID := document["id"]
	idMap := rawID.(map[string]interface{})
	return idMap["value"].(string)
}

// determine if a given SKU is in stock or not by querying the API with the appropriate SKU and parsing the returned HTTP response
// sample curl for fetching SKU data
// curl 'https://www.bestbuy.com/site/canopy/component/fulfillment/add-to-cart-button/v1?skuId=6429440' \
// -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:81.0) Gecko/20100101 Firefox/81.0' -H 'Accept: */*' -H 'Accept-Language: en-US,en;q=0.5' \
// --compressed -H 'DNT: 1' -H 'Connection: keep-alive' -H 'Pragma: no-cache' -H 'Cache-Control: no-cache'
// TODO return a map of the SKUs to booleans
func getSkuStatuses(skus []string, client *http.Client) map[string]bool {
	colorGreen := "\033[32m"
	colorReset := "\033[0m"
	var statuses = map[string]bool{}
	method := "GET"
	for _, sku := range skus { // TODO refactor to use goroutines for parallel execution
		url := fmt.Sprintf("https://www.bestbuy.com/site/canopy/component/fulfillment/add-to-cart-button/v1?skuId=%v", sku)
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			statuses[sku] = false
			continue
		}
		for k, v := range bestBuyHeaders() {
			req.Header.Add(k, v)
		}
		resp, err := client.Do(req)
		if err != nil {
			statuses[sku] = false
			continue
		}
		defer resp.Body.Close()
		inStockStatus := parseHtmlForStatus(sku, resp)
		if err != nil {
			statuses[sku] = false
			continue
		}
		statuses[sku] = inStockStatus
		if inStockStatus {
			log.Println("sku", sku, "in stock at best buy:", inStockStatus)
			url = fmt.Sprintf("https://www.bestbuy.com/site/searchpage.jsp?st=%v&_dyncharset=UTF-8&_dynSessConf=&id=pcat17071&type=page&sc=Global&cp=1&nrp=&sp=&qp=&list=n&af=true&iht=y&usc=All+Categories&ks=960&keys=keys", sku)
			log.Println(string(colorGreen), "Found in stock best buy SKU at", url, string(colorReset))
			util.OpenURL(url)
		} else {
			log.Printf("sku %v not in stock at best buy", sku)
		}
		util.RandomSleep(5, 10)
	}
	return statuses
}

/**
  Sample HTML response:
  <div id="fulfillment-add-to-cart-button-17b74fad-e4fd-42ad-b516-55d57a45d67c" class="None" data-version="21.0.79">

			<div class="fulfillment-add-to-cart-button"><div><button type="button" class="btn btn-disabled v-medium add-to-cart-button" disabled="" style="padding:0 8px">Coming Soon</button></div></div>
			<script defer src="https://assets.bbystatic.com/fulfillment/add-to-cart-button/dist/client/client-fcddaff74f3c36ca2951e596be60b1f8.js"></script>
			<script>
			initializer.initializeComponent({"creatorNamespace":"fulfillment","componentId":"add-to-cart-button","contractVersion":"v1","componentVersion":"21.0.79"}, "fulfillment-add-to-cart-button-17b74fad-e4fd-42ad-b516-55d57a45d67c", "{\"app\":{\"attachModal\":{\"isShowing\":false,\"instanceNumber\":0},\"tests\":{\"sold-out-openbox\":{\"scaled\":true,\"ignore\":false},\"combo-availability-plp\":{\"scaled\":true,\"ignore\":false}},\"xboxAllAccess\":false,\"instanceId\":\"fulfillment-add-to-cart-button-17b74fad-e4fd-42ad-b516-55d57a45d67c\",\"blockLevel\":false,\"size\":\"medium\",\"disabled\":false,\"chatForServiceURL\":\"\\u002Fusw\\u002Fcovertops\",\"learnMoreURL\":\"\\u002Fsite\\u002Fhome-security-solutions\\u002Fsmart-home\\u002Fpcmcat748302047019.c?id=pcmcat748302047019\",\"attachModalEnabled\":true,\"buttonStateSourceRouteEnabled\":true,\"inHomeConsultationURL\":\"https:\\u002F\\u002Fwww.bestbuy.com\\u002Fservices\\u002Fin-home-advisor\\u002Fhome\",\"maxComboSkusExceeded\":false,\"soldOutDispatchDisabled\":true,\"reportingVariables\":{},\"pickupTodayStores\":[],\"disasterFallBackLists\":{},\"xboxAllAccessModal\":{\"isShowing\":false,\"instanceNumber\":0},\"skuPdpURL\":\"\\u002Fsite\\u002Fnvidia-geforce-rtx-3080-10gb-gddr6x-pci-express-4-0-graphics-card-titanium-and-black\\u002F6429440.p?skuId=6429440\"},\"items\":[{\"skuId\":\"6429440\",\"condition\":null,\"tags\":[\"tag\\u002Fprimary\"]}],\"buttonState\":{\"$__path\":[\"shop\",\"buttonstate\",\"v5\",\"item\",\"optional\",\"skus\",\"6429440\",\"conditions\",\"NONE\",\"destinationZipCode\",\"%20\",\"storeId\",\"%20\",\"storeZipCode\",\"%20\",\"context\",\"%20\",\"addAll\",\"false\",\"consolidated\",\"false\",\"source\",\"buttonView\",\"xboxAllAccess\",\"false\",\"buttonStateResponseInfos\",0],\"buttonState\":\"COMING_SOON\",\"displayText\":\"Coming Soon\"},\"metaLayer\":{\"env_assets\":\"https:\\u002F\\u002Fassets.bbystatic.com\",\"env_piscesURL\":\"https:\\u002F\\u002Fpisces.bbystatic.com\",\"env_appServer\":\"https:\\u002F\\u002Fwww.bestbuy.com\"}}", "en-US");
		</script>
		</div>
    Therefore I think we want to open up the first div, find the second script tag, and parse that
*/
func parseHtmlForStatus(sku string, resp *http.Response) bool {
	doc, err := html.Parse(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("unable to parse HTML for status code, assuming  %v out of stock", sku)
		return false
	}
	// recursively iterate through all the HTML nodes
	var f func(*html.Node) bool
	f = func(n *html.Node) bool {
		var result bool
		if isAvailable(n) {
			return true
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			result = f(c)
			if result {
				return result
			}
		}
		return false
	}
	return f(doc)
}

func isAvailable(node *html.Node) bool {
	return strings.Contains(node.Data, "ADD_TO_CART")
}
