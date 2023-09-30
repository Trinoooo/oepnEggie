package mistake

type ServerMistakeCode int

const (
	invalidCode      ServerMistakeCode = 10000
	listenFailedCode ServerMistakeCode = 10001
	serverCloseCode  ServerMistakeCode = 10002
	connExistCode    ServerMistakeCode = 10003
	connCloseCode    ServerMistakeCode = 10004
)

type ServerMistakeMsg string

const (
	invalidMsg      ServerMistakeMsg = "param invalid"
	listenFailedMsg ServerMistakeMsg = "listen to host:port failed"
	serverCloseMsg  ServerMistakeMsg = "server close called"
	connExistMsg    ServerMistakeMsg = "connect already exist"
	connCloseMsg    ServerMistakeMsg = "connect close failed"
)

var (
	invalidMistake = ServerMistake{
		code: invalidCode,
		msg:  invalidMsg,
	}
	listenFailedMistake = ServerMistake{
		code: listenFailedCode,
		msg:  listenFailedMsg,
	}
	serverCloseMistake = ServerMistake{
		code: serverCloseCode,
		msg:  serverCloseMsg,
	}
	connExistMistake = ServerMistake{
		code: connExistCode,
		msg:  connExistMsg,
	}
	connCloseMistake = ServerMistake{
		code: connCloseCode,
		msg:  connCloseMsg,
	}
)
