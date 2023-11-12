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
	"fmt"
	//"net"
	"os"
//	"time"
	"math"
	"encoding/json"
	"strconv"


	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)


type Money_Temp struct {
	units float64
	nanos float64
}


var (
	converted_money Money_Temp
)

//Helper function to read currencies from file
func getCurrencyData() (map[string]float64){



	da, err := os.ReadFile("curr.json")

	data := make(map[string]string)

	currencies := make(map[string]float64)

	json.Unmarshal(da, &data)

	if err != nil {
		fmt.Println("Could not read currencies file")
		fmt.Println(err)
		return nil
	} else {
		for key, val := range data {
			converted, _ := strconv.ParseFloat(val, 64)

			currencies[key] = converted

		}

	}

	return currencies
}

// Convert between currencies
func convert_helper(units int64, nanos int32, currency_code string, to_code string) {

//	log.Info("Units %v", units)

	var currency_list map[string]float64

	currency_list = getCurrencyData()

	var nAmount = Money_Temp{(float64(units)/currency_list[currency_code]), (float64(nanos)/currency_list[currency_code])};

	var ineuros = carry(nAmount)

	ineuros.nanos = math.Round(ineuros.nanos)

	var nResult = Money_Temp{(ineuros.units * currency_list[to_code]), (ineuros.nanos * currency_list[to_code])}

//	log.Info("Before: %v", nResult.units)

	converted_money = carry(nResult)

//	log.Info("Converted Units", converted_money.units)
//	log.Info("Converted Nanso", converted_money.nanos)

//	log.Info("Converted: %v", converted_money.units)

	converted_money.units = math.Floor(converted_money.units)
	converted_money.nanos = math.Floor(converted_money.nanos)

//	log.Info("converted unis", converted_money.units)
//	log.Info("converted nanos", converted_money.nanos)

}

//Handle carry operations
func carry(amount Money_Temp) (Money_Temp) {

	var fractionSize float64
	fractionSize = math.Pow(10, 9)

	amount.nanos += math.Mod(amount.units, 1) * fractionSize

	amount.units = math.Floor(amount.units) + math.Floor(amount.nanos / fractionSize)

	amount.nanos = math.Mod(amount.nanos, fractionSize)

//	log.Info("Units: %d", amount.units)
//	log.Info("Nanos: %d", amount.nanos)

	return amount

}




// Get all supported currencies produces a shipping quote (cost) in USD.
func getAllCurrencies() (*pb.GetSupportedCurrenciesResponse, error) {

//	log.Info("Got request on all currencies on MONO")

	var curr_codes []string

	var all_currencies map[string]float64

	all_currencies = getCurrencyData()

	for key, _ := range all_currencies {
		curr_codes = append(curr_codes, key)
	}
	
	return &pb.GetSupportedCurrenciesResponse{CurrencyCodes: curr_codes}, nil
	

}

func convertCurr(in *pb.CurrencyConversionRequest) (*pb.Money, error) {

//	log.Info("Got request on conversion")

	convert_helper(in.From.Units, in.From.Nanos, in.From.CurrencyCode, in.ToCode);

	return &pb.Money{
		CurrencyCode: in.ToCode,
		Units: int64(converted_money.units),
		Nanos: int32(converted_money.nanos),
		}, nil

}
