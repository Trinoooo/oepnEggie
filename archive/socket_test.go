package archive

import (
	"crypto"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	server                      *Server
	client                      *Client
	applicationReqServerReadyCh chan struct{}
	fileReqServerReadyCh        chan struct{}
)

func TestMain(m *testing.M) {
	server = &Server{
		logger: log.New(os.Stdout, "[Server]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix),
	}
	client = &Client{
		logger: log.New(os.Stdout, "[Client]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix),
	}
	applicationReqServerReadyCh = make(chan struct{})
	fileReqServerReadyCh = make(chan struct{})
	m.Run()
}

func TestServer_ProcessApplicationReq(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	t.Parallel()
	server.processApplicationReq()
}

func TestClient_ApplicationRequest(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var additionOperandTestTable = [][]int32{
		{1, 2, 3}, {4, 5},
	}
	t.Parallel()
	select {
	case <-applicationReqServerReadyCh:
	}
	client.applicationRequest(additionOperandTestTable, '+')
}

func TestServer_ProcessFileReq(t *testing.T) {
	t.Parallel()
	server.processFileReq(fileReqServerReadyCh)
}

func TestClient_FileRequest(t *testing.T) {
	abs, err := filepath.Abs("socket_test")
	if err != nil {
		client.logger.Println("err:", err)
		t.Fail()
	}
	testTable := []string{
		abs,
	}
	t.Parallel()
	select {
	case <-fileReqServerReadyCh:
		client.logger.Println("client started...")
	}
	for _, c := range testTable {
		client.fileRequest(c)
	}
}

func TestFileSame(t *testing.T) {
	hash := crypto.SHA256.New()
	origin, _ := os.Open("trans/socket_test")
	cp, _ := os.Open("trans/socket_test_copy")
	io.Copy(hash, origin)
	originSum := fmt.Sprintf("%x", hash.Sum(nil))
	hash.Reset()
	io.Copy(hash, cp)
	cpSum := fmt.Sprintf("%x", hash.Sum(nil))
	if originSum != cpSum {
		t.Errorf("expect cpsum: %s, got: %s", originSum, cpSum)
	}
}

func TestSequence(t *testing.T) {
	t.Run("parallel", func(t *testing.T) {
		t.Run("server", TestServer_ProcessFileReq)
		t.Run("client", TestClient_FileRequest)
	})

	t.Run("hash", TestFileSame)
}
