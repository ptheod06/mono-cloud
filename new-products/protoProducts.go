package main

import 	(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"sort"
)


type Product struct {
	Name string
	Sku int
	Price float32
	Category []string
	Manufacturer string
	Type string
}

type Currency struct {
	CurrencyCode string `json:"currencyCode"`
	Units int `json:"units"`
	Nanos int `json:"nanos"`
}

type outProduct struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Picture string `json:"picture"`
	PriceUsd Currency `json:"priceUsd"`
	Categories []string `json:"categories"`
}


func main() {

	var products []Product
	var output []outProduct

	myfile, _ := ioutil.ReadFile("final_products.json")
	json.Unmarshal(myfile, &products)

	for i := 0; i < 15; i++ {

		output = append(output, outProduct{products[i].Sku, products[i].Name, "", })

	}


}
