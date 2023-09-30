package mistake

import (
	"fmt"
)

type ServerMistake struct {
	err  *ServerMistake
	code ServerMistakeCode
	msg  ServerMistakeMsg
}

func (se *ServerMistake) Error() string {
	return fmt.Sprintf("[ code: %d, message: %s ]", se.code, se.msg)
}

func (se *ServerMistake) Unwrap(err error) error {
	se, ok := err.(*ServerMistake)
	if !ok {
		return nil
	}
	return se.err
}

func NewInvalid() *ServerMistake {
	return &invalidMistake
}

func NewListenFailed() *ServerMistake {
	return &listenFailedMistake
}

func NewServerClose() *ServerMistake {
	return &serverCloseMistake
}

func NewConnExist() *ServerMistake {
	return &connExistMistake
}

func NewConnClose() *ServerMistake {
	return &connCloseMistake
}

func NewWithCode(err *ServerMistake, code ServerMistakeCode) *ServerMistake {
	return &ServerMistake{
		err:  err,
		code: code,
	}
}

func NewWithMsg(err *ServerMistake, msg ServerMistakeMsg) *ServerMistake {
	return &ServerMistake{
		err: err,
		msg: msg,
	}
}
