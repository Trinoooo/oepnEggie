package oepnEggie

import (
	"github.com/Trinoooo/oepnEggie/logs"
	"github.com/Trinoooo/oepnEggie/mistake"
	"github.com/Trinoooo/oepnEggie/mux"
	"github.com/Trinoooo/oepnEggie/ptlpsr"
	"github.com/Trinoooo/oepnEggie/trans"
	"net"
	"sync"
)

type Server struct {
	trans.Transporter
	ptlpsr.ProtocolParser
	mux.Multiplexer
	cm     map[*net.Conn]struct{}
	mu     sync.Mutex
	wg     sync.WaitGroup
	doneCh chan struct{}
	errCh  chan error
}

func (srv *Server) Serve() error {
	var err error
	return err
}

// serve
// 处理单个conn
// error 指读写conn出错，与业务错误无关，业务错误应该通过conn给出
func (srv *Server) serve(conn *net.Conn) error {
	err := srv.addConn(conn)
	if err != nil {
		return err
	}

	request, err := srv.ProtocolParser.ParseRequest(conn)
	if err != nil {
		return err
	}

	handler := srv.Multiplexer.Match(request.Method(), request.Path())
	err = handler(request, request)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) addConn(conn *net.Conn) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if _, exist := srv.cm[conn]; exist {
		err := mistake.NewConnExist()
		logs.V1().Println(err)
		return err
	}

	srv.cm[conn] = struct{}{}
	return nil
}

func (srv *Server) Close() error {
	srv.doneCh <- struct{}{}
	for conn := range srv.cm {
		go func(c *net.Conn) {
			err := (*c).Close()
			if err != nil {
				logs.V1().Println(err)
			}
		}(conn)
	}
	srv.wg.Wait()

	err := srv.ProtocolParser.Close()
	if err != nil {
		return err
	}

	err = srv.Multiplexer.Close()
	if err != nil {
		return err
	}

	return <-srv.errCh
}

// ServerBuilder Server的工厂类
type ServerBuilder struct {
	trans.Transporter
	ptlpsr.ProtocolParser
	mux.Multiplexer
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

func (sb *ServerBuilder) WithTransporter(tp trans.Transporter) *ServerBuilder {
	sb.Transporter = tp
	return sb
}

func (sb *ServerBuilder) WithProtocolParser(pp ptlpsr.ProtocolParser) *ServerBuilder {
	sb.ProtocolParser = pp
	return sb
}

func (sb *ServerBuilder) WithMultiplexer(mux mux.Multiplexer) *ServerBuilder {
	sb.Multiplexer = mux
	return sb
}

func (sb *ServerBuilder) Build() (*Server, error) {
	srv := &Server{
		cm:     make(map[*net.Conn]struct{}),
		doneCh: make(chan struct{}),
		errCh:  make(chan error),
	}

	if sb.Transporter == nil {
		srv.Transporter = trans.NewDefaultTransportBuilder().Build()
	}

	if sb.ProtocolParser == nil {
		srv.ProtocolParser = ptlpsr.NewDefaultProtoParserBuilder().Build()
	}

	if sb.Multiplexer == nil {
		srv.Multiplexer = mux.NewDefaultMuxBuilder().Build()
	}

	return srv, nil
}
