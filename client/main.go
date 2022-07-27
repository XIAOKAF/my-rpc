package main

import (
	"fmt"
	"log"
	"my-rpc/communicate"
)

var Hello func(name string) (string, error)

func main() {
	rpc := &communicate.RpcClient{}
	conn, err := rpc.NewClient()
	rpc.Conn = conn
	if err != nil {
		panic(err)
	}
	rpc.MakeCall("hello-service", "hello", "hello")
	rpc.Call(conn, "hello")
	info, err := Hello("say hello")
	if err != nil {
		log.Fatalf("not found err:%s", err)
	}
	fmt.Println(info)
}
