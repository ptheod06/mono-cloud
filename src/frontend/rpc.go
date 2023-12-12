// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
//	"context"
//	"time"
	"fmt"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

	"github.com/pkg/errors"
)

const (
	avoidNoopCurrencyConversionRPC = false
)

func (fe *frontendServer) getCurrencies() ([]string, error) {
	fmt.Println("Calling all currencies")
	currs, err := getAllCurrencies()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, c := range currs.CurrencyCodes {
		if _, ok := whitelistedCurrencies[c]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

func (fe *frontendServer) getProducts() ([]*pb.Product, error) {
	resp, err := ListProducts(&pb.Empty{})
	return resp.GetProducts(), err
}

func (fe *frontendServer) getProduct(id string) (*pb.Product, error) {
	resp, err := GetProduct(&pb.GetProductRequest{Id: id})

	return resp, err
}


func (fe *frontendServer) addProduct(money *pb.Money, id, name, desc, img string, categories []string, prodType, manufacturer string) error {


	mon := &pb.ProductNew{Id: id,
                        Name: name,
                        Description: desc,
                        Picture: img,
                        PriceUsd: money,
                        Categories: categories,
			Type: prodType,
			Manufacturer: manufacturer}


	err := AddNewProduct(mon)

	if (err != nil) {
		return err
	}
	return nil
}

func (fe *frontendServer) getCart(userID string) ([]*pb.CartItem, error) {
	fmt.Println("get cart from mono requested!")
	resp, err := GetCart(&pb.GetCartRequest{UserId: userID})
	return resp.GetItems(), err
}

func (fe *frontendServer) emptyCart(userID string) error {
	_, err := EmptyCart(&pb.EmptyCartRequest{UserId: userID})
	return err
}

func (fe *frontendServer) insertCart(userID, productID string, quantity int32) error {
	_, err := AddItem(&pb.AddItemRequest{
		UserId: userID,
		Item: &pb.CartItem{
			ProductId: productID,
			Quantity:  quantity},
	})
	return err
}

func (fe *frontendServer) convertCurrency(money *pb.Money, currency string) (*pb.Money, error) {
	if avoidNoopCurrencyConversionRPC && money.GetCurrencyCode() == currency {
		return money, nil
	}

	return convertCurr(&pb.CurrencyConversionRequest{
		From:   money,
		ToCode: currency})

}

func (fe *frontendServer) getShippingQuote(items []*pb.CartItem, currency string) (*pb.Money, error) {
	quote, err := GetQuote(&pb.GetQuoteRequest{
                        Address: nil,
                        Items:   items})

	if err != nil {
		return nil, err
	}
	localized, err := fe.convertCurrency(quote.GetCostUsd(), currency)
	return localized, errors.Wrap(err, "failed to convert currency for shipping cost")
}

func (fe *frontendServer) getRecommendations(userID string, productIDs []string) ([]*pb.Product, error) {
	resp, err := ListRecommendations(&pb.ListRecommendationsRequest{UserId: userID, ProductIds: productIDs})
	if err != nil {
		return nil, err
	}
	out := make([]*pb.Product, len(resp.GetProductIds()))
	for i, v := range resp.GetProductIds() {
		p, err := fe.getProduct(v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get recommended product info (#%s)", v)
		}
		out[i] = p
	}
	if len(out) > 4 {
		out = out[:4] // take only first four to fit the UI
	}
	return out, err
}

