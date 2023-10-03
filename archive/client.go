package archive

import (
	"errors"
	"github.com/Trinoooo/oepnEggie/trans"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	logger *log.Logger
}

// 应用层协议格式
// ｜ num ｜ operandA ｜...| operandN | operator |
//
//	1byte	4byte	...		4byte		1byte
func (c *Client) applicationRequest(operands [][]int32, operator byte) {
	cln, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		c.logger.Println("err:", err)
		return
	}
	defer func() {
		if e := syscall.Shutdown(cln, syscall.SHUT_WR); e != nil {
			c.logger.Println("err:", e)
			return
		}
		c.logger.Println("Closed")
	}()

	sockAddr := &syscall.SockaddrInet4{
		Port: trans.parsePort(9999),
		Addr: trans.parseAddr(),
	}
	if err := syscall.Connect(cln, sockAddr); err != nil {
		c.logger.Println("err:", err)
		return
	}
	c.logger.Println("Connected")

	for idx, operand := range operands {
		req := c.buildRequest(operand, operator)
		n, err := syscall.Write(cln, req)
		if err != nil {
			c.logger.Println("err:", err)
			return
		} else if l := len(req); l != n {
			c.logger.Printf("send %d bytes, expect %d bytes\n", n, l)
			return
		}

		l, offset, cycle := 4, 0, 0
		resp := make([]byte, l)
		for l > offset {
			n, err = syscall.Read(cln, resp)
			if err != nil {
				c.logger.Println(err)
				return
			} else if n == 0 {
				c.logger.Println("encounter end of connection")
				break
			}
			offset += n
			cycle++
			c.logger.Printf("[req #%d, cycle #%d] recv %d bytes, total offset %d bytes\n", idx, cycle, n, offset)
		}

		result, err := c.parseResponse(resp[:l])
		if err != nil {
			c.logger.Println("err:", err)
		}
		c.logger.Printf("Result: %d\n", result)
	}
}

func (c *Client) buildRequest(operands []int32, operator byte) []byte {
	operandCount := len(operands)
	req := make([]byte, 0)
	req = append(req, uint8(operandCount))
	for _, operand := range operands {
		req = append(req, uint8(operand>>24), uint8(operand>>16), uint8(operand>>8), uint8(operand))
	}
	req = append(req, operator)
	return req
}

func (c *Client) parseResponse(resp []byte) (int32, error) {
	if len(resp) != 4 {
		return 0, errors.New("invalid resp")
	}

	return ((int32)(resp[0]) << 24) | ((int32)(resp[1]) << 16) | ((int32)(resp[2]) << 8) | (int32)(resp[3]), nil
}

// 应用层协议格式
// ｜ length ｜ filepath ｜
//
//	1byte     ...
func (c *Client) fileRequest(path string) {
	cln, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		c.logger.Println("err:", err)
		return
	}

	sa := &syscall.SockaddrInet4{
		Port: trans.parsePort(9998),
		Addr: trans.parseAddr(),
	}
	if err = syscall.Connect(cln, sa); err != nil {
		c.logger.Println("err:", err)
		return
	}

	request := c.buildFileRequest(path)
	n, err := syscall.Write(cln, request)
	if err != nil {
		c.logger.Println("err:", err)
		return
	} else if n != len(request) {
		c.logger.Println("todo...")
	}

	fd, err := syscall.Open(path+"_copy", syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, syscall.S_IWRITE|syscall.S_IREAD|syscall.S_IWGRP|syscall.S_IRGRP|syscall.S_IWOTH|syscall.S_IROTH)
	if err != nil {
		c.logger.Println("err:", err)
		return
	}

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := syscall.Read(cln, buffer)
			if err != nil {
				c.logger.Println("err:", err)
				return
			} else if n == 0 {
				c.logger.Println("meet EOF, break loop...")
				if err := syscall.Close(fd); err != nil {
					c.logger.Println("err:", err)
					return
				}
				break
			}

			_, err = syscall.Write(fd, buffer[:n])
			if err != nil {
				c.logger.Println("err:", err)
				return
			}
		}

		_, err = syscall.Write(cln, []byte("Thank you!"))
		if err != nil {
			c.logger.Println("err:", err)
			return
		}

		if err := syscall.Shutdown(cln, syscall.SHUT_WR); err != nil {
			c.logger.Println("err:", err)
			return
		}
	}()
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-signalCh:
		if err := syscall.Close(cln); err != nil {
			c.logger.Println("err:", err)
		}
	}
	c.logger.Println("gracefully shutdown")
}

func (c *Client) buildFileRequest(path string) []byte {
	pathBytes := []byte(path)
	length := uint8(len(pathBytes))
	req := make([]byte, 0)
	req = append(req, length)
	return append(req, pathBytes...)
}
