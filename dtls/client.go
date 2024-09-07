package dtls

import (
	"fmt"
	"time"

	"go-attested-coap-over-ascon/v3/dtls/server"
	"go-attested-coap-over-ascon/v3/message"
	"go-attested-coap-over-ascon/v3/message/codes"
	"go-attested-coap-over-ascon/v3/message/pool"
	coapNet "go-attested-coap-over-ascon/v3/net"
	"go-attested-coap-over-ascon/v3/net/blockwise"
	"go-attested-coap-over-ascon/v3/net/monitor/inactivity"
	"go-attested-coap-over-ascon/v3/net/responsewriter"
	"go-attested-coap-over-ascon/v3/options"
	"go-attested-coap-over-ascon/v3/udp"
	udpClient "go-attested-coap-over-ascon/v3/udp/client"

	"github.com/pion/dtls/v2"
	dtlsnet "github.com/pion/dtls/v2/pkg/net"
)

var DefaultConfig = func() udpClient.Config {
	cfg := udpClient.DefaultConfig
	cfg.Handler = func(w *responsewriter.ResponseWriter[*udpClient.Conn], r *pool.Message) {
		switch r.Code() {
		case codes.POST, codes.PUT, codes.GET, codes.DELETE:
			if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
				cfg.Errors(fmt.Errorf("dtls client: cannot set response: %w", err))
			}
		}
	}
	return cfg
}()

// Dial creates a client connection to the given target.
func Dial(target string, dtlsCfg *dtls.Config, opts ...udp.Option) (*udpClient.Conn, error) {
	cfg := DefaultConfig
	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}

	c, err := cfg.Dialer.DialContext(cfg.Ctx, cfg.Net, target)
	if err != nil {
		return nil, err
	}

	conn, err := dtls.Client(dtlsnet.PacketConnFromConn(c), c.RemoteAddr(), dtlsCfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts, options.WithCloseSocket())
	return Client(conn, opts...), nil
}

// Client creates client over dtls connection.
func Client(conn *dtls.Conn, opts ...udp.Option) *udpClient.Conn {
	cfg := DefaultConfig
	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}
	if cfg.Errors == nil {
		cfg.Errors = func(error) {
			// default no-op
		}
	}
	if cfg.CreateInactivityMonitor == nil {
		cfg.CreateInactivityMonitor = func() udpClient.InactivityMonitor {
			return inactivity.NewNilMonitor[*udpClient.Conn]()
		}
	}
	if cfg.MessagePool == nil {
		cfg.MessagePool = pool.New(0, 0)
	}
	errorsFunc := cfg.Errors
	cfg.Errors = func(err error) {
		if coapNet.IsCancelOrCloseError(err) {
			// this error was produced by cancellation context or closing connection.
			return
		}
		errorsFunc(fmt.Errorf("dtls: %v: %w", conn.RemoteAddr(), err))
	}

	createBlockWise := func(*udpClient.Conn) *blockwise.BlockWise[*udpClient.Conn] {
		return nil
	}
	if cfg.BlockwiseEnable {
		createBlockWise = func(cc *udpClient.Conn) *blockwise.BlockWise[*udpClient.Conn] {
			v := cc
			return blockwise.New(
				v,
				cfg.BlockwiseTransferTimeout,
				cfg.Errors,
				func(token message.Token) (*pool.Message, bool) {
					return v.GetObservationRequest(token)
				},
			)
		}
	}

	monitor := cfg.CreateInactivityMonitor()
	l := coapNet.NewConn(conn)
	session := server.NewSession(cfg.Ctx,
		l,
		cfg.MaxMessageSize,
		cfg.MTU,
		cfg.CloseSocket,
	)
	cc := udpClient.NewConnWithOpts(session,
		&cfg,
		udpClient.WithBlockWise(createBlockWise),
		udpClient.WithInactivityMonitor(monitor),
		udpClient.WithRequestMonitor(cfg.RequestMonitor),
	)

	cfg.PeriodicRunner(func(now time.Time) bool {
		cc.CheckExpirations(now)
		return cc.Context().Err() == nil
	})

	go func() {
		err := cc.Run()
		if err != nil {
			cfg.Errors(err)
		}
	}()

	return cc
}
