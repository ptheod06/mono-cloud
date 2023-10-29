package main

import (

//	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"math"
	"sort"
	"sync"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

)


type ProdOut struct {

	Sku int
	Similarity float64
}

type SimProducts struct {

	Sku int
	SimilarProducts []ProdOut
}

type ProductComp struct {

	Item int
	Similarity float64
}


var (

	prods []SimProducts
	products []Product
	mu sync.RWMutex

)



func readProducts(allProducts *[]Product) error {
	catalogJSON, err := ioutil.ReadFile("finaldetails_products.json")
        if err != nil {
                log.Fatalf("failed to open products catalog json file: %v", err)
                return err
        }
        if err := json.Unmarshal(catalogJSON, allProducts); err != nil {
                log.Warnf("failed to parse the catalog JSON: %v", err)
                return err
        }


        log.Info("successfully parsed products catalog json for Recommendation service")
        return nil



}

// Not used since recommendations are generated dynamically

// func readRecommFile(catalog *[]SimProducts) error {
// 	catalogJSON, err := ioutil.ReadFile("recommProducts.json")
// 	if err != nil {
// 		log.Fatalf("failed to open recommendations catalog json file: %v", err)
// 		return err
// 	}
// 	if err := json.Unmarshal(catalogJSON, catalog); err != nil {
// 		log.Warnf("failed to parse the catalog JSON: %v", err)
// 		return err
// 	}
// 	log.Info("successfully parsed recommendations catalog json")
// 	return nil
// }

func addNewRecommendation(productToAdd Product) error {

	products = append(products, productToAdd)

	go calculateSimilarities(&prods)

	return nil
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


func calculateSimilarities(allSimilarities *[]SimProducts) {

	var similarities [][]float64
	var topSimilarities [][]ProductComp


	log.Info("Starting to calculate Recommendations")

	// before := time.Now()

	for i := 0; i < len(products); i++ {

		var inner_arr []float64


		for j:= 0; j < len(products); j++ {
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


	for _, arrays := range similarities {

		var mostSimilar = []ProductComp{}

		for i := 0; i < 5; i++ {
			prod := ProductComp{i, arrays[i]}
			mostSimilar = append(mostSimilar, prod)
		}

		sort.Slice(mostSimilar, func(i, j int) bool { return mostSimilar[i].Similarity < mostSimilar[j].Similarity})
		for i := 5; i < len(arrays); i++ {
			if (mostSimilar[0].Similarity < arrays[i]) {
				prodNew := ProductComp{i, arrays[i]}
				mostSimilar[0] = prodNew
				sort.Slice(mostSimilar, func(i, j int) bool { return mostSimilar[i].Similarity < mostSimilar[j].Similarity })
			}

		}

		topSimilarities = append(topSimilarities, mostSimilar)
	}

	// after := time.Since(before)

	var allNewSimilarities = []SimProducts{}


	for i := 0; i < len(products); i++ {
		var currSimilarities = []ProdOut{}

		for _, item := range topSimilarities[i] {

			currSimilarities = append(currSimilarities, ProdOut{products[item.Item].Sku, item.Similarity})
		}

		allNewSimilarities = append(allNewSimilarities, SimProducts{products[i].Sku, currSimilarities})

	}

	sort.Slice(allNewSimilarities, func (i, j int) bool {
		return allNewSimilarities[i].Sku < allNewSimilarities[j].Sku })


	mu.Lock()
	defer mu.Unlock()

	*allSimilarities = allNewSimilarities

}


func ListRecommendations(in *pb.ListRecommendationsRequest) (*pb.ListRecommendationsResponse, error) {

	var similarProducts []string

	mu.RLock()
	defer mu.RUnlock()

//	for _, item := range prods {
//		if (in.ProductIds[0] == strconv.Itoa(item.Sku)) {
//			for i := 0; i < len(item.SimilarProducts); i++ {
//				similarProducts = append(similarProducts, strconv.Itoa(item.SimilarProducts[i].Sku))
//			}
//			break
//		}
//
//	}

	if (len(in.ProductIds) < 1) {
		return &pb.ListRecommendationsResponse{ProductIds: []string{"2099128"}}, nil
	}

	idStr, _ := strconv.Atoi(in.ProductIds[0])

	foundAt := sort.Search(len(prods), func (ind int) bool {
		return prods[ind].Sku >= idStr
	})

	if (foundAt < len(prods) && prods[foundAt].Sku == idStr) {
		for i := 0; i < len(prods[foundAt].SimilarProducts); i++ {
                	similarProducts = append(similarProducts, strconv.Itoa(prods[foundAt].SimilarProducts[i].Sku))
                }
	}


	if (len(similarProducts) < 4) {

		return &pb.ListRecommendationsResponse{ProductIds: []string{"2099128"}}, nil
	} else {

		return &pb.ListRecommendationsResponse{ProductIds: similarProducts[:4]}, nil
	}
}

