package client

import (
	"fmt"
	"time"
	"sync"
	"runtime"

	"Asura/conf"
	. "Asura/app/rpc"
	"Asura/src/logger"
	_ "Asura/src/ecode"

	"google.golang.org/grpc"
)

const (
	_family = "rpc_client"
)

var (
	idx int64 = 0
)

// grpc connection
type Conn struct {
	wg  	*sync.WaitGroup
	conf    *conf.RpcClient
	client  *grpc.ClientConn
	err     error
}

func NewClient(c *conf.RpcClient) *Conn {
	runtime.GOMAXPROCS(runtime.NumCPU())
	currTime := time.Now()

	conn := &Conn{}
	conn.conf = c
	conn.wg = new(sync.WaitGroup)
	// 并行请求
	for i := 0; i < int(c.Retry); i++ {
		conn.wg.Add(1)
		go func(c *Conn) {
			defer conn.wg.Done()
			c.client, c.err = grpc.Dial(c.conf.Addrs[0], grpc.WithInsecure())
		}(conn)
	}
	conn.wg.Wait()
	logger.Info("time taken: %.2f ", time.Now().Sub(currTime).Seconds())
	if conn.err != nil {
		panic(conn.err)
	}

	return conn
}

func (c *Conn) Call(name string, callback func(c TransportClient)) (err error) {
	defer c.client.Close()
	client := NewTransportClient(c.client)
	for i := 0; i < int(c.conf.Times); i++ {
		callback(client)
	}
	return
}

func Welcome(client TransportClient) {
	idx = idx + 1
	fmt.Println("hello world", idx)
}