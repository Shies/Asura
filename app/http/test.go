package http

import (
	blade "Asura/src"
	"Asura/src/ecode"
)

func test(c *blade.Context) {
	var (
		_ = c.Request.ParseForm()
		req = c.Request.Form
		key = req.Get("key")
	)
	if key == "" {
		c.JSON(nil, ecode.AccessKeyErr)
		return
	}
	test, err := srv.Test(c, key)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(test, nil)
}
