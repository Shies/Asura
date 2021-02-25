package conf

import (
	"os"
	"sync"
	"syscall"
	"os/signal"

	"Asura/src/logger"
)

var (
	confLock = new(sync.RWMutex)
)

func Load() *Config {
	confLock.RLock()
	defer confLock.RUnlock()

	return &Conf
}

func Init() {
	if err := ParseConfig(); err != nil {
		logger.Error("conf Init() error(%v)", err)
		panic(err)
	}

	// 热更新配置可能有多种触发方式，这里使用系统信号量sigusr1实现
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
			logger.Info("Reloaded config:", ParseConfig())
		}
	}()

	return
}

func Listen() {
	select {}
}