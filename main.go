package main

import (
	"errors"
	"fmt"
	"my-rpc/communicate"
)

type i interface {
	Hello(req string) string
}

type me struct {
}

func main() {
	s := communicate.NewCenter()
	err := s.Run()
	if err != nil {
		panic(err)
	}
	s.Register("hello-service", "hello", &me{})
}

func (m *me) Hello(req string) error {
	if req != "" {
		fmt.Println("hello")
	}
	return errors.New("invalid request")
}
