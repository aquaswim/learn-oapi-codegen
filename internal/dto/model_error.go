package dto

func (e *Error) Error() string {
	return e.Message
}
