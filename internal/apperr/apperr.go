package apperr

import (
	"bytes"
	"strconv"
)

// ErrStatus is specific error status for the error i.e. USER_NOT_FOUND, USER_EXISTS and not NOT_FOUND (generic)
// ErrStatus can be used by client as well for error handling and can also be used by caller for better control flow
// ErrInternalError is exception to this. BUT only use it in case where the error is not expected at all
type ErrStatus string

const (
	ErrBadRequest         ErrStatus = "BAD_REQUEST"
	ErrUnknown            ErrStatus = "UNKNOWN"
	ErrInternalError      ErrStatus = "INTERNAL_ERROR"
	ErrStatusUserNotFound ErrStatus = "USER_NOT_FOUND"
)

type AppError struct {
	Code   int       `json:"code"`
	Msg    string    `json:"msg"`
	Err    error     `json:"-"`
	Status ErrStatus `json:"status"`
}

func New(code int, err error, msg string, status ErrStatus) *AppError {
	return &AppError{
		Code:   code,
		Msg:    msg,
		Err:    err,
		Status: status,
	}
}

func (a *AppError) Error() string {
	var buf bytes.Buffer

	buf.WriteString("code: ")
	buf.WriteString(strconv.Itoa(a.Code))

	buf.WriteString(", msg: ")
	buf.WriteString(a.Msg)

	buf.WriteString(", status: ")
	buf.WriteString(string(a.Status))

	if a.Err != nil {
		buf.WriteString(", err: ")
		buf.WriteString(a.Err.Error())
	}

	return buf.String()
}
