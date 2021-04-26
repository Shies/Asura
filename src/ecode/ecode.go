package ecode

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/pkg/errors"
)

// All common ecode
var (
	OK = add(0) // 正确

	AppKeyInvalid           = add(-1)   // 应用程序不存在或已被封禁
	AccessKeyErr            = add(-2)   // Access Key错误
	SignCheckErr            = add(-3)   // API校验密匙错误
	MethodNoPermission      = add(-4)   // 调用方对该Method没有权限

	NotModified           = add(-304) // 木有改动
	TemporaryRedirect     = add(-307) // 撞车跳转
	RequestErr            = add(-400) // 请求错误
	Unauthorized          = add(-401) // 未认证
	AccessDenied          = add(-403) // 访问权限不足
	NothingFound          = add(-404) // 啥都木有
	MethodNotAllowed      = add(-405) // 不支持该方法
	Conflict              = add(-409) // 冲突
	ServerErr             = add(-500) // 服务器错误
	ServiceUnavailable    = add(-503) // 过载保护,服务暂不可用
	Deadline              = add(-504) // 服务调用超时
	LimitExceed           = add(-509) // 超出限制

	Degrade     = add(-1200) // 被降级过滤的请求
	RPCNoClient = add(-1201) // rpc服务的client都不可用
	RPCNoAuth   = add(-1202) // rpc服务的client没有授权
)

var (
	_messages atomic.Value         // NOTE: stored map[string]map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

// Register register ecode message map.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

// New new a ecode.Error by int value.
// NOTE: ecode must unique in global, the New will check repeat and then panic.
func New(e int) Error {
	if e <= 0 {
		panic("business ecode must greater than zero")
	}
	return add(e)
}

func add(e int) Error {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

// Error ecode error interface which has a code & message.
type Error interface {
	error
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	// Equal compare whether two errors are equal.
	Equal(error) bool
}

type ecode int

func (e ecode) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

func (e ecode) Code() int {
	return int(e)
}

func (e ecode) Message() (msg string) {
	cm, ok := _messages.Load().(map[int]string)
	if !ok {
		msg = e.Error()
		return
	}
	// get code
	if msg, ok = cm[e.Code()]; ok {
		return
	}
	msg = e.Error()
	return
}

// Equal compare whether two errors are equal.
func (e ecode) Equal(ec error) bool {
	return Cause(ec).Code() == e.Code()
}

// Int parse code int to error.
func Int(i int) Error {
	return ecode(i)
}

// String parse code string to error.
func String(e string) Error {
	if e == "" {
		return OK
	}
	// try error string
	i, err := strconv.Atoi(e)
	if err != nil {
		return ServerErr
	}
	return ecode(i)
}

// Cause cause from error to ecode.
func Cause(e error) Error {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(Error)
	if ok {
		return ec
	}
	return String(e.Error())
}

// Equal equal a and b by code int.
func Equal(a, b Error) bool {
	if a == nil {
		a = OK
	}
	if b == nil {
		b = OK
	}
	return a.Code() == b.Code()
}
