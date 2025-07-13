package middleware

import (
	"context"
	"net/http"

	"github.com/danipurwadi/internal-transfer-system/foundation/customerror"
	"github.com/danipurwadi/internal-transfer-system/foundation/logger"
	"github.com/danipurwadi/internal-transfer-system/foundation/web"
)

var codeStatus [17]int

// init maps out the error codes to http status codes.
func init() {
	codeStatus[customerror.OK.Value()] = http.StatusOK
	codeStatus[customerror.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatus[customerror.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[customerror.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[customerror.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[customerror.NotFound.Value()] = http.StatusNotFound
	codeStatus[customerror.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[customerror.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[customerror.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[customerror.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[customerror.Aborted.Value()] = http.StatusConflict
	codeStatus[customerror.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[customerror.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[customerror.Internal.Value()] = http.StatusInternalServerError
	codeStatus[customerror.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatus[customerror.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[customerror.Unauthenticated.Value()] = http.StatusUnauthorized
}

// Errors executes the errors middleware functionality.
func Errors(log *logger.Logger) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			if err := ConvertError(ctx, log, hdl); err != nil {
				errs := err.(customerror.Error)
				if err := web.Respond(ctx, w, errs, codeStatus[errs.Code.Value()]); err != nil {
					return err
				}
			}
			return nil
		}
		return h
	}
	return m
}

// ConvertError handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func ConvertError(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)
	if err == nil {
		return nil
	}

	log.Error(ctx, "message", "ERROR", err.Error())

	if customerror.IsError(err) {
		return customerror.GetError(err)
	}

	return customerror.Newf(customerror.Unknown, customerror.Unknown.String())
}
