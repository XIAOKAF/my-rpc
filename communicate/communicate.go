package communicate

import (
	"encoding/binary"
	"fmt"
	"io"
	"my-rpc/config"
	"my-rpc/tool"
	"reflect"
	"unsafe"
)

type RpcMsg struct {
	*tool.Header
	ServiceName string
	MethodName  string
	Method      Methods
}

type Methods struct {
	Method reflect.Value
	args   []byte
}

func (r *RpcMsg) NewRpc() *RpcMsg {
	header := tool.Header([config.HeaderLen]byte{})
	header[0] = config.MagicNumber
	header[1] = config.Version
	header[2] = byte(config.FromServer)
	header[3] = byte(config.Gob)
	return &RpcMsg{
		Header: &header,
	}
}

func (r *RpcMsg) Send(writer io.Writer) error {
	_, err := writer.Write(r.Header[:])
	if err != nil {
		return err
	}
	dataLen := config.SplitLen*4 + len(r.Method.args) + len(r.MethodName) + len(r.ServiceName)
	err = binary.Write(writer, binary.BigEndian, uint32(dataLen))
	if err != nil {
		return err
	}
	//写入信息的长度方便解析
	//写入信息
	err = r.write(writer, uint32(len(r.ServiceName)), r.ServiceName)
	if err != nil {
		return err
	}
	err = r.write(writer, uint32(len(r.MethodName)), r.MethodName)
	if err != nil {
		return err
	}
	err = r.write(writer, uint32(len(r.Method.args)), r.Method.args)
	return err
}

func (r *RpcMsg) write(writer io.Writer, length uint32, i interface{}) error {
	err := binary.Write(writer, binary.BigEndian, length)
	if err != nil {
		return err
	}
	byte, err := ser.Encode(i)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, byte)
	return err
}

func (r *RpcMsg) Read(reader io.Reader) (*RpcMsg, error) {
	//解析头
	dataByte, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("解析头失败：%s", err)
		return nil, err
	}
	//header
	headerByte := make([]byte, config.HeaderLen)
	headerByte = dataByte[:config.HeaderLen]
	_, err = io.ReadFull(reader, headerByte)
	if err != nil {
		return nil, err
	}
	bodyLen := len(dataByte) - config.HeaderLen
	//body
	data := make([]byte, bodyLen)
	data = dataByte[config.HeaderLen:]
	//服务名
	start := 0
	end := start + config.SplitLen
	serviceLen := binary.BigEndian.Uint32(data[start:end])
	start = end
	end = start + int(serviceLen)
	rpc.ServiceName = rpc.byteToString(data[start:end])
	//方法名
	start = end
	end = start + config.SplitLen
	methodLen := binary.BigEndian.Uint32(data[start:end])
	start = end
	end = start + int(methodLen)
	rpc.MethodName = rpc.byteToString(data[start:end])
	//传参
	start = end
	end = start + config.SplitLen
	payloadLen := binary.BigEndian.Uint32(data[start:end])
	start = end
	end = start + int(payloadLen)
	rpc.Method.args = data[start:end]
	return &rpc, nil
}

func (r *RpcMsg) byteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
