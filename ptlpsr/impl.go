package ptlpsr

import (
	"github.com/Trinoooo/oepnEggie"
	"net"
)

type DefaultProtoParser struct {
	srv *oepnEggie.Server
}

func (dpp *DefaultProtoParser) ParseRequest(conn *net.Conn) (HttpRequester, error) {
	return nil, nil
}

func (dpp *DefaultProtoParser) Close() error {
	return nil
}

type DefaultProtoParserBuilder struct {
}

func NewDefaultProtoParserBuilder() *DefaultProtoParserBuilder {
	return &DefaultProtoParserBuilder{}
}

func (ppb *DefaultProtoParserBuilder) Build() *DefaultProtoParser {
	return nil
}
