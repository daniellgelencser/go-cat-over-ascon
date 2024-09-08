package connection

import (
	"fmt"
	"net"
	"time"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/codes"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/monitor/inactivity"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/responsewriter"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/options/config"
)

type Config struct {
	config.Common[*Conn]
	CreateInactivityMonitor        CreateInactivityMonitorFunc
	Net                            string
	GetMID                         GetMIDFunc
	Handler                        HandlerFunc
	Dialer                         *net.Dialer
	OnNewConn                      OnNewConnFunc
	TransmissionNStart             uint32
	TransmissionAcknowledgeTimeout time.Duration
	TransmissionMaxRetransmit      uint32
	CloseSocket                    bool
	MTU                            uint16
}

func NewConfig(
	createInactivityMonitorFunc CreateInactivityMonitorFunc,
	dialer *net.Dialer,
	net string,
	onNewConnFunc OnNewConnFunc,
	handlerFunc HandlerFunc,
) Config {
	opts := Config{
		Common:                         config.NewCommon[*Conn](),
		CreateInactivityMonitor:        createInactivityMonitorFunc,
		Dialer:                         dialer,
		Net:                            net,
		OnNewConn:                      onNewConnFunc,
		TransmissionNStart:             1,
		TransmissionAcknowledgeTimeout: time.Second * 2,
		TransmissionMaxRetransmit:      4,
		GetMID:                         message.GetMID,
		MTU:                            1472,
		Handler:                        handlerFunc,
	}
	return opts
}

func NewNilInactivityMonitor() InactivityMonitor {
	return inactivity.NewNilMonitor[*Conn]()
}

var DefaultServerConfig = func() Config {
	var options Config
	options = NewConfig(
		func() InactivityMonitor {
			timeout := time.Second * 16
			onInactive := func(cc *Conn) {
				_ = cc.Close()
			}
			return inactivity.New(timeout, onInactive)
		},
		nil,
		"",
		func(cc *Conn) {
			// do nothing by default
			// cc.serverHandshake()
		},
		func(w *responsewriter.ResponseWriter[*Conn], r *pool.Message) {
			if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
				options.Errors(fmt.Errorf("udp server: cannot set response: %w", err))
			}
		},
	)
	return options
}()

var DefaultClientConfig = func() Config {
	var options Config
	options = NewConfig(
		func() InactivityMonitor {
			return inactivity.NewNilMonitor[*Conn]()
		},
		&net.Dialer{Timeout: time.Second * 3},
		"udp",
		nil,
		func(w *responsewriter.ResponseWriter[*Conn], r *pool.Message) {
			switch r.Code() {
			case codes.POST, codes.PUT, codes.GET, codes.DELETE:
				if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
					options.Errors(fmt.Errorf("udp client: cannot set response: %w", err))
				}
			}
		},
	)

	return options
}()
