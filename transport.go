package mbserver

type transport interface {
	Close() error
	ExecuteRequest(*pdu) (*pdu, error)
	// ReadRequest() (*pdu, error)
	HandleRequest([]byte) (*pdu, error)
	WriteResponse(*pdu) error
}
