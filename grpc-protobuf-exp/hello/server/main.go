package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/haitu-honor/go-rpc-code/grpc-protobuf-exp/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
// 很多教程实现了mustEmbedUnimplementedHelloServer 这个方法，但是由于是小写， 同目录下是好的，跨了目录就会有问题。
// 所以直接在 helloService 里加入 pb.UnimplementedHelloServer，不同项目这个方法的名字也不一样
type helloService struct {
	pb.UnimplementedHelloServer
}

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)
	return resp, nil
}

// func (h helloService) mustEmbedUnimplementedHelloServer() {}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	// 实例化grpc Server
	s := grpc.NewServer()
	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)
	grpclog.Infoln("Listen on " + Address)
	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
