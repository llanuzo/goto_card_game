package controller

import (
	"fmt"
)

type ErrApiResponse struct {
	HttpStatusCode int
	Message        string
}

func newErrApiResponse(httpStatusCode int, format string, args ...any) ErrApiResponse {
	return ErrApiResponse{
		HttpStatusCode: httpStatusCode,
		Message:        fmt.Sprintf(format, args...),
	}
}

func (e ErrApiResponse) Error() string {
	return e.Message
}
