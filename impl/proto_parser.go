package impl

import (
	"github.com/Trinoooo/oepnEggie/types"
	"net"
)

type DefaultProtoParser struct {
	srv *Server
}

func NewDefaultProtoParser(server *Server) *DefaultProtoParser {
	return &DefaultProtoParser{
		srv: server,
	}
}

func (dpp *DefaultProtoParser) ParseRequest(conn *net.Conn) (types.HttpRequester, error) {
	return nil, nil
}

func (dpp *DefaultProtoParser) Close() error {
	return nil
}
