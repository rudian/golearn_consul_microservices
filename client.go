package main

import (
	"context"
	"fmt"
	"golearn_consul_microservices/consul"
	"golearn_consul_microservices/proto_generated/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	//1. 启动consul，注入consul服务地址
	consulService, err := consul.NewService("127.0.0.1:8500")
	if err != nil {
		log.Fatal(err)
		return
	}

	//2. 从consul获取服务的地址和端口
	serviceAddress, err := consulService.GetServiceAddress("hello_service")
	if err != nil {
		log.Fatal(err)
		return
	}

	//3. 往微服务端口尽量链接
	conn, err := grpc.NewClient(
		serviceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	//4. 把grpc链接和微服务function做绑定
	grpcClient := pb.NewHelloClient(conn)

	//5. 调用微服务的function，hello获得service微服务的回传
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
