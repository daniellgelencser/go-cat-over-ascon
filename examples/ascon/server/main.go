package main

import (
	"bytes"
	"log"

	coap "go-attested-coap-over-ascon/v3"
	"go-attested-coap-over-ascon/v3/message"
	"go-attested-coap-over-ascon/v3/message/codes"
	"go-attested-coap-over-ascon/v3/mux"
)

func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Conn().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func handleA(w mux.ResponseWriter, r *mux.Message) {
	err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("mattis aliquam faucibus purus in massa tempor nec feugiat nisl pretium fusce id velit ut tortor pretium viverra suspendisse")))
	if err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func handleB(w mux.ResponseWriter, r *mux.Message) {
	customResp := w.Conn().AcquireMessage(r.Context())
	defer w.Conn().ReleaseMessage(customResp)
	customResp.SetCode(codes.Content)
	customResp.SetToken(r.Token())
	customResp.SetContentFormat(message.TextPlain)
	customResp.SetBody(bytes.NewReader([]byte("B hello world")))
	err := w.Conn().WriteMessage(customResp)
	if err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Handle("/a", mux.HandlerFunc(handleA))
	r.Handle("/b", mux.HandlerFunc(handleB))

	log.Fatal(coap.ListenAndServerASCON("udp", ":5688", r))
}
