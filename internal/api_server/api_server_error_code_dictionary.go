package api_server

import (
	"net/http"
	errorCode "todo-codegen/internal/error_code"
)

var errorCodeToHttpCodeDict = map[string]int{
	errorCode.ErrCodeInternal:   http.StatusInternalServerError,
	errorCode.ErrCodeNotFound:   http.StatusNotFound,
	errorCode.ErrCodeValidation: http.StatusBadRequest,
}
