package ascon

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go-attested-coap-over-ascon/v3/ascon/coder"
	"go-attested-coap-over-ascon/v3/ascon/connection"
	"go-attested-coap-over-ascon/v3/message"
	"go-attested-coap-over-ascon/v3/message/codes"
	"go-attested-coap-over-ascon/v3/message/pool"
	coapNet "go-attested-coap-over-ascon/v3/net"
	"go-attested-coap-over-ascon/v3/net/blockwise"
	"go-attested-coap-over-ascon/v3/net/monitor/inactivity"
	"go-attested-coap-over-ascon/v3/options"
)

// A Option sets options such as credentials, keepalive parameters, etc.
type ClientOption interface {
	ASCONClientApply(cfg *connection.Config)
}

// Dial creates a client connection to the given target.
func Dial(target string, opts ...ClientOption) (*connection.Conn, error) {
	cfg := connection.DefaultClientConfig
	for _, o := range opts {
		o.ASCONClientApply(&cfg)
	}
	c, err := cfg.Dialer.DialContext(cfg.Ctx, cfg.Net, target)
	if err != nil {
		return nil, err
	}
	conn, ok := c.(*net.UDPConn)
	if !ok {
		return nil, fmt.Errorf("unsupported connection type: %T", c)
	}
	opts = append(opts, options.WithCloseSocket())
	return client(conn, opts...), nil
}

// Client creates client over ascon connection.
func client(conn *net.UDPConn, opts ...ClientOption) *connection.Conn {
	cfg := connection.DefaultClientConfig
	for _, o := range opts {
		o.ASCONClientApply(&cfg)
	}
	if cfg.Errors == nil {
		cfg.Errors = func(error) {
			// default no-op
		}
	}
	if cfg.CreateInactivityMonitor == nil {
		cfg.CreateInactivityMonitor = func() connection.InactivityMonitor {
			return inactivity.NewNilMonitor[*connection.Conn]()
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
		errorsFunc(fmt.Errorf("ascon: %v: %w", conn.RemoteAddr(), err))
	}
	addr, _ := conn.RemoteAddr().(*net.UDPAddr)
	createBlockWise := func(cc *connection.Conn) *blockwise.BlockWise[*connection.Conn] {
		return nil
	}
	if cfg.BlockwiseEnable {
		createBlockWise = func(cc *connection.Conn) *blockwise.BlockWise[*connection.Conn] {
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
	l := coapNet.NewUDPConn(cfg.Net, conn, coapNet.WithErrors(cfg.Errors))

	session := connection.NewSession(cfg.Ctx,
		context.Background(),
		l,
		addr,
		cfg.MaxMessageSize,
		cfg.MTU,
		cfg.CloseSocket,
	)

	cc := connection.NewConn(
		session,
		createBlockWise,
		monitor,
		&cfg,
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

	err := handshake(cc)
	if err != nil {
		log.Fatalf("Could not handshake %s", err)
	}

	return cc
}

func handshake(cc *connection.Conn) error {
	fmt.Println("\nClient Handshake")

	// return nil
	clientPrivateKey := coder.RandomBytes(32)
	clientPublicKey := coder.ComputePublicKey(clientPrivateKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//send client hello
	request := cc.AcquireMessage(ctx)
	token, err := cc.Client.GetToken()
	if err != nil {
		return fmt.Errorf("cannot get token: %w", err)
	}

	request.SetCode(codes.HANDSHAKE)
	request.SetToken(token)
	request.SetBody(bytes.NewReader(clientPublicKey))

	defer cc.ReleaseMessage(request)

	response, err := cc.LimitParallelRequests.Do(request)
	if err != nil {
		return fmt.Errorf("cannot send client hello: %w", err)
	}

	serverPublicKey, err := response.ReadBody()
	if err != nil {
		return fmt.Errorf("cannot read server hello %w", err)
	}

	if len(serverPublicKey) != 32 {
		return fmt.Errorf("invalid server hello")
	}

	sharedKey, err := coder.DeriveSharedKey(clientPrivateKey, serverPublicKey)
	if err != nil {
		return fmt.Errorf("Could not compute shared key: %w\n", err)
	}

	fmt.Printf("Client public: %X\n", clientPublicKey)
	fmt.Printf("Server public: %X\n", serverPublicKey)
	fmt.Printf("Shared key: %X\n", sharedKey)

	// save shared secret
	cc.SetClientSecret(sharedKey[:16])

	//send client ack (non-confirmable)

	fmt.Println("Finish Handshake")
	return nil
}
