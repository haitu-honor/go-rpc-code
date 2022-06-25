package main

import (
	"context"
	"flag"
	"log"
	"time"

	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/client"
)

var (
	consulAddr = flag.String("consulAddr", "localhost:8500", "consul address")
	basePath   = flag.String("base", "/rpcx_consul", "prefix path")
)

func main() {
	flag.Parse()

	// option := client.DefaultOption
	// option.Heartbeat = true   // 启用心跳
	// option.HeartbeatInterval = time.Second * 5   // 设置心跳间隔
	d, _ := client.NewConsulDiscovery(*basePath, "Arith", []string{*consulAddr}, nil)
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)

	defer xclient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	for {
		reply := &example.Reply{}
		// context.WithTimeout() 设置请求超时时间
		err := xclient.Call(context.Background(), "Mul", args, reply)
		if err != nil {
			log.Printf("ERROR failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(1e9)
	}

}
