package apperr

import "fmt"

type AppError struct {
	Code    int               // stable code (exposed)
	Message string            // safe message for client
	Meta    map[string]string // safe metadata (optional)
	Err     error             // internal cause (not exposed)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code : %d, message : %s, err : %v", e.Code, e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

func NewWithMeta(code int, msg string, meta map[string]string, err error) *AppError {
	appErr := New(code, msg, err)
	appErr.Meta = meta

	return appErr
}
