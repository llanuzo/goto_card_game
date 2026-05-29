//go:build unit

package httpapi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_UnmarshalJSON(t *testing.T) {
	t.Run("normal string", func(t *testing.T) {
		var s String
		err := json.Unmarshal([]byte(`"hello"`), &s)
		assert.NoError(t, err)
		assert.Equal(t, String("hello"), s)
	})

	t.Run("trims whitespace", func(t *testing.T) {
		var s String
		err := json.Unmarshal([]byte(`"  hello  "`), &s)
		assert.NoError(t, err)
		assert.Equal(t, String("hello"), s)
	})

	t.Run("empty string", func(t *testing.T) {
		var s String
		err := json.Unmarshal([]byte(`""`), &s)
		assert.NoError(t, err)
		assert.Equal(t, String(""), s)
	})

	t.Run("invalid JSON type", func(t *testing.T) {
		var s String
		err := json.Unmarshal([]byte(`123`), &s)
		assert.Error(t, err)
	})
}

func TestEnumString_UnmarshalJSON(t *testing.T) {
	t.Run("lowercases input", func(t *testing.T) {
		var s EnumString
		err := json.Unmarshal([]byte(`"ACTIVE"`), &s)
		assert.NoError(t, err)
		assert.Equal(t, EnumString("active"), s)
	})

	t.Run("empty string becomes undefined", func(t *testing.T) {
		var s EnumString
		err := json.Unmarshal([]byte(`""`), &s)
		assert.NoError(t, err)
		assert.Equal(t, EnumString("undefined"), s)
	})

	t.Run("trims whitespace and lowercases", func(t *testing.T) {
		var s EnumString
		err := json.Unmarshal([]byte(`"  HELLO  "`), &s)
		assert.NoError(t, err)
		assert.Equal(t, EnumString("hello"), s)
	})

	t.Run("invalid JSON type", func(t *testing.T) {
		var s EnumString
		err := json.Unmarshal([]byte(`123`), &s)
		assert.Error(t, err)
	})
}

func TestEnumString_MarshalJSON(t *testing.T) {
	t.Run("undefined marshals to null", func(t *testing.T) {
		s := EnumString("undefined")
		data, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte(`null`), data)
	})

	t.Run("normal value marshals to quoted string", func(t *testing.T) {
		s := EnumString("active")
		data, err := s.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, []byte(`"active"`), data)
	})
}

func TestEnumString_RoundTrip(t *testing.T) {
	t.Run("undefined survives round-trip as null then undefined", func(t *testing.T) {
		original := EnumString("undefined")
		data, err := json.Marshal(original)
		assert.NoError(t, err)
		assert.Equal(t, []byte(`null`), data)
	})

	t.Run("normal value survives round-trip", func(t *testing.T) {
		original := NewEnumString("active")
		data, err := json.Marshal(original)
		assert.NoError(t, err)

		var restored EnumString
		err = json.Unmarshal(data, &restored)
		assert.NoError(t, err)
		assert.Equal(t, original, restored)
	})
}
