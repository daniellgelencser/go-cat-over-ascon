package udp

import "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp/server"

func NewServer(opt ...server.Option) *server.Server {
	return server.New(opt...)
}
