package mux

import (
	"github.com/Trinoooo/oepnEggie"
	"github.com/Trinoooo/oepnEggie/logs"
	"github.com/Trinoooo/oepnEggie/mistake"
	"github.com/Trinoooo/oepnEggie/types"
	"sync"
)

// HandlerSet 多路复用处理器集合
type HandlerSet map[types.Method]*Item

type Item struct {
	m    map[string]Handler
	once sync.Once // 懒加载 m
	mu   sync.Mutex
}

type DefaultMux struct {
	srv *oepnEggie.Server
	hs  HandlerSet
}

func (dm *DefaultMux) Get(path string, handler Handler) error {
	return dm.commonRegister(types.Get, path, handler)
}

func (dm *DefaultMux) Post(path string, handler Handler) error {
	return dm.commonRegister(types.Post, path, handler)
}

func (dm *DefaultMux) Match(method types.Method, path string) Handler {
	return nil
}

func (dm *DefaultMux) Close() error {
	return nil
}

func (dm *DefaultMux) commonRegister(method types.Method, path string, handler Handler) error {
	item := dm.hs[method]
	item.once.Do(func() {
		item.m = make(map[string]Handler)
	})
	item.mu.Lock()
	defer item.mu.Unlock()
	if _, exist := item.m[path]; exist {
		err := mistake.NewInvalid()
		logs.V1().Println(err.Error())
		return err
	}
	item.m[path] = handler
	return nil
}

type DefaultMuxBuilder struct {
}

func NewDefaultMuxBuilder() *DefaultMuxBuilder {
	return &DefaultMuxBuilder{}
}

func (dmb *DefaultMuxBuilder) Build() *DefaultMux {
	return nil
}
