package ptlpsr

import (
	"net"
)

// ProtocolParser 协议解析器
type ProtocolParser interface {
	ParseRequest(conn *net.Conn) (HttpRequester, error)
	Close() error
}

type HttpRequester interface {
	Method() Method
	Path() string
	Header(key string) string
	Query(key string) string
	Body() []byte
}

type HttpResponser interface {
}

type Method string

const (
	Get  Method = "GET"
	Post Method = "POST"
)
