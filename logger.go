package blade

import (
	"fmt"

	"blade/ecode"
	"blade/log"
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
		lv := log.D{
			"method": req.Method,
			"ip":     ip,
			"path":   path,
			"params": params.Encode(),
			"ret":    cerr.Code(),
			"msg":    cerr.Message(),
			"stack":  fmt.Sprintf("%+v", err),
		}
		lf := log.Infov
		if err != nil {
			lv["err"] = err.Error()
			lf = log.Errorv
		}
		lf(c, lv)
	}
}