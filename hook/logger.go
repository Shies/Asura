package hook

import (
	"fmt"

	"blade"
	"blade/ecode"
	"blade/slog"
)

// Logger init middleware
func Logger() blade.HandlerFunc {
	return func(c *blade.Context) {
		ip := c.RemoteIP()
		req := c.Request
		path := req.URL.Path
		params := req.Form

		c.Next()

		err := c.Error
		cerr := ecode.Cause(err)
		lv := slog.D{
			"method": req.Method,
			"ip":     ip,
			"path":   path,
			"params": params.Encode(),
			"ret":    cerr.Code(),
			"msg":    cerr.Message(),
			"stack":  fmt.Sprintf("%+v", err),
		}
		lf := slog.Infov
		if err != nil {
			lv["err"] = err.Error()
			lf = slog.Errorv
		}
		lf(c, lv)
	}
}