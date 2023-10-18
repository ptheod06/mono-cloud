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


func main() {

	var products []Product


	myfile, _ := ioutil.ReadFile("final_products.json")
	json.Unmarshal(myfile, &products)

//	fmt.Println(products[51640:])

	sort.Slice(products, func(i, j int) bool { return products[i].Sku < products[j].Sku })

	before := time.Now()

	prod := sort.Search(len(products), func(i int) bool { return products[i].Sku >= 2736642})


//	for i := 0; i < len(products); i++ {
//		if (products[i].Sku == 150115) {
//			fmt.Println(products[i])
//			break
//		}

//	}

//	fmt.Println(products[51640:])

	after := time.Since(before)


	if (prod < len(products) && products[prod].Sku == 2736642) {
		fmt.Println(products[prod])
	} else {
		fmt.Println("not found")
	}

	fmt.Println(after)
}
