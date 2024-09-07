package server

import (
	"fmt"
	"time"

	"go-attested-coap-over-ascon/v3/message"
	"go-attested-coap-over-ascon/v3/message/codes"
	"go-attested-coap-over-ascon/v3/message/pool"
	"go-attested-coap-over-ascon/v3/net/monitor/inactivity"
	"go-attested-coap-over-ascon/v3/net/responsewriter"
	"go-attested-coap-over-ascon/v3/options/config"
	udpClient "go-attested-coap-over-ascon/v3/udp/client"
)

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as COAP handlers.
type HandlerFunc = func(*responsewriter.ResponseWriter[*udpClient.Conn], *pool.Message)

type ErrorFunc = func(error)

// OnNewConnFunc is the callback for new connections.
type OnNewConnFunc = func(cc *udpClient.Conn)

type GetMIDFunc = func() int32

var DefaultConfig = func() Config {
	opts := Config{
		Common: config.NewCommon[*udpClient.Conn](),
		CreateInactivityMonitor: func() udpClient.InactivityMonitor {
			timeout := time.Second * 16
			onInactive := func(cc *udpClient.Conn) {
				_ = cc.Close()
			}
			return inactivity.New(timeout, onInactive)
		},
		OnNewConn: func(*udpClient.Conn) {
			// do nothing by default
		},
		RequestMonitor: func(*udpClient.Conn, *pool.Message) (bool, error) {
			return false, nil
		},
		TransmissionNStart:             1,
		TransmissionAcknowledgeTimeout: time.Second * 2,
		TransmissionMaxRetransmit:      4,
		GetMID:                         message.GetMID,
		MTU:                            udpClient.DefaultMTU,
	}
	opts.Handler = func(w *responsewriter.ResponseWriter[*udpClient.Conn], _ *pool.Message) {
		if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
			opts.Errors(fmt.Errorf("udp server: cannot set response: %w", err))
		}
	}
	return opts
}()

type Config struct {
	config.Common[*udpClient.Conn]
	CreateInactivityMonitor        udpClient.CreateInactivityMonitorFunc
	GetMID                         GetMIDFunc
	Handler                        HandlerFunc
	OnNewConn                      OnNewConnFunc
	RequestMonitor                 udpClient.RequestMonitorFunc
	TransmissionNStart             uint32
	TransmissionAcknowledgeTimeout time.Duration
	TransmissionMaxRetransmit      uint32
	MTU                            uint16
}
