package errors

import (
	"context"
	"io"

	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoNotWrap reports whether an error should not be wrapped in the Error
// type from this package.
// It returns true if err is a retry error, a context error, io.EOF, or if it wraps
// one of those.
func DoNotWrap(err error) bool {
	if xerrors.Is(err, io.EOF) {
		return true
	}
	if xerrors.Is(err, context.Canceled) {
		return true
	}
	if xerrors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}

var (
	grpcCodeMap = map[codes.Code]ErrorCode{
		codes.AlreadyExists:      AlreadyExists,
		codes.Aborted:            Aborted,
		codes.Canceled:           Canceled,
		codes.DataLoss:           DataLoss,
		codes.DeadlineExceeded:   DeadlineExceeded,
		codes.FailedPrecondition: FailedPrecondition,
		codes.Internal:           Internal,
		codes.InvalidArgument:    InvalidArgument,
		codes.NotFound:           NotFound,
		codes.OK:                 OK,
		codes.OutOfRange:         OutOfRange,
		codes.PermissionDenied:   PermissionDenied,
		codes.Unauthenticated:    Unauthenticated,
		codes.Unavailable:        Unavailable,
		codes.Unimplemented:      Unimplemented,
		codes.Unknown:            Unknown,
	}
)

// GRPCCode extracts the gRPC status code and converts it into an ErrorCode.
// It returns Unknown if the error isn't from gRPC.
func GRPCCode(err error) ErrorCode {
	if code, ok := grpcCodeMap[status.Code(err)]; ok {
		return code
	}
	return Unknown
}
