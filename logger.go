package blade

import (
	"flag"
	"os"
	"fmt"
	"log"
)

// Logger init middleware
func Logger() HandlerFunc {
	return func(c *Context) {
		flag.Parse()
		outfile, err := os.OpenFile("/tmp", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(*outfile, "open failed")
			os.Exit(1)
		}
		log.SetOutput(outfile)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
}