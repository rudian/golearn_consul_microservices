package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"golearn_consul_microservices/proto_generated/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

func main() {
	consulConfig := api.DefaultConfig()

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	service, _, err := consulClient.Health().Service("hello service", "aa", true, nil)
	if err != nil {
		return
	}

	fmt.Println(service[0].Service)
	address := service[0].Service.Address + ":" + strconv.Itoa(service[0].Service.Port)

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewHelloClient(conn)
	hello, err := grpcClient.SayHello(context.TODO(), &pb.Person{
		Name:   "kiat",
		Age:    18,
		Length: 15.6,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(hello)
}
