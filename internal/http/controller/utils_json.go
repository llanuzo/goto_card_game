package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func writeJson(w http.ResponseWriter, code int, data any) error {
	if reflect.TypeOf(data).Kind() != reflect.Pointer {
		return fmt.Errorf("data must be a pointer")
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal %T as json: %w", data, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)

	return nil
}
