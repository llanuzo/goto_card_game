package httpapi

import (
	"encoding/json"
	"fmt"
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

type EnumString string

func NewEnumString(s string) EnumString {
	return EnumString(s)
}

func (s EnumString) String() string {
	return string(s)
}

func (s *EnumString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		*s = "undefined"
	} else {
		*s = EnumString(strings.TrimSpace(strings.ToLower(str)))
	}

	return nil
}

func (s EnumString) MarshalJSON() ([]byte, error) {
	if s == "undefined" {
		return []byte(`null`), nil
	}

	return fmt.Appendf(nil, "%q", s), nil
}
