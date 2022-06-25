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

	d, _ := client.NewConsulDiscoveryTemplate(*basePath, []string{*consulAddr}, nil)
	oneClient := client.NewOneClient(client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer oneClient.Close()

	args := &example.Args{
		A: 10,
		B: 20,
	}

	reply := &example.Reply{}
	for {
		err := oneClient.Call(context.Background(), "Arith", "Mul", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("%d * %d = %d", args.A, args.B, reply.C)
		time.Sleep(time.Second)
	}

}
