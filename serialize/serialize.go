package serialize

import (
	"bytes"
	"encoding/gob"
	"my-rpc/config"
	"my-rpc/tool"
)

type Serialize interface {
	Encode(i interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type GobSerialize struct {
}

var check tool.Check

func (g *GobSerialize) Encode(i interface{}) ([]byte, error) {
	flag := check.ForInter(i)
	if !flag {
		return nil, config.InvalidInterface
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(i)
	return buffer.Bytes(), err
}

func (g *GobSerialize) Decode(data []byte, i interface{}) error {
	flag := check.ForInter(i)
	if !flag {
		return config.InvalidInterface
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(i)
}
