package types

import (
  "net/http"
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
}
