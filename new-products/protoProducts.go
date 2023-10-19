package main

import 	(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)


type FinalOut struct {
	Products []outProduct `json:"products"`

}

type Product struct {
	Name string
	Sku int
	Price float32
	Category []string
	Manufacturer string
	Type string
	Description string
	Image string
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

	myfile, _ := ioutil.ReadFile("cleaned_products.json")
	json.Unmarshal(myfile, &products)

	for i := 0; i < len(products); i++ {

		units := int(math.Trunc(float64(products[0].Price)))
	        nanos := int(products[0].Price * 100) % 100 * 10000000

		output = append(output, outProduct{strconv.Itoa(products[i].Sku), products[i].Name, products[i].Description, products[i].Image, Currency{"USD", units, nanos}, products[i].Category})

	}

	Outputfinal := FinalOut{output}


	file_content, _ := json.MarshalIndent(Outputfinal, "", " ")
	err := ioutil.WriteFile("newwproducts.json", file_content, 0644)
	if err != nil {
		fmt.Println("oops")
	}

}
