package coder

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/codes"
)

type Coder struct {
	secret []byte
}

var DefaultCoder = new(Coder)

// func MakeDefaultCoder(isSecretReady bool, secret []byte) *Coder {
// 	cdr := DefaultCoder
// 	if isSecretReady {
// 		DefaultCoder.SetSecret(secret)
// 	}
// 	return cdr
// }

func (c *Coder) SetSecret(secret []byte) *Coder {
	fmt.Printf("Set Coder secret: %x\n", secret)
	c.secret = secret
	return c
}

// size in bytes
func (c *Coder) Size(m message.Message) (int, error) {
	if len(m.Token) > message.MaxTokenSize {
		return -1, message.ErrInvalidTokenLen
	}
	size := 4 + len(m.Token)
	payloadLen := len(m.Payload)
	optionsLen, err := m.Options.Marshal(nil)
	if !errors.Is(err, message.ErrTooSmall) {
		return -1, err
	}
	if payloadLen > 0 {
		// for separator 0xff
		payloadLen++
	}
	size += payloadLen + optionsLen
	return size, nil
}

func (c *Coder) Encode(m message.Message, buf []byte) (int, error) {
	/*
	     0                   1                   2                   3
	    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |Ver| T |  TKL  |      Code     |          Message ID           |
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   Token (if any, TKL bytes) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |   Options (if any) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	   |1 1 1 1 1 1 1 1|    Payload (if any) ...
	   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	if !message.ValidateMID(m.MessageID) {
		return -1, fmt.Errorf("invalid MessageID(%v)", m.MessageID)
	}
	if !message.ValidateType(m.Type) {
		return -1, fmt.Errorf("invalid Type(%v)", m.Type)
	}
	size, err := c.Size(m)
	if err != nil {
		return -1, err
	}
	if len(buf) < size {
		return size, message.ErrTooSmall
	}

	fullBuf := buf
	tmpbuf := []byte{0, 0}
	binary.BigEndian.PutUint16(tmpbuf, uint16(m.MessageID))

	buf[0] = (1 << 6) | byte(m.Type)<<4 | byte(0xf&len(m.Token))
	buf[1] = byte(m.Code)
	buf[2] = tmpbuf[0]
	buf[3] = tmpbuf[1]
	buf = buf[4:]

	if len(m.Token) > message.MaxTokenSize {
		return -1, message.ErrInvalidTokenLen
	}
	copy(buf, m.Token)
	buf = buf[len(m.Token):]

	optionsLen, err := m.Options.Marshal(buf)
	switch {
	case err == nil:
	case errors.Is(err, message.ErrTooSmall):
		return size, err
	default:
		return -1, err
	}
	buf = buf[optionsLen:]

	if len(m.Payload) > 0 {
		buf[0] = 0xff
		buf = buf[1:]
	}

	copy(buf, m.Payload)

	if c.secret != nil {
		fmt.Println("Encrypting")
		nonce := RandomBytes(16)

		tmp, tag := Encrypt(c.secret, nonce, fullBuf[:size])

		copy(fullBuf, tmp)

		buf = fullBuf[:size]
		buf = append(buf, tag...)
		size += 16
		buf = append(buf, nonce[:]...)
		size += 16 //nonce 16 + tag 16
	}

	return size, nil
}

func (c *Coder) Decode(data []byte, m *message.Message) (int, error) {
	size := len(data)
	if size < 4 {
		return -1, ErrMessageTruncated
	}

	if c.secret != nil {
		fmt.Println("Decrypting")
		copy(data, Decrypt(c.secret, data[size-16:], data[:size-32], data[size-32:size-16]))

		data = data[:size-32] // local
		size -= 32            //nonce 16 + tag 16
	}

	if data[0]>>6 != 1 {
		return -1, ErrMessageInvalidVersion
	}

	typ := message.Type((data[0] >> 4) & 0x3)
	tokenLen := int(data[0] & 0xf)
	if tokenLen > 8 {
		return -1, message.ErrInvalidTokenLen
	}

	code := codes.Code(data[1])
	messageID := binary.BigEndian.Uint16(data[2:4])
	data = data[4:]
	if len(data) < tokenLen {
		return -1, ErrMessageTruncated
	}
	token := data[:tokenLen]
	if len(token) == 0 {
		token = nil
	}
	data = data[tokenLen:]

	optionDefs := message.CoapOptionDefs
	proc, err := m.Options.Unmarshal(data, optionDefs)
	if err != nil {
		return -1, err
	}
	data = data[proc:]
	if len(data) == 0 {
		data = nil
	}

	m.Payload = data
	m.Code = code
	m.Token = token
	m.Type = typ
	m.MessageID = int32(messageID)

	return size, nil
}

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}
