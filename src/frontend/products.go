package main

import (
	"fmt"
	"io/ioutil"
//	"encoding/json"
	"strconv"
	"strings"
	"sort"
	"math/rand"
	"sync"
	"bytes"
//	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	cat pb.ListProductsResponse
	catalogMutex *sync.Mutex
	reloadCatalog bool
)

type Product struct {
        Name string
        Sku int
        Price float32
        Category []string
        Manufacturer string
        Type string
}


func readCatalogFile(catalog *pb.ListProductsResponse) error {
	catalogMutex.Lock()
	defer catalogMutex.Unlock()
	catalogJSON, err := ioutil.ReadFile("bestProducts.json")
	if err != nil {
		log.Fatalf("failed to open product catalog json file: %v", err)
		return err
	}
	if err := jsonpb.Unmarshal(bytes.NewReader(catalogJSON), catalog); err != nil {
		log.Warnf("failed to parse the catalog JSON: %v", err)
		return err
	}


	sort.Slice(catalog.Products, func(i, j int) bool { 
		firstVal, _ := strconv.Atoi(catalog.Products[i].Id)
		secondVal, _ := strconv.Atoi(catalog.Products[j].Id)
		return firstVal < secondVal })

	isSorted := sort.SliceIsSorted(catalog.Products, func (i, j int) bool {
		firstVal, _ := strconv.Atoi(catalog.Products[i].Id)
                secondVal, _ := strconv.Atoi(catalog.Products[j].Id)
                return firstVal < secondVal })

	if (isSorted) {
		log.Info("Sorted Successfully")
	} else {
		log.Info("Sort did not work!!!!")
	}

	log.Info("successfully parsed product catalog json")
	return nil
}


func parseCatalog() []*pb.Product {
	if reloadCatalog || len(cat.Products) == 0 {
		err := readCatalogFile(&cat)
		if err != nil {
			return []*pb.Product{}
		}
	}
	return cat.Products
}

func insertProduct(index int, prod *pb.Product) error {
	products := parseCatalog()

	if (len(products) == index) {
		products = append(products, prod)
		return nil
	}

	products = append(products[:index+1], products[index:]...)
	products[index] = prod

	return nil


}


func ListProducts(*pb.Empty) (*pb.ListProductsResponse, error) {
	//Should be remove from microservices as well
	//time.Sleep(extraLatency)

	startNum := rand.Intn(1000)
	endNum := startNum + 10
	return &pb.ListProductsResponse{Products: parseCatalog()[startNum:endNum]}, nil
}


func GetProduct(req *pb.GetProductRequest) (*pb.Product, error) {
//	time.Sleep(extraLatency)
	var found *pb.Product

//	for i := 0; i < len(parseCatalog()); i++ {
//		if req.Id == parseCatalog()[i].Id {
//			found = parseCatalog()[i]
//		}
//	}

	i := sort.Search(len(parseCatalog()), func (ind int) bool {
		firstVal, _ := strconv.Atoi(parseCatalog()[ind].Id)
		secondVal, _ := strconv.Atoi(req.Id)
		return firstVal >= secondVal })

	if (i < len(parseCatalog()) && parseCatalog()[i].Id == req.Id) {

//		log.Info(fmt.Sprintf("Found at index: %d", i))
		found = parseCatalog()[i]
	} else {

//		log.Info("Product Not Found")
//		log.Info(fmt.Sprintf("index: %d", i))
	}


	if found == nil {
		return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.Id)
	}
	return found, nil
}


func SearchProducts(req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
//	time.Sleep(extraLatency)
	// Intepret query as a substring match in name or description.
	var ps []*pb.Product
	for _, p := range parseCatalog() {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(req.Query)) {
			ps = append(ps, p)
		}
	}
	return &pb.SearchProductsResponse{Results: ps}, nil
}


func AddNewProduct(req *pb.ProductNew) (*pb.Empty, error) {
	var found bool
	found = false

//	log.Info(fmt.Sprintf("+%v", req))

//	for i := 0; i < len(parseCatalog()); i++ {
//                if req.Id == parseCatalog()[i].Id {
//                        found = true
//			break
//                }
//        }


	i := sort.Search(len(parseCatalog()), func (ind int) bool {
                firstVal, _ := strconv.Atoi(parseCatalog()[ind].Id)
                secondVal, _ := strconv.Atoi(req.Id)
                return firstVal >= secondVal })

        if (i < len(parseCatalog()) && parseCatalog()[i].Id == req.Id) {

        	log.Info(fmt.Sprintf("Found at index: %d", i))
                found = true
        } else {

//              log.Info("Product Not Found")
//              log.Info(fmt.Sprintf("index: %d", i))
        }


        if found == true {
		log.Info("Product already exists")
                return nil, status.Errorf(codes.NotFound, "product with ID %s already exists!", req.Id)
        }

	freshProd := &pb.Product{Id: req.Id,
                        Name: req.Name,
                        Description: req.Description,
                        Picture: req.Picture,
                        PriceUsd: req.PriceUsd,
                        Categories: req.Categories}

	err := insertProduct(i, freshProd)


	if (err != nil) {
		log.Info("product added successfully")
        }

//	strSku, _ := strconv.Atoi(req.Id)
	
/*
	err = sendMsgToQueue(Product{
		Sku: strSku,
		Price: float32(req.PriceUsd.Units),
		Name: req.Name,
		Manufacturer: req.Manufacturer,
		Category: req.Categories,
		Type: req.Type})

	if (err != nil) {

		log.Info("Message failed to be send to RabbitMQ")
		return nil, err
	} else {

		log.Info("Message sent to RabbitMQ successfully")
	}

*/

	log.Info(fmt.Sprintf("%+v", req))

	return &pb.Empty{}, nil
}

