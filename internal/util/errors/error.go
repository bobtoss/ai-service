package errors

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	ErrorCode string `json:"ErrorCode"`
	ErrorDesc string `json:"ErrorDesc"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorDesc
}

func NewCustomErrorResponse(code int, msg string) *echo.HTTPError {
	if len(msg) == 0 {
		msg = http.StatusText(code)
	}

	err := echo.NewHTTPError(code, ErrorResponse{
		ErrorCode: strconv.Itoa(code),
		ErrorDesc: msg,
	})

	return err.SetInternal(errors.WithStack(errors.New(msg)))
}

func NewInternalErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorResponse(http.StatusInternalServerError, msg)
}

func NewBadRequestErrorRsp(msg string) *echo.HTTPError {
	return NewCustomErrorResponse(http.StatusBadRequest, msg)
}
