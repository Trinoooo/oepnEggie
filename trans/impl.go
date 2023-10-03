package trans

type DefaultTransport struct {
}

type DefaultTransportBuilder struct {
}

func NewDefaultTransportBuilder() *DefaultTransportBuilder {
	return &DefaultTransportBuilder{}
}

func (dcb *DefaultTransportBuilder) Build() *DefaultTransport {
	return &DefaultTransport{}
}
