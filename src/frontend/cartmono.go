package main

import (
	"fmt"
//	"net"
//	"os"
//	"time"
	//"math"
	"encoding/json"
	//"strconv"

	"github.com/redis/go-redis/v9"

	"golang.org/x/net/context"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

)

//var rClient *redis.Client


// const (
// 	defaultPort = "7070"
// 	redisAddr = "redis-cach:6379"
// )


func AddItem(in *pb.AddItemRequest) (*pb.Empty, error) {

	fmt.Println("Got Add Item request");

//	log.Infof("Session: %s", in.UserId)

	redisCtx := context.Background()

	userCart, err := rClient.Get(redisCtx, in.UserId).Result()

	var newCart pb.Cart
	
	if err == redis.Nil  {

		newItem := &pb.CartItem{
			ProductId: in.Item.ProductId,
			Quantity: in.Item.Quantity,
		}

		var allItems []*pb.CartItem

		allItems = append(allItems, newItem)

		newCart = pb.Cart{UserId: in.UserId, Items: allItems}

	} else {

		unerr := json.Unmarshal([]byte(userCart), &newCart)
		if unerr != nil {
			log.Info("Redis retutned error")
		} else {

			var exists bool

			for _, item := range newCart.Items {
				if item.ProductId == in.Item.ProductId {
					exists = true
					item.Quantity = item.Quantity + in.Item.Quantity
				}
			}

			if exists != true {
				anotherItem := &pb.CartItem{ProductId: in.Item.ProductId, Quantity: in.Item.Quantity}
				newCart.Items = append(newCart.Items, anotherItem)
			}

		}

	}

	b, _ := json.Marshal(newCart)

	seterr := rClient.Set(redisCtx, in.UserId, string(b), 0).Err()
	if seterr != nil {
		log.Info("Failed to redis")
	}

	//fmt.Println(newCart)

	return &pb.Empty{}, nil
}

func GetCart(in *pb.GetCartRequest) (*pb.Cart, error) {

	fmt.Println("Got Get cart request from MONO!!!");

	redisCtx := context.Background()

	clientCart, err := rClient.Get(redisCtx, in.UserId).Result()
	if err != nil {
		//log.Info("Something bad happened with Redis")
		return &pb.Cart{}, nil
	}

	var cleanCart pb.Cart

	reserr := json.Unmarshal([]byte(clientCart),&cleanCart)
	if reserr != nil {
		return &pb.Cart{}, nil
	}

	return &cleanCart, nil
}

func EmptyCart(in *pb.EmptyCartRequest) (*pb.Empty, error) {

	fmt.Println("Got empty cart request");

	redisCtx := context.Background()

	_, err := rClient.Do(redisCtx, "del", in.UserId).Result()

	if err != nil {
		return &pb.Empty{}, nil
	}

	return &pb.Empty{}, nil
}
