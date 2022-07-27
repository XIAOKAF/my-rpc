package communicate

import (
	"log"
	"my-rpc/config"
	"net"
	"reflect"
)

type RpcClient struct {
	Msg  *RpcMsg
	Conn net.Conn
}

var rpc RpcMsg

func (r *RpcClient) NewClient() (net.Conn, error) {
	conn, err := net.Dial("tcp", ":8001")
	return conn, err
}

func (r *RpcClient) MakeCall(serviceName, methodName string, i interface{}) {
	err := r.sendData(serviceName, methodName, i)
	if err != nil {
		log.Fatalf("%s:%s", config.SendDataErr, err)
	}
}

func (r *RpcClient) ReadData(conn net.Conn) (reflect.Value, error) {
	_, err := r.handleConn(conn)
	if err != nil {
		return r.Msg.Method.Method, err
	}
	//判断消息类型
	flag := IsReqFromCli(r.Msg.Header[config.MsgPlace])
	if !flag {
		log.Fatalln(config.MsgTypeErr)
	}
	//解析消息
	rpcMsg, err := r.decodeData(conn)
	if err != nil {
		return rpcMsg.Method.Method, err
	}
	//查询方法是否存在
	s, ok := serviceMap[rpcMsg.ServiceName]
	if !ok {
		return rpcMsg.Method.Method, config.ServiceNotFound
	}
	m, ok := s[rpc.MethodName]
	if !ok {
		return rpcMsg.Method.Method, config.MethodIsNotFound
	}
	//返回方法
	return m, nil
}

func (r *RpcClient) Call(conn net.Conn, i ...interface{}) []reflect.Value {
	msg, err := r.decodeData(conn)
	if err != nil {
		log.Fatalf("%s:%s", config.DecodeErr, err)
	}
	flag := IsReqFromCli(msg.Header[config.MsgPlace])
	if flag {
		log.Fatalln(config.MsgTypeErr)
	}
	m, err := r.ReadData(conn)
	if err != nil {
		log.Fatalf("%s:%s", config.ReadDataErr, err)
	}
	args := make([]reflect.Value, len(i))
	for _, p := range i {
		args = append(args, reflect.ValueOf(p))
	}
	resp := m.Call(args)
	return resp
}

func (r *RpcClient) handleConn(conn net.Conn) ([]interface{}, error) {
	//解析服务端发送的消息
	rpcMsg, err := r.decodeData(conn)
	if err != nil {
		return nil, err
	}
	inArgs := make([]interface{}, 0)
	err = ser.Decode(rpcMsg.Method.args, &inArgs)
	return inArgs, err
}

func (r *RpcClient) decodeData(conn net.Conn) (*RpcMsg, error) {
	rpcMsg, err := r.Msg.Read(conn)
	return rpcMsg, err
}

func (r *RpcClient) sendData(serviceName, methodName string, payload interface{}) error {
	var err error
	msg := rpc.NewRpc()
	msg.ServiceName = serviceName
	msg.MethodName = methodName
	msg.Method.args, err = ser.Encode(&payload)
	if err != nil {
		return err
	}
	err = msg.Send(r.Conn)
	return err
}

func (r *RpcClient) readData() {

}
