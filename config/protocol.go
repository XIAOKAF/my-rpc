package config

const (
	HeaderLen   = 5
	MagicNumber = 0x06
	SplitLen    = 4
	Version     = 1.0
	MsgPlace    = 2
)

type MsgType byte

const (
	FromServer MsgType = iota
	FromClient
)

type CompressType byte

const (
	None CompressType = iota
	Gzip
)

type SerializeType byte

const (
	Gob SerializeType = iota
	Json
)
