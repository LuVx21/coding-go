package errorx

import "net/http"

const defaultCode = 1001

type CodeError struct {
    Code int
    Msg  string
}

func (e *CodeError) Error() string {
    return e.Msg
}

func NewCodeError(code int) error {
    msg := http.StatusText(http.StatusNotFound)
    return NewCodeMsgError(code, msg)
}

func NewCodeMsgError(code int, msg string) error {
    return &CodeError{Code: code, Msg: msg}
}

func NewDefaultError(msg string) error {
    return NewCodeMsgError(defaultCode, msg)
}
