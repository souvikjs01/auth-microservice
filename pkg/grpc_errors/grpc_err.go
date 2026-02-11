package grpc_errors

import (
	"database/sql"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
)

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "password"):
		return codes.InvalidArgument
	}

	return codes.Internal
}
