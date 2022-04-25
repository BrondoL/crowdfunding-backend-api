package helper

import (
	"errors"
	"io"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  status,
		},
		Data: data,
	}
}

func FormatValidationError(err error) []string {
	var errorMessage []string

	if !errors.Is(err, io.EOF) {
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Field()+" "+e.Tag())
		}
		return errorMessage
	}
	return []string{"No Data Found"}
}
