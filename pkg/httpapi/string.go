package httpapi

import (
	"encoding/json"
	"strings"
)

type String string

func NewString(s string) String {
	return String(s)
}

func (s String) String() string {
	return string(s)
}

func (s *String) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*s = String(strings.TrimSpace(str))

	return nil
}
