//go:build unit

package httpapi

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testUUIDStr = "550e8400-e29b-41d4-a716-446655440000"

func TestUUID_UnmarshalJSON(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		var u UUID
		err := json.Unmarshal([]byte(`"`+testUUIDStr+`"`), &u)
		assert.NoError(t, err)
		assert.Equal(t, uuid.MustParse(testUUIDStr), u.UUID)
	})

	t.Run("trims whitespace", func(t *testing.T) {
		var u UUID
		err := json.Unmarshal([]byte(`"  `+testUUIDStr+`  "`), &u)
		assert.NoError(t, err)
		assert.Equal(t, uuid.MustParse(testUUIDStr), u.UUID)
	})

	t.Run("invalid UUID returns error", func(t *testing.T) {
		var u UUID
		err := json.Unmarshal([]byte(`"not-a-uuid"`), &u)
		assert.Error(t, err)
	})

	t.Run("invalid JSON type returns error", func(t *testing.T) {
		var u UUID
		err := json.Unmarshal([]byte(`123`), &u)
		assert.Error(t, err)
	})
}

func TestUUID_MarshalJSON(t *testing.T) {
	t.Run("nil UUID marshals to null", func(t *testing.T) {
		u := UUID{}
		data, err := u.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte(`null`), data)
	})

	t.Run("valid UUID marshals to quoted string", func(t *testing.T) {
		id := uuid.MustParse(testUUIDStr)
		u := NewUuid(id)
		data, err := u.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte(`"`+testUUIDStr+`"`), data)
	})
}

func TestUUID_RoundTrip(t *testing.T) {
	original := uuid.MustParse(testUUIDStr)
	u := NewUuid(original)

	data, err := json.Marshal(u)
	assert.NoError(t, err)

	var restored UUID
	err = json.Unmarshal(data, &restored)
	assert.NoError(t, err)
	assert.Equal(t, original, restored.UUID)
}
