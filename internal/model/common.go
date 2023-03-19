package model

import "errors"

var (
	ErrTokenInvalid = errors.New("invalid token")
)

// Response :nodoc:
type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// NewResponse :nodoc:
func NewResponse() *Response {
	return new(Response)
}

// NewDefaultResponse :nodoc:
func NewDefaultResponse() *Response {
	return &Response{
		Message: "OK",
	}
}

func WithBadRequestResponse(errors any) *Response {
	return &Response{
		Message: "Bad Request",
		Errors:  errors,
	}
}

// WithMessage :nodoc:
func (m *Response) WithMessage(message string) *Response {
	m.Message = message
	return m
}

// WithErrorMessage :nodoc:
func (m *Response) WithErrorMessage(err error) *Response {
	m.Message = err.Error()
	return m
}

// WithData :nodoc:
func (m *Response) WithData(data any) *Response {
	m.Data = data
	return m
}
