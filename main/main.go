package main

import (
	blade "Asura/src"
)

func main() {
	engine := blade.Default()
	engine.GET("/welcome", func(c *blade.Context) {
		c.String(200, "%s", "hello world !!!")
	})
	engine.Run(":8080")
}
