package connection

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/ascon/coder"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"
	coapNet "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net"
)

type ISession interface {
	Context() context.Context
	Close() error
	MaxMessageSize() uint32
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	// NetConn returns the underlying connection that is wrapped by Session. The Conn returned is shared by all invocations of NetConn, so do not modify it.
	NetConn() net.Conn
	WriteMessage(req *pool.Message) error
	// WriteMulticast sends multicast to the remote multicast address.
	// By default it is sent over all network interfaces and all compatible source IP addresses with hop limit 1.
	// Via opts you can specify the network interface, source IP address, and hop limit.
	WriteMulticastMessage(req *pool.Message, address *net.UDPAddr, opts ...coapNet.MulticastOption) error
	Run(cc *Conn) error
	AddOnClose(f EventFunc)
	SetContextValue(key interface{}, val interface{})
	Done() <-chan struct{}
	SetServerSecret(secret []byte)
}

type Session struct {
	onClose []EventFunc

	ctx atomic.Value // TODO: change to atomic.Pointer[context.Context] for go1.19

	doneCtx    context.Context
	connection *coapNet.UDPConn
	doneCancel context.CancelFunc

	cancel context.CancelFunc
	raddr  *net.UDPAddr

	mutex          sync.Mutex
	maxMessageSize uint32
	mtu            uint16

	serverSecret  []byte
	isSecretReady bool

	closeSocket bool
}

func NewSession(
	ctx context.Context,
	doneCtx context.Context,
	connection *coapNet.UDPConn,
	raddr *net.UDPAddr,
	maxMessageSize uint32,
	mtu uint16,
	closeSocket bool,
) *Session {
	ctx, cancel := context.WithCancel(ctx)

	doneCtx, doneCancel := context.WithCancel(doneCtx)
	s := &Session{
		cancel:         cancel,
		connection:     connection,
		raddr:          raddr,
		maxMessageSize: maxMessageSize,
		mtu:            mtu,
		closeSocket:    closeSocket,
		doneCtx:        doneCtx,
		doneCancel:     doneCancel,
		isSecretReady:  false,
	}
	s.ctx.Store(&ctx)
	return s
}

func (s *Session) SetServerSecret(secret []byte) {
	s.serverSecret = secret
}

func (s *Session) popOnClose() []EventFunc {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tmp := s.onClose
	s.onClose = nil
	return tmp
}

func (s *Session) Shutdown() {
	defer s.doneCancel()
	for _, f := range s.popOnClose() {
		f()
	}
}

func (s *Session) AddOnClose(f EventFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onClose = append(s.onClose, f)
}

func (s *Session) Close() error {
	s.cancel()
	if s.closeSocket {
		return s.connection.Close()
	}
	return nil
}

func (s *Session) Context() context.Context {
	return *s.ctx.Load().(*context.Context) //nolint:forcetypeassert
}

// Done signalizes that connection is not more processed.
func (s *Session) Done() <-chan struct{} {
	return s.doneCtx.Done()
}

func (s *Session) LocalAddr() net.Addr {
	return s.connection.LocalAddr()
}

func (s *Session) MaxMessageSize() uint32 {
	return s.maxMessageSize
}

// NetConn returns the underlying connection that is wrapped by s. The Conn returned is shared by all invocations of NetConn, so do not modify it.
func (s *Session) NetConn() net.Conn {
	return s.connection.NetConn()
}

func (s *Session) RemoteAddr() net.Addr {
	return s.raddr
}

func (s *Session) Run(cc *Conn) (err error) {
	defer func() {
		err1 := s.Close()
		if err == nil {
			err = err1
		}
		s.Shutdown()
	}()
	m := make([]byte, s.mtu)
	for {
		buf := m
		n, _, err := s.connection.ReadWithContext(s.Context(), buf)
		if err != nil {
			return err
		}
		buf = buf[:n]
		err = cc.Process(buf)
		if err != nil {
			return err
		}
	}
}

// SetContextValue stores the value associated with key to context of connection.
func (s *Session) SetContextValue(key interface{}, val interface{}) {
	ctx := context.WithValue(s.Context(), key, val)
	s.ctx.Store(&ctx)
}

func (s *Session) resolveCoder() *coder.Coder {
	cdr := coder.DefaultCoder
	if s.serverSecret == nil {
		return cdr
	}

	if s.isSecretReady {
		cdr.SetSecret(s.serverSecret)
	} else {
		s.isSecretReady = true
	}

	return cdr
}

func (s *Session) WriteMessage(req *pool.Message) error {
	data, err := req.MarshalWithEncoder(s.resolveCoder())
	if err != nil {
		return fmt.Errorf("cannot marshal: %w", err)
	}
	return s.connection.WriteWithContext(req.Context(), s.raddr, data)
}

// WriteMulticastMessage sends multicast to the remote multicast address.
// By default it is sent over all network interfaces and all compatible source IP addresses with hop limit 1.
// Via opts you can specify the network interface, source IP address, and hop limit.
func (s *Session) WriteMulticastMessage(req *pool.Message, address *net.UDPAddr, opts ...coapNet.MulticastOption) error {
	fmt.Printf("WriteMulticastMessage: ")
	data, err := req.MarshalWithEncoder(s.resolveCoder())
	if err != nil {
		return fmt.Errorf("cannot marshal: %w", err)
	}
	return s.connection.WriteMulticast(req.Context(), address, data, opts...)
}
