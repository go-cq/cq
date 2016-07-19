package cq

import (
  "net/http"
  "github.com/origininvesting/cq/types"
)

var (
  transport http.RoundTripper = &http.Transport{}
  client = &http.Client{
    Transport: transport,
  }
)

func SetTransport(rt http.RoundTripper) {
	transport = rt
	client.Transport = transport
  types.SetTransport(rt)
}
