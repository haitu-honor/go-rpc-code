package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	pb "github.com/haitu-honor/go-rpc-code/grpc-protobuf-exp/proto/hello" // 引入proto包
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	// 公钥中读取和解析公钥/私钥对
	pair, err := tls.LoadX509KeyPair("../../keys/client.crt", "../../keys/client.key")
	if err != nil {
		fmt.Println("LoadX509KeyPair error ", err)
		return
	}
	// 创建一组根证书
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../keys/ca.crt")
	if err != nil {
		fmt.Println("ReadFile ca.crt error ", err)
		return
	}
	// 解析证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("certPool.AppendCertsFromPEM error ")
		return
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair}, // 放入客户端证书
		ServerName:   "grpctls.com",           // 这里的参数为server.conf里的 DNS
		RootCAs:      certPool,                // 证书池
	})
	// tls连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()
	// 初始化客户端
	c := pb.NewHelloClient(conn)
	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		grpclog.Fatalln(err)
	}
	grpclog.Infoln(res.Message) // grcp 日志打印
	fmt.Println(res.Message)
}
