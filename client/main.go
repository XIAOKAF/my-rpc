package main

import (
	"fmt"
	"log"
	"my-rpc/communicate"
	"net"
)

var Hello func(name string) (string, error)

func main() {
	rpc := &communicate.RpcClient{}
	conn, err := rpc.NewClient()
	rpc.Conn = conn
	if err != nil {
		panic(err)
	}
	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}
	c, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	rpc.MakeCall("hello-service", "hello", "hello")
	rpc.Call(c, "hello")
	info, err := Hello("say hello")
	if err != nil {
		log.Fatalf("not found err:%s", err)
	}
	fmt.Println(info)
}
