package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	example "github.com/rpcxio/rpcx-examples"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

var (
	addr       = flag.String("addr", "localhost:8972", "server address")
	consulAddr = flag.String("consulAddr", "localhost:8500", "consul address")
	basePath   = flag.String("base", "/rpcx_consul", "prefix path")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)
	// meta := "state=inactive"   设置服务状态：服务虽然活着，但不提供服务
	meta := ""
	s.RegisterName("Arith", new(example.Arith), meta)
	err := s.Serve("tcp", *addr)
	if err != nil {
		fmt.Println(err)
	}
}

func addRegistryPlugin(s *server.Server) {

	r := &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,        // 本机对外暴露的监听地址， 格式为tcp@ipaddress:port
		ConsulServers:  []string{*consulAddr}, // 集群的地址
		BasePath:       *basePath,             //服务前缀。 如果有多个项目同时使用consul，避免命名冲突，可以设置这个参数，为当前的服务设置命名空间
		Metrics:        metrics.NewRegistry(), //用来更新服务的TPS
		UpdateInterval: time.Minute,           //服务的刷新间隔， 如果在一定间隔内(当前设为2 * UpdateInterval)没有刷新,服务就会从consul中删除
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
