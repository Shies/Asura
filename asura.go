package Asura

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(c *Context)
}

// HandlerFunc http request handler function.
type HandlerFunc func(*Context)

// ServeHTTP calls f(ctx).
func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

var (
	_     IRouter = &Engine{}
)

// Config is the engine configuration struct.
type Config struct {
	// Timeout is used to deliver global timeout to every handler
	Timeout time.Duration
}

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup

	lock sync.RWMutex
	conf *Config

	address string

	mux       *http.ServeMux                    // http mux router
	server    atomic.Value                      // store *http.Server
	metastore map[string]map[string]interface{} // metastore is the path as key and the metadata of this path as value, it export via /metadata
}

// New returns a new blank Engine instance without any middleware attached.
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		address: "0.0.0.0",
		conf: &Config{
			Timeout: time.Second,
		},
		mux:       http.NewServeMux(),
		metastore: make(map[string]map[string]interface{}),
	}
	engine.RouterGroup.engine = engine
	// NOTE add prometheus monitor location
	engine.addRoute("GET", "/metrics", Monitor())
	engine.addRoute("GET", "/metadata", engine.metadata())
	return engine
}

// Default returns an Engine instance with the Recovery, Logger and CSRF middleware already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger(), CSRF(), Mobile())
	return engine
}

func (engine *Engine) addRoute(method, path string, handlers ...HandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if method == "" {
		panic("HTTP method can not be empty")
	}
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}
	if _, ok := engine.metastore[path]; !ok {
		engine.metastore[path] = make(map[string]interface{})
	}
	engine.metastore[path]["method"] = method
	engine.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		c := &Context{
			Context:  nil,
			engine:   engine,
			index:    -1,
			handlers: nil,
			Keys:     nil,
			method:   "",
			Error:    nil,
		}

		c.Request = req
		c.Writer = w
		c.handlers = handlers
		c.method = method

		engine.handleContext(c)
	})
}

// SetConfig is used to set the engine configuration.
// Only the valid config will be loaded.
func (engine *Engine) SetConfig(conf *Config) (err error) {
	if conf.Timeout <= 0 {
		return errors.New("Asura: config timeout must greater than 0")
	}
	engine.lock.Lock()
	engine.conf = conf
	engine.lock.Unlock()
	return
}

func (engine *Engine) handleContext(c *Context) {
	var cancel func()
	// get derived timeout from http request header,
	// compare with the engine configured,
	// and use the minimum one
	engine.lock.RLock()
	tm := engine.conf.Timeout
	engine.lock.RUnlock()
	if tm > 0 {
		c.Context, cancel = context.WithTimeout(context.Background(), tm)
	} else {
		c.Context, cancel = context.WithCancel(context.Background())
	}
	defer cancel()
	c.Next()
}

// Router return a http.Handler for using http.ListenAndServe() directly.
func (engine *Engine) Router() http.Handler {
	return engine.mux
}

// Server return a http.Server for using server.Shutdown() directly.
func (engine *Engine) Server() *http.Server {
	s, ok := engine.server.Load().(*http.Server)
	if !ok {
		return nil
	}
	return s
}

// UseFunc attachs a global middleware to the router. ie. the middleware attached though UseFunc() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) UseFunc(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.UseFunc(middleware...)
	return engine
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...Handler) IRoutes {
	engine.RouterGroup.Use(middleware...)
	return engine
}

// Ping is used to set the general HTTP ping handler.
func (engine *Engine) Ping(handler HandlerFunc) {
	engine.GET("/monitor/ping", handler)
}

// Register is used to export metadata to discovery.
func (engine *Engine) Register(handler HandlerFunc) {
	engine.GET("/register", handler)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	server := &http.Server{
		Addr:    address,
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.ListenAndServe(); err != nil {
		err = errors.Wrapf(err, "addrs: %v", addr)
	}
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.ListenAndServeTLS(certFile, keyFile); err != nil {
		err = errors.Wrapf(err, "tls: %s/%s:%s", addr, certFile, keyFile)
	}
	return
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
		return
	}
	defer listener.Close()
	server := &http.Server{
		Handler: engine.mux,
	}
	engine.server.Store(server)
	if err = server.Serve(listener); err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
	}
	return
}

// RunServer will serve and start listening HTTP requests by given server and listener.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunServer(server *http.Server, l net.Listener) (err error) {
	server.Handler = engine.mux
	engine.server.Store(server)
	if err = server.Serve(l); err != nil {
		err = errors.Wrapf(err, "listen server: %+v/%+v", server, l)
		return
	}
	return
}

func (engine *Engine) metadata() HandlerFunc {
	return func(c *Context) {
		c.JSON(engine.metastore, nil)
	}
}
