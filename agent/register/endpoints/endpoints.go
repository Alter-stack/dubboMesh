package endpoints


type EndPoint struct {
	host string
	port int
}

func (e *EndPoint) GetHost() string {
	return e.host
}

func (e *EndPoint) GetPort() int {
	return e.port
}

func NewEndPoint(host string, port int) *EndPoint {
	e := EndPoint{
		host:host,
		port:port,
	}
	return &e
}

