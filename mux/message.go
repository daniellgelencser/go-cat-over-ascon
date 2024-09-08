package mux

import "github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"

// RouteParams contains all the information related to a route
type RouteParams struct {
	Path         string
	Vars         map[string]string
	PathTemplate string
}

// Message contains message with sequence number.
type Message struct {
	*pool.Message
	RouteParams *RouteParams
}
