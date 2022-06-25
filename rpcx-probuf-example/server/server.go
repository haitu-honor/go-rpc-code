package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/haitu-honor/go-rpcx-code/rpcx-probuf-example/pb"

	"github.com/smallnest/rpcx/server"
)

var (
	// flag go的内置库，可以参考这篇文章 https://darjun.github.io/2020/01/10/godailylib/flag/
	addr = flag.String("addr", "localhost:8972", "server address")
)

type Arith int

// 接受3个参数，第一个是 context.Context类型，其他2个都是可导出（或内置）的类型。
func (t *Arith) Mul(ctx context.Context, args *pb.ProtoArgs, reply *pb.ProtoReply) error {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func main() {
	// 从os.Args[1:]中解析注册的flag。必须在所有flag都注册好而未访问其值时执行
	flag.Parse()
	// 注册一个服务命为 Arith 的服务
	s := server.NewServer()
	//s.RegisterName("Arith", new(example.Arith), "")
	s.Register(new(Arith), "")
	s.Serve("tcp", *addr)
}
