package http

import (
	"net"
	"net/http"
	"strconv"
	"time"

	xtime "Asura/app/time"
	"Asura/app/service"
	"Asura/conf"
	blade "Asura/src"
	"Asura/src/logger"
	"Asura/src/render"

	"github.com/pkg/errors"
)

var (
	srv *service.Service
	cnf = &conf.Config{}
)

const (
	_defaultAddr = ":9701"
)

func ParseInt(value string) int64 {
	intval, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		intval = 0
	}

	return intval
}

func Atoi(value string) int {
	intval, err := strconv.Atoi(value)
	if err != nil {
		intval = 0
	}

	return intval
}

func initService(c *conf.Config, s *service.Service) {
	srv, cnf = s, c
}

func initRouter(engine *blade.Engine) {
	engine.GET("/ping", func(c *blade.Context) {
		c.String(200, "%s", "hello world")
	})
	v1 := engine.Group("/v1")
	{
		v1.GET("/ping", ping)
		v1.GET("/test", test)
	}
}

func Init(c *conf.Config, s *service.Service) (err error) {
	initService(c, s)
	if c == nil {
		c = &conf.Config{
			HttpServer: &conf.HTTPServer{
				Addrs:    []string{_defaultAddr},
				Timeout: xtime.Duration(time.Second),
			},
		}
	}
	l, err := net.Listen("tcp", c.HttpServer.Addrs[0])
	if err != nil {
		err = errors.Wrapf(err, "listen tcp: %d", c.HttpServer.Addrs[0])
		return
	}

	engine := blade.Default()
	if err = engine.SetConfig(&blade.Config{Timeout: time.Duration(c.HttpServer.Timeout)}); err != nil {
		return
	}

	initRouter(engine)
	logger.Info("start http listen addr: %s", c.HttpServer.Addrs[0])
	server := &http.Server{
		ReadTimeout:  time.Duration(c.HttpServer.ReadTimeout),
		WriteTimeout: time.Duration(c.HttpServer.WriteTimeout),
	}
	go func() {
		if err = engine.RunServer(server, l); err != nil {
			logger.Error("Fatal: engine.ListenServer(%+v, %+v) error(%v)", server, l, err)
		}
	}()

	return nil
}

func ping(c *blade.Context) {
	c.JSON(srv.Commits, nil)
	/*
	var err error
	if err = srv.Ping(c); err != nil {
		logger.Error("service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
	*/
}

func renderErrMsg(c *blade.Context, code int, msg string) {
	data := map[string]interface{} {
		"code": code,
		"message": msg,
	}
	c.Render(http.StatusOK, render.MapJSON(data))
}