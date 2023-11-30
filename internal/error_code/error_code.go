package error_code

import "todo-codegen/internal/dto"

func NewError(code, message string) *dto.Error {
	return &dto.Error{
		Code:    code,
		Message: message,
	}
}

func NewErrorWithData(code, message string, data interface{}) *dto.Error {
	return &dto.Error{
		Code:    code,
		Message: message,
		Detail:  &data,
	}
}

const (
	ErrCodeInternal   = "001"
	ErrCodeNotFound   = "002"
	ErrCodeValidation = "003"
)

var (
	ErrorTodoNotFound = NewError(ErrCodeNotFound, "Todo not found")
)
