package main

import (
//	"fmt"
	//"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/durango/go-credit-card"

	"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"
)


func Charge(in *pb.ChargeRequest) (*pb.ChargeResponse, error) {

//	fmt.Println("Got Charge request");

	cardYear := strconv.Itoa(int(in.CreditCard.CreditCardExpirationYear))
	cardMonth := strconv.Itoa(int(in.CreditCard.CreditCardExpirationMonth))

	cardNo := strings.ReplaceAll(in.CreditCard.CreditCardNumber, "-", "")

	card := creditcard.Card{Number: cardNo, Cvv: "123", Month: cardMonth, Year: cardYear}

	err := card.Method()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to validate credit card number")
	}


	if card.Company.Short != "visa" && card.Company.Short != "mastercard" {
		return nil, status.Errorf(codes.Internal, "Only accepting visa and mastercard")
	}

	expErr := card.ValidateExpiration()

	if expErr != nil {
		return nil, status.Errorf(codes.Internal, "Card is expired")
	}

	val, _ := uuid.NewRandom()

	cleaned := val.String()

	return &pb.ChargeResponse{TransactionId: cleaned}, nil


}
