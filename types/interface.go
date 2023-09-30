package types

import (
	"github.com/Trinoooo/oepnEggie/impl"
	"net"
)

// ProtocolParser 协议解析器
type ProtocolParser interface {
	Close() error
	ParseRequest(conn *net.Conn) (HttpRequester, error)
}

// Multiplexer 多路复用器
type Multiplexer interface {
	// 注册handler相关
	Get(path string, handler impl.Handler) error
	Post(path string, handler impl.Handler) error

	// 匹配handler相关
	Match(method Method, path string) impl.Handler

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
