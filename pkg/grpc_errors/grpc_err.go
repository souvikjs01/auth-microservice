package grpc_errors

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrNoCtxMetaData    = errors.New("no context metaData")
	ErrInvalidSessionId = errors.New("invalid sessionId")
	ErrEmailExists      = errors.New("email already exists")
)

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "password"):
		return codes.InvalidArgument
	case strings.Contains(err.Error(), "redis"):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrEmailExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrNoCtxMetaData):
		return codes.Unauthenticated
	case errors.Is(err, ErrInvalidSessionId):
		return codes.PermissionDenied
	}

	return codes.Internal
}

// map gRPC error code to HTTP status code
func MapGRPCCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.AlreadyExists:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
