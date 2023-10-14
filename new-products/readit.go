package main

import 	(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)



type Product struct {
	Name string
	Sku int
	Price float32
	Category []string
	Manufacturer string
	Type string
}

func intersection(arr1, arr2 []string) []string {

	commons := []string{}

	for _, i := range arr1 {
		for _, j := range arr2 {
			if i==j {
				commons = append(commons, i)
			}

		}

	}
	return commons
}


func main() {

	var products []Product
	var similarities [][]float64

	fmt.Println(similarities)

	myfile, _ := ioutil.ReadFile("first_products.json")
	json.Unmarshal(myfile, &products)

	fmt.Println(intersection(products[0].Category, products[1].Category))

//	fmt.Println(products[0].Category)

	for i := 0; i < len(products); i++ {

		var inner_arr []float64

		for j:= 0; j < len(products); j++ {
			if i == j {
				continue
			}

			numerator := len(intersection(products[i].Category, products[j].Category))
			denominator := 0.0

			denominator = math.Sqrt(float64(len(products[i].Category) + 3)) * math.Sqrt(float64(len(products[j].Category) + 3))

			if products[i].Price == products[j].Price {
				numerator += 1
			}

			if products[i].Type == products[j].Type {
                                numerator += 1
                        }

			if products[i].Manufacturer == products[j].Manufacturer { 
                                numerator += 1
                        }

			similarity := float64(numerator) / denominator
			inner_arr = append(inner_arr, similarity)
		}

		similarities = append(similarities, inner_arr)

	}

	fmt.Println(similarities)
}
