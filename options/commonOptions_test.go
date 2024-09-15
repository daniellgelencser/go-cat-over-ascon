package options_test

import (
	"context"
	"net"
	"testing"
	"time"

	dtlsServer "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/dtls/server"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/mux"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/blockwise"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/responsewriter"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/options"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/options/config"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/pkg/runner/periodic"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp"
	udpClient "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp/client"
	udpServer "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/udp/server"

	"github.com/stretchr/testify/require"
)

func TestCommonUDPServerApply(t *testing.T) {
	cfg := udpServer.Config{}
	handler := func(*responsewriter.ResponseWriter[*udpClient.Conn], *pool.Message) {
		// no-op
	}
	ctx := context.Background()
	errs := func(error) {
		// no-op
	}
	processRecvMessage := func(*pool.Message, *udpClient.Conn, config.HandlerFunc[*udpClient.Conn]) {
		// no-op
	}
	inactivityMonitor := func(*udpClient.Conn) {
		// no-op
	}
	periodicRunner := periodic.New(ctx.Done(), time.Millisecond*10)
	onNewConn := func(*udpClient.Conn) {
		// no-op
	}
	requestMonitor := func(*udpClient.Conn, *pool.Message) (bool, error) {
		return false, nil
	}
	mp := pool.New(1024, 1600)
	getToken := func() (message.Token, error) {
		return nil, nil
	}
	opts := []udpServer.Option{
		options.WithHandlerFunc(handler),
		options.WithContext(ctx),
		options.WithMaxMessageSize(1024),
		options.WithErrors(errs),
		options.WithProcessReceivedMessageFunc(processRecvMessage),
		options.WithInactivityMonitor(time.Minute, inactivityMonitor),
		options.WithPeriodicRunner(periodicRunner),
		options.WithBlockwise(true, blockwise.SZX16, time.Second),
		options.WithOnNewConn(onNewConn),
		options.WithRequestMonitor(requestMonitor),
		options.WithMessagePool(mp),
		options.WithGetToken(getToken),
		options.WithLimitClientParallelRequest(42),
		options.WithLimitClientEndpointParallelRequest(43),
		options.WithReceivedMessageQueueSize(10),
	}

	for _, o := range opts {
		o.UDPServerApply(&cfg)
	}
	// WithHandlerFunc
	require.NotNil(t, cfg.Handler)
	// WithContext
	require.Equal(t, ctx, cfg.Ctx)
	// WithMaxMessageSize
	require.Equal(t, uint32(1024), cfg.MaxMessageSize)
	// WithErrors
	require.NotNil(t, cfg.Errors)
	// WithProcessReceivedMessageFunc
	require.NotNil(t, cfg.ProcessReceivedMessage)
	// WithInactivityMonitor
	require.NotNil(t, cfg.CreateInactivityMonitor)
	// WithPeriodicRunner
	require.NotNil(t, cfg.PeriodicRunner)
	// WithBlockwise
	require.True(t, cfg.BlockwiseEnable)
	require.Equal(t, blockwise.SZX16, cfg.BlockwiseSZX)
	require.Equal(t, time.Second, cfg.BlockwiseTransferTimeout)
	// WithOnNewConn
	require.NotNil(t, cfg.OnNewConn)
	// WithRequestMonitor
	require.NotNil(t, cfg.RequestMonitor)
	// WithMessagePool
	require.Equal(t, mp, cfg.MessagePool)
	// WithGetToken
	require.NotNil(t, cfg.GetToken)
	// WithLimitClientParallelRequest
	require.Equal(t, int64(42), cfg.LimitClientParallelRequests)
	// WithLimitClientEndpointParallelRequest
	require.Equal(t, int64(43), cfg.LimitClientEndpointParallelRequests)
	// WithReceivedMessageQueueSize
	require.Equal(t, 10, cfg.ReceivedMessageQueueSize)

	m := mux.NewRouter()
	keepAlive := func(*udpClient.Conn) {
		// no-op
	}
	cfg = udpServer.Config{}
	opts = []udpServer.Option{
		options.WithMux(m),
		options.WithKeepAlive(16, time.Second, keepAlive),
	}
	for _, o := range opts {
		o.UDPServerApply(&cfg)
	}
	// WithMux
	require.NotNil(t, cfg.Handler)
	// WithKeepAlive
	require.NotNil(t, cfg.CreateInactivityMonitor)
}

