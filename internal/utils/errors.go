package utils

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(
	code string,
	message string,
	status int,
) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}
