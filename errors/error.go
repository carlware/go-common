package errors

import (
	"bytes"
	"fmt"
)

// Based on the post of https://middlemost.com/failure-is-your-domain/

// An Error describes a Go CDK error.
type Error struct {
	Code    ErrorCode
	Message string
	Op      string
	Err     error
	Details []string
}

// New creates and returns a new error
func New(code ErrorCode, op, message string, err error) *Error {
	return &Error{Op: op, Code: code, Message: message, Err: err}
}

// NewD creates and returns a new error with details
func NewD(code ErrorCode, op, message string, details []string, err error) *Error {
	return &Error{Op: op, Code: code, Message: message, Err: err, Details: details}
}

func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		fmt.Fprintf(&buf, "<%s> %s", e.Code, e.Details)
		if e.Message != "" {
			buf.WriteString(" " + e.Message)
		}
	}
	return buf.String()

}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

// ErrorDetails returns details if available.
func ErrorDetails(err error) []string {
	if err == nil {
		return []string{}
	} else if e, ok := err.(*Error); ok {
		return e.Details
	}
	return []string{}
}

// As return an Error or transform common error to Error.
func As(err error) *Error {
	def := New(Unimplemented, "", "", nil)
	if err == nil {
		return def
	} else if e, ok := err.(*Error); ok {
		return e
	} else if !ok {
		return &Error{
			Code:    Internal,
			Message: err.Error(),
			Err:     err,
		}
	}
	return def
}
