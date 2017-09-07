package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"go-orderbook/trades"
)

type Product struct {
	ID             string `json:"id"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseMinSize    string `json:"base_min_size"`
	BaseMaxSize    string `json:"base_max_size"`
	QuoteIncrement string `json:"quote_increment"`
	DisplayName    string `json:"display_name"`
	MarginEnabled  bool   `json:"margin_enabled"`
}

// GetProducts
func getProducts() (*[]Product, error) {
	resp, err := http.Get("https://api.gdax.com/products")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var products []Product
	if err = json.Unmarshal(body, &products); err != nil {
		return nil, err
	}
	return &products, nil
}

func main() {

	var download bool
	flag.BoolVar(&download, "d", false, "download second level orders into csv")
	flag.Parse()

	products, err := getProducts()
	if err != nil {
		log.Fatal(err)
	}

	if download {
		ch := make(chan string)
		for _, product := range *products {
			fmt.Println("Start fetch for ", product.ID)
			go processLevel2OrderBook(product.ID, ch)
		}

		for i := 0; i < len(*products); i++ {
			status := <-ch
			fmt.Println(status)
		}
		log.Println("Done")
		return
	} else {
		for _, product := range *products {
			processLevel1OrderBook(product.ID)
		}
		fmt.Println("Run make client")
		fmt.Println("Example: 127.0.0.1:6379> HGETALL LTC-EUR")
	}
}

func processLevel1OrderBook(productID string) {
	fmt.Println("Getting top order for", productID)
	orders, err := trades.GetOrderBook("1", productID)
	if err != nil {
		log.Println(err)
	}
	err = orders.RecordTopBidinRedis(productID)
	if err != nil {
		log.Println(err)
	}
}

func processLevel2OrderBook(productID string, ch chan string) {
	orders, err := trades.GetOrderBook("2", productID)
	if err != nil {
		log.Println(err)
	}
	path := filepath.Join("csv", productID)
	err = orders.WriteToCSV(path)
	if err != nil {
		log.Println(err)
	}
	ch <- productID + "- Done"
}
