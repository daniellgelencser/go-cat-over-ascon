// Package coap provides a CoAP client and server.
package coap

import (
	"errors"
	"fmt"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/ascon"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/dtls"
	dtlsServer "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/dtls/server"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/mux"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/options"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp"
	udpServer "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp/server"

	piondtls "github.com/pion/dtls/v2"
)

// ListenAndServe Starts a server on address and network specified Invoke handler
// for incoming queries.
func ListenAndServe(network string, addr string, handler mux.Handler) (err error) {
	switch network {
	case "udp", "udp4", "udp6", "":
		l, err := net.NewListenUDP(network, addr)
		if err != nil {
			return err
		}
		defer func() {
			if errC := l.Close(); errC != nil && err == nil {
				err = errC
			}
		}()
		s := udp.NewServer(options.WithMux(handler))
		return s.Serve(l)
	default:
		return fmt.Errorf("invalid network (%v)", network)
	}
}

// ListenAndServeDTLS Starts a server on address and network over DTLS specified Invoke handler
// for incoming queries.
func ListenAndServeDTLS(network string, addr string, config *piondtls.Config, handler mux.Handler) (err error) {
	l, err := net.NewDTLSListener(network, addr, config)
	if err != nil {
		return err
	}
	defer func() {
		if errC := l.Close(); errC != nil && err == nil {
			err = errC
		}
	}()
	s := dtls.NewServer(options.WithMux(handler))
	return s.Serve(l)
}

func ListenAndServerASCON(network string, addr string, handler mux.Handler) (err error) {
	l, err := net.NewListenUDP(network, addr)
	if err != nil {
		return err
	}
	defer func() {
		if errC := l.Close(); errC != nil && err == nil {
			err = errC
		}
	}()
	s := ascon.NewServer(options.WithMux(handler))
	return s.Serve(l)
}

// ListenAndServeWithOption Starts a server on address and network specified Invoke options
// for incoming queries. The options is only support tcpServer.Option and udpServer.Option
func ListenAndServeWithOptions(network, addr string, opts ...any) (err error) {
	udpOptions := []udpServer.Option{}
	for _, opt := range opts {
		switch o := opt.(type) {
		case udpServer.Option:
			udpOptions = append(udpOptions, o)
		default:
			return errors.New("only support tcpServer.Option and udpServer.Option")
		}
	}

	switch network {
	case "udp", "udp4", "udp6", "":
		l, err := net.NewListenUDP(network, addr)
		if err != nil {
			return err
		}
		defer func() {
			if errC := l.Close(); errC != nil && err == nil {
				err = errC
			}
		}()
		s := udp.NewServer(udpOptions...)
		return s.Serve(l)
	default:
		return fmt.Errorf("invalid network (%v)", network)
	}
}

// ListenAndServeDTLSWithOptions Starts a server on address and network over DTLS specified Invoke options
// for incoming queries.
func ListenAndServeDTLSWithOptions(network string, addr string, config *piondtls.Config, opts ...dtlsServer.Option) (err error) {
	l, err := net.NewDTLSListener(network, addr, config)
	if err != nil {
		return err
	}
	defer func() {
		if errC := l.Close(); errC != nil && err == nil {
			err = errC
		}
	}()
	s := dtls.NewServer(opts...)
	return s.Serve(l)
}
