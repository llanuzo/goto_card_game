package httpapi

type ErrorResponse struct {
	Message string `json:"message"`

	StatusCode int `json:"-"`
}
