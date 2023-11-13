package main

import (

	"io/ioutil"
	"encoding/json"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
//        "go.mongodb.org/mongo-driver/mongo/readpref"
//        "go.mongodb.org/mongo-driver/bson"

)

type Product struct {
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Picture     string `protobuf:"bytes,4,opt,name=picture,proto3" json:"picture,omitempty"`
	PriceUsd    *Money `protobuf:"bytes,5,opt,name=price_usd,json=priceUsd,proto3" json:"priceUsd,omitempty"`
	// Categories such as "clothing" or "kitchen" that can be used to look up
	// other related products.
	Categories           []string `protobuf:"bytes,6,rep,name=categories,proto3" json:"categories,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}


type ListProductsResponse struct {
	Products             []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

type Money struct {
	// The 3-letter currency code defined in ISO 4217.
	CurrencyCode string `protobuf:"bytes,1,opt,name=currency_code,json=currencyCode,proto3" json:"currencyCode,omitempty"`
	// The whole units of the amount.
	// For example if `currencyCode` is `"USD"`, then 1 unit is one US dollar.
	Units int64 `protobuf:"varint,2,opt,name=units,proto3" json:"units,omitempty"`
	// Number of nano (10^-9) units of the amount.
	// The value must be between -999,999,999 and +999,999,999 inclusive.
	// If `units` is positive, `nanos` must be positive or zero.
	// If `units` is zero, `nanos` can be positive, zero, or negative.
	// If `units` is negative, `nanos` must be negative or zero.
	// For example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.
	Nanos                int32    `protobuf:"varint,3,opt,name=nanos,proto3" json:"nanos,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}


var catalog ListProductsResponse


func main() {


	myfile, _ := ioutil.ReadFile("bestProducts.json")
	json.Unmarshal(myfile, &catalog)

	fmt.Printf("%+v", catalog.Products[10].PriceUsd)


	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongodb:27017"))
        if err != nil {
                panic(err)
        }


        dbprods := client.Database("mydb").Collection("products")


	var dbProducts []interface{}

	for _, item := range catalog.Products {
                dbProducts = append(dbProducts, item)
        }

        _, err = dbprods.InsertMany(context.TODO(), dbProducts)

        if err != nil {
                panic(err)
        }


}
