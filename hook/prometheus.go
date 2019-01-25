package hook

import (
	"blade"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Monitor() blade.HandlerFunc {
	return func(c *blade.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	}
}
