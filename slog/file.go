package slog

import (
	"encoding/json"
	"math"
	"path"

	"golang/log4go"
)

// FileHandler .
type FileHandler struct {
	logger log4go.Logger
}

// NewFile crete a file logger.
func NewFile(dir string) *FileHandler {
	var l = log4go.Logger{}
	log4go.LogBufferLength = 10240
	// new info writer
	iw := log4go.NewFileLogWriter(path.Join(dir, "info.log"), false)
	iw.SetRotateDaily(true)
	iw.SetRotateSize(math.MaxInt32)
	iw.SetFormat("[%D %T] [%L] [%S] %M")
	l["info"] = &log4go.Filter{
		Level:     log4go.INFO,
		LogWriter: iw,
	}
	// new warning writer
	ww := log4go.NewFileLogWriter(path.Join(dir, "warning.log"), false)
	ww.SetRotateDaily(true)
	ww.SetRotateSize(math.MaxInt32)
	ww.SetFormat("[%D %T] [%L] [%S] %M")
	l["warning"] = &log4go.Filter{
		Level:     log4go.WARNING,
		LogWriter: ww,
	}
	// new error writer
	ew := log4go.NewFileLogWriter(path.Join(dir, "error.log"), false)
	ew.SetRotateDaily(true)
	ew.SetRotateSize(math.MaxInt32)
	ew.SetFormat("[%D %T] [%L] [%S] %M")
	l["error"] = &log4go.Filter{
		Level:     log4go.ERROR,
		LogWriter: ew,
	}
	return &FileHandler{logger: l}
}

// Log loggint to file .
func (h *FileHandler) Log(lv Level, d D) {
	msg, err := json.Marshal(d)
	if err != nil {
		return
	}
	switch lv {
	case _debugLevel:
		h.logger.Debug(string(msg))
	case _infoLevel:
		h.logger.Info(string(msg))
	case _warnLevel:
		h.logger.Warn(string(msg))
	case _errorLevel:
		h.logger.Error(string(msg))
	case _fatalLevel:
		h.logger.Critical(string(msg))
	default:
	}
}

// Close .
func (h *FileHandler) Close() (err error) {
	h.logger.Close()
	return nil
}
