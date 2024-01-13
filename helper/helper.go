package helper

import "github.com/go-playground/validator/v10"

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseContainer(msg string, code int, status string, data interface{}) Response {
	responseJson := Response{
		Message: msg,
		Code:    code,
		Status:  status,
		Data:    data,
	}

	return responseJson
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
