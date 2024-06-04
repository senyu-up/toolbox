package su_error

import "fmt"

type ErrCode interface {
	Error() string
	Msg() string
	Code() uint32
}

// IsErrCode
// 判断是否为 ErrCode 类型的 err
func IsErrCode(err error) bool {
	if _, ok := err.(ErrCode); ok {
		return true
	} else {
		return false
	}
}

// 尝试获取 ErrCode 类型的 err
func GetErrCode(err error) (ErrCode, bool) {
	if val, ok := err.(ErrCode); ok {
		return val, true
	} else {
		return nil, false
	}
}

type errCode struct {
	msg  string
	err  string
	code uint32
}

func NewErrCode(err, msg string, code uint32) *errCode {
	return &errCode{msg: msg, err: err, code: code}
}

func ErrToErrCode(err error, msg string, code uint32) *errCode {
	if val, ok := err.(ErrCode); ok {
		var newErr = &errCode{err: val.Error(), code: val.Code(), msg: val.Msg()}
		if msg != "" {
			newErr.msg = msg
		}
		return newErr
	}
	if val, ok := err.(error); ok {
		return &errCode{err: val.Error(), code: code, msg: msg}
	}
	return &errCode{err: fmt.Sprintf("%v", err), code: code, msg: msg}
}

func (e *errCode) Error() string {
	return e.err
}

func (e *errCode) Msg() string {
	return e.msg
}

func (e *errCode) Code() uint32 {
	return e.code
}
