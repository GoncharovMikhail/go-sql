package errors

import "fmt"

type Errors interface {
	GetErrMsg() string
	Get() error
	Cause() Errors
}

type errors struct {
	errMsg string
	err    error
	next   Errors
}

func (e *errors) GetErrMsg() string {
	return e.errMsg
}

func (e *errors) Get() error {
	return e.err
}

func (e *errors) Cause() Errors {
	return e.next
}

func NewErrors(msg string, err error, next Errors) Errors {
	return &errors{
		errMsg: msg,
		err:    err,
		next:   next,
	}
}

func BuildSimpleErrMsg(errName string, err error) string {
	return fmt.Sprintf(errName+": <%v>", err)
}
