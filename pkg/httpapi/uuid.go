package httpapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func NewUuid(id uuid.UUID) UUID {
	return UUID{
		UUID: id,
	}
}

func (t *UUID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	str = strings.TrimSpace(str)

	var err error
	t.UUID, err = uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("invalid uuid input %s in json payload", str)
	}

	return nil
}

func (m UUID) MarshalJSON() ([]byte, error) {
	if m.UUID == uuid.Nil {
		return []byte(`null`), nil
	}

	return fmt.Appendf(nil, "%q", m.UUID.String()), nil
}
