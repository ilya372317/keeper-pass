package v1

import (
	"database/sql"
	"errors"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkErr(kind string, id int64, err error) *status.Status {
	if errors.Is(err, domain.ErrNotSupportedOperation) {
		return status.Newf(codes.InvalidArgument, "data with id [%d] is not %s type", id, kind)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return status.Newf(codes.NotFound, "%s data with id [%d] not found", kind, id)
	}
	if errors.Is(err, domain.ErrAccesDenied) {
		return status.Newf(codes.PermissionDenied, "you can`t access %s data with id [%d]", kind, id)
	}
	return status.Newf(codes.Internal, "failed get %s data for show: %v", kind, err)
}
