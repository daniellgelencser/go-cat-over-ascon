package connection

import (
	"sync"
	"time"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"

	"go.uber.org/atomic"
)

type midElement struct {
	handler    HandlerFunc
	start      time.Time
	deadline   time.Time
	retransmit atomic.Int32

	private struct {
		sync.Mutex
		msg *pool.Message
	}
}

func (m *midElement) ReleaseMessage(cc *Conn) {
	m.private.Lock()
	defer m.private.Unlock()
	if m.private.msg != nil {
		cc.ReleaseMessage(m.private.msg)
		m.private.msg = nil
	}
}

func (m *midElement) IsExpired(now time.Time, maxRetransmit int32) bool {
	if !m.deadline.IsZero() && now.After(m.deadline) {
		// remove element if deadline is exceeded
		return true
	}
	retransmit := m.retransmit.Load()
	return retransmit >= maxRetransmit
}

func (m *midElement) Retransmit(now time.Time, acknowledgeTimeout time.Duration) bool {
	if now.After(m.start.Add(acknowledgeTimeout * time.Duration(m.retransmit.Load()+1))) {
		m.retransmit.Inc()
		// retransmit
		return true
	}
	// wait for next retransmit
	return false
}

func (m *midElement) GetMessage(cc *Conn) (*pool.Message, bool, error) {
	m.private.Lock()
	defer m.private.Unlock()
	if m.private.msg == nil {
		return nil, false, nil
	}
	msg := cc.AcquireMessage(m.private.msg.Context())
	if err := m.private.msg.Clone(msg); err != nil {
		cc.ReleaseMessage(msg)
		return nil, false, err
	}
	return msg, true, nil
}
