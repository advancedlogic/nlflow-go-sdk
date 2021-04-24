package core

type Handler struct {
	//logger    *zap.Logger
	endpoint     string
	method       string
	transportIn  TransportIn
	transportOut TransportOut
	logic        Logic

	Error error
}

func NewHandler() *Handler {
	return &Handler{
		method:       "GET",
		Error:        nil,
		transportIn:  &defaultTransportIn{},
		transportOut: &defaultTransportOut{},
	}
}

//func (handler *Handler) WithLogger(logger *zap.Logger) *Handler {
//	handler.logger = logger
//	return handler
//}

func (handler *Handler) WithEndpoint(endpoint string) *Handler {
	handler.endpoint = endpoint
	return handler
}

func (handler *Handler) WithMethod(method string) *Handler {
	handler.method = method
	return handler
}

func (handler *Handler) WithTransportIn(transport TransportIn) *Handler {
	handler.transportIn = transport
	return handler
}

func (handler *Handler) WithTransportOut(transport TransportOut) *Handler {
	handler.transportOut = transport
	return handler
}

func (handler *Handler) WithLogic(logic Logic) *Handler {
	handler.logic = logic
	return handler
}

func (handler *Handler) GetTransportIn() TransportIn {
	return handler.transportIn
}

func (handler *Handler) GetTransportOut() TransportOut {
	return handler.transportOut
}

func (handler *Handler) GetLogic() Logic {
	return handler.logic
}

func (handler *Handler) GetMethod() string {
	return handler.method
}

func (handler *Handler) GetEndpoint() string {
	return handler.endpoint
}
