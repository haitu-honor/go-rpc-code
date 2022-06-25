package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	pb "github.com/haitu-honor/go-rpc-code/grpc-protobuf-exp/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	// 公钥中读取和解析公钥/私钥对
	pair, err := tls.LoadX509KeyPair("../../keys/server.crt", "../../keys/server.key")
	if err != nil {
		fmt.Println("LoadX509KeyPair error", err)
		return
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../keys/ca.crt")
	if err != nil {
		fmt.Println("read ca pem error ", err)
		return
	}
	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("AppendCertsFromPEM error ")
		return
	}
	// 创建credentials 对象
	// ClientAuth有5种类型，如果要进行双向认证必须是RequireAndVerifyClientCert
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair},        //服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert, // 需要并且验证客户端证书
		ClientCAs:    certPool,                       // 客户端证书池
	})
	// 实例化grpc Server,并开启TLS认证
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	// 监听端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	// 注册HelloService
	pb.RegisterHelloServer(grpcServer, HelloService)
	grpclog.Infoln("Listen on " + Address) // grcp 日志打印
	fmt.Println("Listen on " + Address + "with TLS")
	// 启动服务
	grpcServer.Serve(listen)
}
