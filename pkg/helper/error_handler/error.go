package error_handler

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type ErrorString struct {
	code       int
	stacktrace string
	message    string
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Stacktrace() string {
	return e.stacktrace
}

func DefaultError(cause error, vals ...interface{}) error {
	if cause == nil {
		return nil
	}

	_, ok := cause.(*ErrorString)
	if ok {
		return cause
	}

	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf("\nMessage: %s %s", cause.Error(), j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf("\n--- at %s:%d ---", file, line)
	}

	return &ErrorString{500, stacktrace, cause.Error()}
}

func NewError(code int, message string, vals ...interface{}) error {
	j, _ := json.Marshal(vals)
	stacktrace := fmt.Sprintf("\nMessage: %s %s", message, j)

	for i := 1; i <= 3; i++ {
		pc, file, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		if f == nil || line == 0 {
			break
		}

		stacktrace += fmt.Sprintf("\n--- at %s:%d ---", file, line)
	}

	return &ErrorString{code, stacktrace, message}
}