func TestCommonDTLSServerApply(t *testing.T) {
	cfg := dtlsServer.Config{}
	handler := func(*responsewriter.ResponseWriter[*udpClient.Conn], *pool.Message) {
		// no-op
	}
	ctx := context.Background()
	errs := func(error) {
		// no-op
	}
	processRecvMessage := func(*pool.Message, *udpClient.Conn, config.HandlerFunc[*udpClient.Conn]) {
		// no-op
	}
	inactivityMonitor := func(*udpClient.Conn) {
		// no-op
	}
	periodicRunner := periodic.New(ctx.Done(), time.Millisecond*10)
	onNewConn := func(*udpClient.Conn) {
		// no-op
	}
	requestMonitor := func(*udpClient.Conn, *pool.Message) (bool, error) {
		return false, nil
	}
	mp := pool.New(1024, 1600)
	getToken := func() (message.Token, error) {
		return nil, nil
	}
	opts := []dtlsServer.Option{
		options.WithHandlerFunc(handler),
		options.WithContext(ctx),
		options.WithMaxMessageSize(1024),
		options.WithErrors(errs),
		options.WithProcessReceivedMessageFunc(processRecvMessage),
		options.WithInactivityMonitor(time.Minute, inactivityMonitor),
		options.WithPeriodicRunner(periodicRunner),
		options.WithBlockwise(true, blockwise.SZX16, time.Second),
		options.WithOnNewConn(onNewConn),
		options.WithRequestMonitor(requestMonitor),
		options.WithMessagePool(mp),
		options.WithGetToken(getToken),
		options.WithLimitClientParallelRequest(42),
		options.WithLimitClientEndpointParallelRequest(43),
		options.WithReceivedMessageQueueSize(10),
	}

	for _, o := range opts {
		o.DTLSServerApply(&cfg)
	}
	// WithHandlerFunc
	require.NotNil(t, cfg.Handler)
	// WithContext
	require.Equal(t, ctx, cfg.Ctx)
	// WithMaxMessageSize
	require.Equal(t, uint32(1024), cfg.MaxMessageSize)
	// WithErrors
	require.NotNil(t, cfg.Errors)
	// WithProcessReceivedMessageFunc
	require.NotNil(t, cfg.ProcessReceivedMessage)
	// WithInactivityMonitor
	require.NotNil(t, cfg.CreateInactivityMonitor)
	// WithPeriodicRunner
	require.NotNil(t, cfg.PeriodicRunner)
	// WithBlockwise
	require.True(t, cfg.BlockwiseEnable)
	require.Equal(t, blockwise.SZX16, cfg.BlockwiseSZX)
	require.Equal(t, time.Second, cfg.BlockwiseTransferTimeout)
	// WithOnNewConn
	require.NotNil(t, cfg.OnNewConn)
	// WithRequestMonitor
	require.NotNil(t, cfg.RequestMonitor)
	// WithMessagePool
	require.Equal(t, mp, cfg.MessagePool)
	// WithGetToken
	require.NotNil(t, cfg.GetToken)
	// WithLimitClientParallelRequest
	require.Equal(t, int64(42), cfg.LimitClientParallelRequests)
	// WithLimitClientEndpointParallelRequest
	require.Equal(t, int64(43), cfg.LimitClientEndpointParallelRequests)
	// WithReceivedMessageQueueSize
	require.Equal(t, 10, cfg.ReceivedMessageQueueSize)

	m := mux.NewRouter()
	keepAlive := func(*udpClient.Conn) {
		// no-op
	}
	cfg = dtlsServer.Config{}
	opts = []dtlsServer.Option{
		options.WithMux(m),
		options.WithKeepAlive(16, time.Second, keepAlive),
	}
	for _, o := range opts {
		o.DTLSServerApply(&cfg)
	}
	// WithMux
	require.NotNil(t, cfg.Handler)
	// WithKeepAlive
	require.NotNil(t, cfg.CreateInactivityMonitor)
}

