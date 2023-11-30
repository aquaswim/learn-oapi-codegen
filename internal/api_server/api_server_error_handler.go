package api_server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo-codegen/internal/dto"
	errorCode "todo-codegen/internal/error_code"
)

// ErrorHandler will translate the business logic error to echo error
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// For valid credentials call next
		err := next(c)
		if err == nil {
			return err
		}
		c.Logger().Infof("Error handler called with err %+v", err)
		switch e := err.(type) {
		case *dto.Error:
			return appErrorToEchoErrorConverter(e)
		case *echo.HTTPError:
			e.Message = dto.Error{
				Code:    errorCode.ErrCodeInternal,
				Message: e.Unwrap().Error(),
			}
			return e
		default:
			// this default case should be the last one expected by this
			return echo.NewHTTPError(http.StatusInternalServerError, dto.N500{
				Code:    errorCode.ErrCodeInternal,
				Message: fmt.Sprintf("%+v", err),
			})
		}
	}
}

func appErrorToEchoErrorConverter(d *dto.Error) *echo.HTTPError {
	httpCode := http.StatusInternalServerError
	if val, found := errorCodeToHttpCodeDict[d.Code]; found {
		httpCode = val
	}

	return echo.NewHTTPError(httpCode, *d)
}
