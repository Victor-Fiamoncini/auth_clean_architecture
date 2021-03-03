package contracts

import (
	shared_custom_errors "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/errors"

	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/types"
)

// IHTTPResponse interface
type IHTTPResponse interface {
	GetStatusCode() int
	SetStatusCode(statusCode int)
	GetErrorObject() error
	SetErrorObject(errorObject error)
	GetErrorName() string
	SetErrorName(errorName string)
	GetBody() types.Map
	SetBody(body types.Map)

	Success(data types.Map) IHTTPResponse
	BadRequest(err shared_custom_errors.IDefaultError) IHTTPResponse
	Unauthorized() IHTTPResponse
	ServerError() IHTTPResponse
}