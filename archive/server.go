package archive

import (
	"encoding/binary"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	logger *log.Logger
}

func (s *Server) processApplicationReq() {
	srv, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		s.logger.Println("err:", err)
		return
	}
	defer func() {
		if e := syscall.Close(srv); e != nil {
			s.logger.Println("err:", e)
			return
		}
		s.logger.Println("Closed")
	}()

	sockAddr := &syscall.SockaddrInet4{
		Port: parsePort(9999),
		Addr: parseAddr(),
	}
	if err = syscall.Bind(srv, sockAddr); err != nil {
		s.logger.Println("err:", err)
	}

	if err = syscall.Listen(srv, 10); err != nil {
		s.logger.Println("err:", err)
		return
	}

	for {
		in, _, err := syscall.Accept(srv)
		if err != nil {
			s.logger.Println("err:", err)
			return
		}

		go func(fd int) {
			defer func() {
				if e := syscall.Shutdown(in, syscall.SHUT_RDWR); e != nil {
					s.logger.Println("err:", e)
					return
				}
				s.logger.Printf("Server Connection #%d Closed\n", fd)
			}()

			req := make([]byte, 1024)
			length, offset, cycle := 0, 0, 1
			for {
				n, err := syscall.Read(in, req)
				if err != nil {
					s.logger.Println(err)
					return
				} else if n == 0 {
					s.logger.Println("encounter end of connection")
					break
				}

				length = int(req[0])*4 + 2
				offset += n
				s.logger.Printf("[cycle #%d] recv %d bytes, total offset %d bytes\n", cycle, n, offset)
				for length > offset {
					n, err := syscall.Read(in, req[offset:])
					if err != nil {
						s.logger.Println(err)
						return
					} else if n == 0 {
						s.logger.Println("encounter end of connection")
						break
					}
					offset += n
					cycle++
					s.logger.Printf("[cycle #%d] recv %d bytes, total offset %d bytes\n", cycle, n, offset)
				}

				s.logger.Printf("Client send msg: %b\n", req[:length])

				result, err := s.calculate(s.logger, req[:length])
				if err != nil {
					return
				}
				s.logger.Println("Result: ", result)

				resp := make([]byte, 4)
				binary.BigEndian.PutUint32(resp, uint32(result))
				n, err = syscall.Write(in, resp)
				if err != nil {
					s.logger.Println("err:", err)
					return
				} else if l := len(resp); l != n {
					s.logger.Printf("send %d bytes, expect %d bytes\n", n, l)
				}

				// 长链接
				copy(req[:], req[length:offset])
				offset = offset - length
				length, cycle = 0, 1
			}
		}(in)
	}
}

func (s *Server) calculate(logger *log.Logger, req []byte) (int32, error) {
	length := int(req[0])
	operator := req[len(req)-1]
	logger.Printf("operator: %c\n", operator)
	calculator := s.buildCalculator(operator)
	operand := int32(0)
	for i := 0; i < length; i++ {
		offset := i*4 + 1
		operand = (int32(req[offset]) << 24) | (int32(req[offset+1]) << 16) | (int32(req[offset+2]) << 8) | int32(req[offset+3])
		logger.Println("operand:", operand)
		calculator.Cal(operand)
	}
	return calculator.Result(), nil
}

func (s *Server) buildCalculator(operator byte) Calculator {
	switch operator {
	case '+':
		return &Addition{}
	case '-':
	case '*':
	case '\\':
	}
	return nil
}

type Calculator interface {
	Cal(int32)
	Result() int32
}

type Addition struct {
	mu  sync.Mutex
	sum int32
}

func (add *Addition) Cal(operand int32) {
	add.mu.Lock()
	defer add.mu.Unlock()
	add.sum += operand
}

func (add *Addition) Result() int32 {
	return add.sum
}

func (s *Server) processFileReq(readyCh chan struct{}) {
	srvFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		s.logger.Println("err:", err)
		return
	}

	if err = syscall.SetsockoptInt(srvFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		s.logger.Println("err:", err)
		return
	}
	s.logger.Println("allow reuse time-waited addr:port...")

	sa := &syscall.SockaddrInet4{
		Port: parsePort(9998),
		Addr: parseAddr(),
	}
	if err = syscall.Bind(srvFd, sa); err != nil {
		s.logger.Println("err:", err)
		return
	}

	if err := syscall.Listen(srvFd, syscall.SOMAXCONN); err != nil {
		s.logger.Println("err:", err)
		return
	}
	readyCh <- struct{}{}

	go func() {
		wg := sync.WaitGroup{}
		for {
			in, _, err := syscall.Accept(srvFd)
			if err == syscall.ECONNABORTED || err == syscall.EINTR {
				s.logger.Println("accept interrupted. break loop...")
				break
			} else if err != nil {
				s.logger.Println("err:", err)
				continue
			}

			wg.Add(1)
			go func(inFd int) {
				defer func() {
					wg.Done()
					if e := syscall.Close(inFd); e != nil {
						s.logger.Println("err:", e)
					}
				}()
				s.processFile(inFd)
			}(in)
		}
		wg.Wait()
	}()
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-signalCh:
		if err := syscall.Close(srvFd); err != nil {
			s.logger.Println("err:", err)
		}
	}
	s.logger.Println("gracefully shutdown")
}

func (s *Server) processFile(inFd int) {
	req := make([]byte, 1024)
	n, err := syscall.Read(inFd, req)
	if err != nil {
		s.logger.Println("err:", err)
		return
	} else if n == 0 {
		s.logger.Println("peer shutdown WR endpoint...")
		return
	}

	offset, length := n, req[0]
	for int(length) > offset {
		n, err := syscall.Read(inFd, req[offset:])
		if err != nil {
			s.logger.Println("err:", err)
			return
		} else if n == 0 {
			s.logger.Println("peer shutdown WR endpoint...")
			return
		}

		offset += n
	}

	filePath := string(req[1 : length+1])
	s.logger.Println("file path:", filePath)
	fd, err := syscall.Open(filePath, syscall.O_RDONLY, 0)
	if err != nil {
		s.logger.Println("err:", err)
		return
	}

	buffer := make([]byte, 128)
	for {
		fn, err := syscall.Read(fd, buffer)
		if err != nil {
			s.logger.Println("err:", err)
			return
		} else if fn == 0 {
			s.logger.Println("meet EOF, break loop...")
			// 文件传输完毕，关闭socket写端，以便对端可以感知到close信号
			if err := syscall.Shutdown(inFd, syscall.SHUT_WR); err != nil {
				s.logger.Println("err:", err)
			}
			break
		}

		wn, err := syscall.Write(inFd, buffer[:fn])
		if err != nil {
			s.logger.Println("err:", err)
			return
		} else if wn != fn {
			s.logger.Println("todo...")
		}
	}

	// 读取 "Thank you!"
	offset = 0
	for {
		n, err := syscall.Read(inFd, req[offset:])
		if err != nil {
			s.logger.Println("err:", err)
			return
		} else if n == 0 {
			s.logger.Println("peer shutdown WR endpoint...")
			break
		}
		offset += n
	}

	s.logger.Println("recv msg from client:", string(req[:offset]))
	return
}
