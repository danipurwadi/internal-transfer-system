package tests

import "github.com/danipurwadi/internal-transfer-system/foundation/customerror"

func toErrorPtr(err customerror.Error) *customerror.Error {
	return &err
}
