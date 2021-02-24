package Asura

import (
	"net/url"
	"strings"

	"Asura/logger"
)

var (
	_allowReferer = []string{
	}
)

// CSRF returns the csrf middleware to prevent invalid cross site request.
// Only referer is checked currently.
func CSRF() HandlerFunc {
	return func(c *Context) {
		referer := c.Request.Header.Get("Referer")
		params := c.Request.Form
		cross := (params.Get("callback") != "" && params.Get("jsonp") == "jsonp") || (params.Get("cross_domain") != "")
		if referer == "" {
			if !cross {
				return
			}
			logger.Error("The request's Referer header is empty.")
			c.AbortWithStatus(403)
			return
		}
		illegal := true
		if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
			for _, r := range _allowReferer {
				if strings.HasSuffix(strings.ToLower(uri.Host), r) {
					illegal = false
					break
				}
			}
		}
		if illegal {
			logger.Error("The request's Referer header `%s` does not match any of allowed referers.", referer)
			c.AbortWithStatus(403)
			return
		}
	}
}
