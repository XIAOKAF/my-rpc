package tool

import (
	"my-rpc/config"
)

type Header [config.HeaderLen]byte

func (h *Header) CheckMagic() bool {
	return h[0] == config.MagicNumber
}

func (h *Header) CheckCompress() bool {
	return h[3] == byte(config.None) || h[3] == byte(config.Gzip)
}

func (h *Header) SetVersion() {
	h[1] = config.Version
}

func (h *Header) CheckProtocolVersion() bool {
	return h[2] == config.Version
}

func (h *Header) CheckMsgType() bool {
	return h[2] == byte(config.FromServer) || h[2] == byte(config.FromClient)
}

func (h *Header) CheckSerialize() bool {
	return h[4] == byte(config.Gob) || h[4] == byte(config.Json)
}
