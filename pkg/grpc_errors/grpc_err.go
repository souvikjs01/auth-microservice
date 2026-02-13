package grpc_errors

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
)

var (
	ErrNotFound         = errors.New("ErrNotFound")
	ErrNoCtxMetaData    = errors.New("ErrNoCtxMetaData")
	ErrInvalidSessionId = errors.New("ErrInvalidSessionId")
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
	}

	return codes.Internal
}
