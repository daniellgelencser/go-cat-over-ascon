package options_test

import (
	"crypto/tls"
	"testing"

	"go-attested-coap-over-ascon/v3/options"
	"go-attested-coap-over-ascon/v3/tcp"
	"go-attested-coap-over-ascon/v3/tcp/client"
	"go-attested-coap-over-ascon/v3/tcp/server"

	"github.com/stretchr/testify/require"
)

func TestTCPClientApply(t *testing.T) {
	cfg := client.Config{}
	tlsCfg := &tls.Config{}
	opt := []tcp.Option{
		options.WithDisablePeerTCPSignalMessageCSMs(),
		options.WithDisableTCPSignalMessageCSM(),
		options.WithTLS(tlsCfg),
		options.WithConnectionCacheSize(100),
	}
	for _, o := range opt {
		o.TCPClientApply(&cfg)
	}
	require.True(t, cfg.DisablePeerTCPSignalMessageCSMs)
	require.True(t, cfg.DisableTCPSignalMessageCSM)
	require.Equal(t, tlsCfg, cfg.TLSCfg)
	require.Equal(t, uint16(100), cfg.ConnectionCacheSize)
}

func TestTCPServerApply(t *testing.T) {
	cfg := server.Config{}
	opt := []server.Option{
		options.WithDisablePeerTCPSignalMessageCSMs(),
		options.WithDisableTCPSignalMessageCSM(),
		options.WithConnectionCacheSize(100),
	}
	for _, o := range opt {
		o.TCPServerApply(&cfg)
	}
	require.True(t, cfg.DisablePeerTCPSignalMessageCSMs)
	require.True(t, cfg.DisableTCPSignalMessageCSM)
	require.Equal(t, uint16(100), cfg.ConnectionCacheSize)
}
