package ecode_test

import (
	"fmt"
	"Asura/src/ecode"

	"github.com/pkg/errors"
)

func Example_ecode_Message() {
	ecode.NotModified.Message()
}

func Example_ecode_Code() {
	ecode.NotModified.Code()
}

func Example_ecode_Error() {
	_ = ecode.NotModified.Error()
}

func ExampleCause() {
	err := errors.WithStack(ecode.AccessDenied)
	ecode.Cause(err)
}

func ExampleInt() {
	ecode.Int(500)
}

func ExampleString() {
	ecode.String("500")
}

// ExampleStack package error with stack.
func Example() {
	err := errors.New("dao error")
	errors.Wrap(err, "some message")
	// package ecode with stack.
	errCode := ecode.AccessDenied
	err = errors.Wrap(errCode, "some message")

	//get ecode from package error
	code := errors.Cause(err).(ecode.Error)
	fmt.Printf("%d: %s\n", code.Code(), code.Message())
}
