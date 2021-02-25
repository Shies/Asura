package server

import (
	"net"
	"context"
	"time"

	"Asura/conf"
	. "Asura/app/rpc"
	"Asura/app/service"
	"Asura/src/ecode"
	"Asura/src/logger"

	"google.golang.org/grpc"
)

var (
	srv *service.Service
	cnf = &conf.Config{}
)

const (
	Port = ":41005"
)

type Transport struct {
	// to do sth.
}

func initService(s *service.Service, c *conf.Config) {
	srv, cnf = s, c
}

func Init(c *conf.Config, s *service.Service) {
	initService(s, c)

	lis, err := net.Listen(c.RPCServer2.Servers[0].Proto, c.RPCServer2.Servers[0].Addr)
	if err != nil {
		logger.Error("failed to listen: %v", err)
	}
	g := grpc.NewServer()
	RegisterTransportServer(g, &Transport{})
	err = g.Serve(lis)
	if err != nil {
		logger.Info("grpc server in: %s, (%v)", Port, err)
	}
	return
}

func (t *Transport) Ping(c context.Context, req *Request) (res *Response, err error) {
	return &Response{
		Code: int32(ecode.OK.Code()),
		Message: ecode.OK.Message(),
		Ttl: time.Now().Unix(),
		Data: string("ok welcome!"),
	}, nil
}