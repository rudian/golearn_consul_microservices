package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
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
	consulConfig := api.DefaultConfig()

	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("consul client error:", err)
		return
	}

	healthCheckId := "check1"
	HealthCheckTTl := 10 * time.Second
	reg := api.AgentServiceRegistration{
		Tags:    []string{"aa", "bb"},
		Name:    "hello service",
		Address: "127.0.0.1",
		Port:    3000,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: HealthCheckTTl.String(),
			TLSSkipVerify:                  true,
			TTL:                            HealthCheckTTl.String(),
			CheckID:                        healthCheckId,
		},
	}

	errAgent := consulClient.Agent().ServiceRegister(&reg)
	if errAgent != nil {
		fmt.Println("consul register error:", errAgent)
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			err := consulClient.Agent().UpdateTTL(
				healthCheckId,
				"online",
				api.HealthPassing,
			)
			if err != nil {
				log.Fatalln(err)
				return
			}
			<-ticker.C
		}
	}()

	//--------------

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServer(grpcServer, &Hello{})

	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println("listener error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Service starting successfully")

	errGrpc := grpcServer.Serve(listener)
	if errGrpc != nil {
		fmt.Println("grpc server error:", errGrpc)
		return
	}
}
