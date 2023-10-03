package mux

import "github.com/Trinoooo/oepnEggie/ptlpsr"

// Multiplexer 多路复用器
type Multiplexer interface {
	// 注册handler相关
	Get(path string, handler Handler) error
	Post(path string, handler Handler) error

	// 匹配handler相关
	Match(method ptlpsr.Method, path string) Handler

	Close() error
}

// Handler 请求处理器
type Handler func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error
