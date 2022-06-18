package main

import (
	"context"
	"flag"
	"log"

	"github.com/haitu-honor/go-rpcx-code/probuf-example/pb"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	// flag.Parse()

	// register customized codec
	option := client.DefaultOption
	// 使用 ProtoBuffer编解码器 进行编解码
	option.SerializeType = protocol.ProtoBuffer

	// 1、 定义了使用什么方式来实现服务发现。
	// 在这里我们使用最简单的 Peer2PeerDiscovery（点对点）。客户端直连服务器来获取服务地址。
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	// 第一个参数：要连接的服务名  第二、三个参数：分别设置服务治理的失败模式与负载均衡
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := &pb.ProtoArgs{
		A: 10,
		B: 20,
	}

	reply := &pb.ProtoReply{}
	// context.Background()返回一个非零的空上下文。 它永远不会被取消，没有值，也没有最后期限。
	// 它通常由主函数、初始化和测试使用，并作为传入的顶级上下文
	// 第二个参数：服务的处理函数名
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Printf("%d * %d = %d", args.A, args.B, reply.C)

}
