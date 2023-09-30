package impl

import (
	"fmt"
	"github.com/Trinoooo/oepnEggie/logs"
	"github.com/Trinoooo/oepnEggie/mistake"
	"github.com/Trinoooo/oepnEggie/types"
	"net"
	"sync"
)

type Server struct {
	types.ProtocolParser
	types.Multiplexer
	ln     net.Listener
	cm     map[*net.Conn]struct{}
	mu     sync.Mutex
	wg     sync.WaitGroup
	doneCh chan struct{}
	errCh  chan error
}

func (srv *Server) Serve() error {
	var err error
	for {
		conn, e := srv.ln.Accept()
		if e != nil {
			select {
			case <-srv.doneCh:
				err = mistake.NewServerClose()
			default:
				// 出错不阻塞，继续接受下一个请求
				logs.V1().Println(e)
				continue
			}
		}

		srv.wg.Add(1)
		go func() {
			defer srv.wg.Done()
			if err := srv.serve(&conn); err != nil {
				logs.V1().Println(err)
			}
		}()

		// server.Close() 或者 单个conn执行出错
		if err != nil {
			srv.errCh <- err
			break
		}
	}
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
		err := (*conn).Close()
		if err != nil {
			logs.V1().Println(err)
		}
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
	types.ProtocolParser
	types.Multiplexer
	host string
	port string
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		host: "localhost", // 监听host默认本地
		port: "80",        // 监听端口默认80
	}
}

func (sb *ServerBuilder) WithHost(host string) *ServerBuilder {
	// 非空才赋值，否则用默认的
	if host != "" {
		sb.host = host
	}
	return sb
}

func (sb *ServerBuilder) WithPort(port string) *ServerBuilder {
	// 非空才赋值，否则用默认的
	if port != "" {
		sb.port = port
	}
	return sb
}

func (sb *ServerBuilder) WithProtocolParser(pp types.ProtocolParser) *ServerBuilder {
	// 非空才赋值，否则用默认的
	if pp != nil {
		sb.ProtocolParser = pp
	}
	return sb
}

func (sb *ServerBuilder) WithMultiplexer(mux types.Multiplexer) *ServerBuilder {
	// 非空才赋值，否则用默认的
	if mux != nil {
		sb.Multiplexer = mux
	}
	return sb
}

func (sb *ServerBuilder) Build() (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", sb.host, sb.port))
	if err != nil {
		logs.V1().Println(err)
		return nil, mistake.NewListenFailed()
	}

	srv := &Server{
		cm:     make(map[*net.Conn]struct{}),
		doneCh: make(chan struct{}),
		errCh:  make(chan error),
		ln:     listener,
	}

	if sb.ProtocolParser == nil {
		srv.ProtocolParser = NewDefaultProtoParser(srv)
	}

	if sb.Multiplexer == nil {
		srv.Multiplexer = NewDefaultMux(srv)
	}

	return srv, nil
}
