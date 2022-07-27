package communicate

import (
	"log"
	"my-rpc/config"
	"my-rpc/serialize"
	"net"
	"reflect"
)

type CenterSvr interface {
	Register(string, interface{})
}

type Center struct {
	Listener net.Listener
	*RpcMsg
	*serialize.GobSerialize
}

//map的key为服务名，value为一个map，key为方法名，value为具体方法
var (
	serviceMap map[string]map[string]reflect.Value
	methodMap  map[string]reflect.Value
	ser        serialize.GobSerialize
)

func NewCenter() *Center {
	c := &Center{}
	listener, err := net.Listen("tcp", ":8001")
	c.Listener = listener
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Center) Run() error {
	for {
		conn, err := c.Listener.Accept()
		if err != nil {
			panic(err)
		}
		//处理连接，解析消息
		_, err = c.handleConn(conn)
		if err != nil {
			log.Fatalf("%s:%s", config.HandleErr, err)
		}
	}
}

func isReqFromSvr(msgType byte) bool {
	if msgType == byte(config.FromServer) {
		return true
	}
	return false
}

func (c *Center) Register(serviceName, methodName string, i interface{}) {
	//判断消息类型
	flag := isReqFromSvr(c.Header[config.MsgPlace])
	if !flag {
		log.Fatalln(config.MsgTypeErr)
	}
	methodMap = make(map[string]reflect.Value)
	serviceMap = make(map[string]map[string]reflect.Value)
	methodMap[methodName] = reflect.ValueOf(i).MethodByName(methodName)
	serviceMap[serviceName] = methodMap
}

func IsReqFromCli(msgType byte) bool {
	if msgType == byte(config.FromClient) {
		return true
	}
	return false
}

func (c *Center) handleConn(conn net.Conn) ([]interface{}, error) {
	//解析服务端发送的消息
	rpcMsg, err := c.decodeData(conn)
	if err != nil {
		return nil, err
	}
	inArgs := make([]interface{}, 0)
	err = ser.Decode(rpcMsg.Method.args, &inArgs)
	return inArgs, err
}

func (c *Center) decodeData(conn net.Conn) (*RpcMsg, error) {
	rpcMsg, err := c.Read(conn)
	return rpcMsg, err
}

func (c *Center) SendData(conn net.Conn, resp []byte) error {
	c.RpcMsg.Method.args = resp
	return c.RpcMsg.Send(conn)
}
