package Asura

import (
	"fmt"

	"Asura/src/ecode"
	"Asura/src/logger"
)

// Logger init middleware
func Logger() HandlerFunc {
	return func(c *Context) {
		ip := c.RemoteIP()
		req := c.Request
		path := req.URL.Path
		params := req.Form

		c.Next()

		err := c.Error
		cerr := ecode.Cause(err)
		lv := logger.D{
			"method": req.Method,
			"ip":     ip,
			"path":   path,
			"params": params.Encode(),
			"ret":    cerr.Code(),
			"msg":    cerr.Message(),
			"stack":  fmt.Sprintf("%+v", err),
		}
		lf := logger.Infov
		if err != nil {
			lv["err"] = err.Error()
			lf = logger.Errorv
		}
		lf(c, lv)
	}
}