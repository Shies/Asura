package main

import (
	blade "Asura/src"
)

func main() {
	engine := blade.Default()
	engine.GET("/ping", func(c *blade.Context) {
		c.String(200, "%s", "pong")
	})
	engine.Run(":8080")
}
