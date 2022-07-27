package config

import "errors"

var (
	InvalidInterface = errors.New("invalid interface")
	MsgTypeErr       = errors.New("message type is wrong")

	LackOfParams     = errors.New("params are not enough")
	SendDataErr      = errors.New("fail to send data")
	ServiceNotFound  = errors.New("service is not found")
	MethodIsNotFound = errors.New("method is not found")

	HandleErr   = errors.New("fail to handle connection")
	ReadDataErr = errors.New("fail to read data")

	DecodeErr = errors.New("fail to decode data")
)
