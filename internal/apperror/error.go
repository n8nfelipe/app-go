package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeNotFound     Code = "not_found"
	CodeUnauthorized Code = "unauthorized"
	CodeInvalidInput Code = "invalid_input"
	CodeInternal     Code = "internal"
)

type Error interface {
	error
	Code() Code
	PublicMessage() string
}

type appError struct {
	code    Code
	message string
	cause   error
}

func (e *appError) Error() string {
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s: %v", e.message, e.cause)
}

func (e *appError) Unwrap() error {
	return e.cause
}

func (e *appError) Code() Code {
	return e.code
}

func (e *appError) PublicMessage() string {
	return e.message
}

func New(code Code, publicMessage string, cause error) Error {
	return &appError{code: code, message: publicMessage, cause: cause}
}

func Wrap(err error, context string) error {
	if err == nil {
		return nil
	}
	if Is(err) {
		return fmt.Errorf("%s: %w", context, err)
	}
	return New(CodeInternal, "Erro interno no servidor", fmt.Errorf("%s: %w", context, err))
}

func Is(err error) bool {
	var appErr Error
	return errors.As(err, &appErr)
}
