package main

import (
	"context"
	"fmt"
	"golearn_consul_microservices/consul"
	"golearn_consul_microservices/proto_generated/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Hello struct {
	pb.UnimplementedHelloServer
}

func (this *Hello) mustEmbedUnimplementedHelloServer() {
	//TODO implement me
	panic("implement me")
}

func (this *Hello) SayHello(ctx context.Context, p *pb.Person) (*pb.Person, error) {
	fmt.Println(p.Name)
	p.Name = "hello" + p.Name
	return p, nil
}

func main() {
	//1. 启动consul，注入consul服务地址
	consulService, _ := consul.NewService("127.0.0.1:8500")
	//2. 注册服务
	err := consulService.RegisterService(consul.RegisterService{
		ServiceName:   "hello_service",
		Address:       "127.0.0.1",
		Port:          3000,
		HeathCheckTTL: 10 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	//3. 启动grpc服务器
	grpcServer := grpc.NewServer()

	//4. 把hello和grpc服务器绑定
	pb.RegisterHelloServer(grpcServer, &Hello{})

	//5. 启动端口监听
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("listener error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Service starting successfully")

	//6. grpc和端口监听绑定
	errGrpc := grpcServer.Serve(listener)
	if errGrpc != nil {
		fmt.Println("grpc server error:", errGrpc)
		return
	}
}