func TestCommonUDPClientApply(t *testing.T) {
	cfg := udpClient.Config{}
	handler := func(*responsewriter.ResponseWriter[*udpClient.Conn], *pool.Message) {
		// no-op
	}
	ctx := context.Background()
	errs := func(error) {
		// no-op
	}
	processRecvMessage := func(*pool.Message, *udpClient.Conn, config.HandlerFunc[*udpClient.Conn]) {
		// no-op
	}
	inactivityMonitor := func(*udpClient.Conn) {
		// no-op
	}
	network := "udp4"
	periodicRunner := periodic.New(ctx.Done(), time.Millisecond*10)
	dialer := &net.Dialer{Timeout: time.Second * 3}

	mp := pool.New(1024, 1600)
	getToken := func() (message.Token, error) {
		return nil, nil
	}
	opts := []udp.Option{
		options.WithHandlerFunc(handler),
		options.WithContext(ctx),
		options.WithMaxMessageSize(1024),
		options.WithErrors(errs),
		options.WithProcessReceivedMessageFunc(processRecvMessage),
		options.WithInactivityMonitor(time.Minute, inactivityMonitor),
		options.WithNetwork(network),
		options.WithPeriodicRunner(periodicRunner),
		options.WithBlockwise(true, blockwise.SZX16, time.Second),
		options.WithCloseSocket(),
		options.WithDialer(dialer),
		options.WithMessagePool(mp),
		options.WithGetToken(getToken),
		options.WithLimitClientParallelRequest(42),
		options.WithLimitClientEndpointParallelRequest(43),
		options.WithReceivedMessageQueueSize(10),
	}

	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}
	// WithHandlerFunc
	require.NotNil(t, cfg.Handler)
	// WithContext
	require.Equal(t, ctx, cfg.Ctx)
	// WithMaxMessageSize
	require.Equal(t, uint32(1024), cfg.MaxMessageSize)
	// WithErrors
	require.NotNil(t, cfg.Errors)
	// WithProcessReceivedMessageFunc
	require.NotNil(t, cfg.ProcessReceivedMessage)
	// WithInactivityMonitor
	require.NotNil(t, cfg.CreateInactivityMonitor)
	// WithNetwork
	require.Equal(t, network, cfg.Net)
	// WithPeriodicRunner
	require.NotNil(t, cfg.PeriodicRunner)
	// WithBlockwise
	require.True(t, cfg.BlockwiseEnable)
	require.Equal(t, blockwise.SZX16, cfg.BlockwiseSZX)
	require.Equal(t, time.Second, cfg.BlockwiseTransferTimeout)
	// WithCloseSocket
	require.True(t, cfg.CloseSocket)
	// WithDialer
	require.Equal(t, dialer, cfg.Dialer)
	// WithMessagePool
	require.Equal(t, mp, cfg.MessagePool)
	// WithGetToken
	require.NotNil(t, cfg.GetToken)
	// WithLimitClientParallelRequest
	require.Equal(t, int64(42), cfg.LimitClientParallelRequests)
	// WithLimitClientEndpointParallelRequest
	require.Equal(t, int64(43), cfg.LimitClientEndpointParallelRequests)
	// WithReceivedMessageQueueSize
	require.Equal(t, 10, cfg.ReceivedMessageQueueSize)

	m := mux.NewRouter()
	keepAlive := func(*udpClient.Conn) {
		// no-op
	}
	cfg = udpClient.Config{}
	opts = []udp.Option{
		options.WithMux(m),
		options.WithKeepAlive(16, time.Second, keepAlive),
	}
	for _, o := range opts {
		o.UDPClientApply(&cfg)
	}
	// WithMux
	require.NotNil(t, cfg.Handler)
	// WithKeepAlive
	require.NotNil(t, cfg.CreateInactivityMonitor)
}
