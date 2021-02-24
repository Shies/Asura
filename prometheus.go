package Asura

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Monitor() HandlerFunc {
	return func(c *Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	}
}
