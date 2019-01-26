package logger

import (
	"context"
	"fmt"
	"time"
)

// Config log config.
type Config struct {
	Dir	string
}

// D represents a map of entry level data used for structured logging.
type D map[string]interface{}

var (
	h Handler
	c *Config
)

// Init create logger with context.
func Init(conf *Config) {
	if conf == nil {
		conf = &Config{}
	}
	var (
		hs Handlers
	)
	c = conf
	if conf.Dir != "" {
		hs = append(hs, NewFile(conf.Dir))
	}
	h = hs
}

// Info logs a message at the info log level.
func Info(format string, args ...interface{}) {
	logf(_infoLevel, format, args...)
}

// Warn logs a message at the warning log level.
func Warn(format string, args ...interface{}) {
	logf(_warnLevel, format, args...)
}

// Error logs a message at the error log level.
func Error(format string, args ...interface{}) {
	logf(_errorLevel, format, args...)
}

// Infov logs a message at the info log level.
func Infov(c context.Context, d D) {
	logv(c, _infoLevel, d)
}

// Warnv logs a message at the warning log level.
func Warnv(c context.Context, d D) {
	logv(c, _warnLevel, d)
}

// Errorv logs a message at the error log level.
func Errorv(c context.Context, d D) {
	logv(c, _errorLevel, d)
}

func logf(lv Level, format string, args ...interface{}) {
	if h == nil {
		return
	}
	now := time.Now()
	d := D{}
	d[_level] = lv.String()
	d[_time] = now.Format(_timeFormat)
	d[_source] = funcName()
	d[_log] = fmt.Sprintf(format, args...)
	h.Log(lv, d)
}

func logv(ctx context.Context, lv Level, d D) {
	if h == nil {
		return
	}
	now := time.Now()
	d[_levelValue] = lv
	d[_level] = lv.String()
	d[_time] = now.Format(_timeFormat)
	d[_source] = funcName()
	d[_ts] = now.Unix()
	h.Log(lv, d)
}

// Close close resource.
func Close() (err error) {
	if h == nil {
		return
	}
	return h.Close()
}
