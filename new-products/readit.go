package main

import 	(
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
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

type ProductComp struct {

	Sku string
	Similarity float64
	Index int
	With int
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
	var topSimilarities [][]ProductComp
	//var similarSKU [][]string
	//fmt.Println(similarities)

	myfile, _ := ioutil.ReadFile("final_products.json")
	json.Unmarshal(myfile, &products)

	//fmt.Println(intersection(products[0].Category, products[1].Category))

//	fmt.Println(products[0].Category)

	before := time.Now()

	for i := 0; i < 100; i++ {

		var inner_arr []float64


		for j:= 0; j < 100; j++ {
			if i == j {
				inner_arr = append(inner_arr, -1.0)
				continue
			}

			numerator := len(intersection(products[i].Category, products[j].Category))

			denominator := math.Sqrt(float64(len(products[i].Category) + 3)) * math.Sqrt(float64(len(products[j].Category) + 3))

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


	for index, arrays := range similarities {

		var mostSimilar = []ProductComp{}

		for i := 0; i < 3; i++ {
			prod := ProductComp{products[i].Sku, arrays[i], index, i}
			mostSimilar = append(mostSimilar, prod)
		}

		sort.Slice(mostSimilar, func(i, j int) bool { return mostSimilar[i].Similarity < mostSimilar[j].Similarity})
		for _, item := range arrays {
			if (mostSimilar[0].Similarity < item) {
				prodNew := ProductComp{products[i].Sku, arrays[i], index, i}
				mostSimilar[0] = prodNew
				sort.Slice(mostSimilar, func(i, j int) bool { return mostSimilar[i].Similarity < mostSimilar[j].Similarity })
			}

		}

		topSimilarities = append(topSimilarities, mostSimilar)
	}

	after := time.Since(before)

	fmt.Println(after)

//	fmt.Println(topSimilarities)

	file_content, _ := json.Marshal(similarities)
	err := ioutil.WriteFile("output.json", file_content, 0644)
	if err != nil {
		fmt.Println("oops")
	}
//	fmt.Println(similarities)


	//mine := ProductComp{"aaaww", 12.33, 1}
	//fmt.Println(mine)

}
