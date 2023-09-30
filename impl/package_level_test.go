package impl

import "testing"

var server *Server

func TestMain(m *testing.M) {
	var err error
	server, err = NewServerBuilder().Build()
	if err != nil {
		panic(err)
	}
	m.Run()
}
