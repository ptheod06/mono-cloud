package main

import (
	"fmt"
	"io/ioutil"
//	"encoding/json"
	"strconv"
	"strings"
//	"sort"
	"math/rand"
	"sync"
	"bytes"
	"context"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

	"go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
//        "go.mongodb.org/mongo-driver/mongo/readpref"
        "go.mongodb.org/mongo-driver/bson"

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



//	var dbProducts []interface{}


	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongodb:27017"))
        if err != nil {
                panic(err)
        }

	log.Info("Connected to MongoDB")

	mongoConn = client

	dbprods := mongoConn.Database("mydb").Collection("products")

	indexModel := mongo.IndexModel{
	    Keys: bson.D{{"id", 1}}}

/*
	for _, item := range catalog.Products {
		dbProducts = append(dbProducts, item)
	}

	_, err = dbprods.InsertMany(context.TODO(), dbProducts)

	if err != nil {
        	panic(err)
	}

*/

	_, err = dbprods.Indexes().CreateOne(context.TODO(), indexModel)
        if err != nil {
                panic(err)
        }

	log.Info("successfully inserted products to MongoDB")


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
	var found pb.Product

//	for i := 0; i < len(parseCatalog()); i++ {
//		if req.Id == parseCatalog()[i].Id {
//			found = parseCatalog()[i]
//		}
//	}


	start := time.Now()

	dbprods := mongoConn.Database("mydb").Collection("products")

	err := dbprods.FindOne(context.TODO(), bson.D{{"id", req.Id}}).Decode(&found)

	end := time.Since(start)
	log.Info(end)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.Id)
		}
	}


	return &found, nil
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


func AddNewProduct(req *pb.ProductNew) (error) {
	var found bool
	found = false


	var foundProd pb.Product

//	log.Info(fmt.Sprintf("+%v", req))

//	for i := 0; i < len(parseCatalog()); i++ {
//                if req.Id == parseCatalog()[i].Id {
//                        found = true
//			break
//                }
//        }


	dbprods := mongoConn.Database("mydb").Collection("products")


	err := dbprods.FindOne(context.TODO(), bson.D{{"id", req.Id}}).Decode(&foundProd)

	if err == nil {
		found = true
	}

        if err != nil {
		if err == mongo.ErrNoDocuments {
			found = false
		} else {
			log.Info(err)
		}
        }

        if found == true {
		log.Info("Product already exists")
                return status.Errorf(codes.NotFound, "product with ID %s already exists!", req.Id)
        }

	freshProd := &pb.Product{Id: req.Id,
                        Name: req.Name,
                        Description: req.Description,
                        Picture: req.Picture,
                        PriceUsd: req.PriceUsd,
                        Categories: req.Categories}

	_, err = dbprods.InsertOne(context.TODO(), freshProd)

	if (err != nil) {
		log.Info(err)
		return err
        }



	strSku, _ := strconv.Atoi(req.Id)


	err = addNewRecommendation(Product{
		Sku: strSku,
		Price: float32(req.PriceUsd.Units),
		Name: req.Name,
		Manufacturer: req.Manufacturer,
		Category: req.Categories,
		Type: req.Type})

	if (err != nil) {
		//Change this for monolith implementation
		log.Info("Message failed to be send to RabbitMQ")
		return err
	} else {

		log.Info("Message sent to RabbitMQ successfully")
	}

	log.Info(fmt.Sprintf("%+v", req))

	return nil
}

