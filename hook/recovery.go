package hook

import (
	"fmt"
	"net/http/httputil"
	"os"
	"runtime"

	"blade"
	"blade/log"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() blade.HandlerFunc {
	return func(c *blade.Context) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if c.Request != nil {
					rawReq, _ = httputil.DumpRequest(c.Request, false)
				}
				pl := fmt.Sprintf("http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				fmt.Fprintf(os.Stderr, pl)
				log.Info(pl)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
