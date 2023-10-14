package main

import 	(
	"encoding/json"
	"fmt"
	"io/ioutil"
)



type Product struct {
	Name string
	Sku int
	Price float32
	Category []string
	Manufacturer string
	Type string
}


func main() {

	var products []Product

	myfile, _ := ioutil.ReadFile("first_products.json")
	json.Unmarshal(myfile, &products)

	fmt.Println(products)

}
