package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
)

const (
	ExitSuccess       = 0
	ExitInvalidArgs   = 101
	ExitAuth          = 201
	ExitTransport     = 301
	ExitConflict      = 409
	ExitCompatibility = 551
)

const (
	CodeInvalidArgument       = "Cli.InvalidArgument"
	CodeMissingRequestBody    = "Cli.MissingRequestBody"
	CodeInvalidJSONBody       = "Cli.InvalidJSONBody"
	CodeBodyFileReadFailed    = "Cli.BodyFileReadFailed"
	CodeUnauthorized          = "Cli.Unauthorized"
	CodeUserMissingOrDisabled = "Cli.UserMissingOrDisabled"
	CodeTransportTimeout      = "Cli.TransportTimeout"
	CodeTransportConnection   = "Cli.TransportConnectionFailed"
	CodeBackendInvalidResp    = "Cli.BackendInvalidResponse"
	CodeConflict              = "Cli.Conflict"
	CodeUnsupportedVersion    = "Cli.UnsupportedBackendVersion"
)

type Error struct {
	Code     string
	ExitCode int
	Message  string
	Cause    error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code string, exitCode int, message string) *Error {
	return &Error{Code: code, ExitCode: exitCode, Message: message}
}

func Wrap(code string, exitCode int, message string, cause error) *Error {
	return &Error{Code: code, ExitCode: exitCode, Message: message, Cause: cause}
}

func As(err error, target any) bool {
	return stderrors.As(err, target)
}

func ExitCodeOf(err error) int {
	if err == nil {
		return ExitSuccess
	}
	var cliErr *Error
	if As(err, &cliErr) {
		return cliErr.ExitCode
	}
	if isCLIUsageError(err) {
		return ExitInvalidArgs
	}
	return ExitTransport
}

func ClassifyHTTPStatus(status int) *Error {
	switch status {
	case 401:
		return New(CodeUnauthorized, ExitAuth, "backend returned 401 unauthorized")
	case 403:
		return New(CodeUserMissingOrDisabled, ExitAuth, "backend returned 403 user missing or disabled")
	case 409:
		return New(CodeConflict, ExitConflict, "backend returned 409 conflict")
	default:
		return New(CodeBackendInvalidResp, ExitTransport, fmt.Sprintf("unexpected backend status %d", status))
	}
}

func isCLIUsageError(err error) bool {
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	needles := []string{
		"unknown command",
		"unknown shorthand flag",
		"unknown flag",
		"invalid argument",
		"accepts",
		"requires",
		"requires at least",
		"required flag",
		"invalid value",
	}
	for _, needle := range needles {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}
