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
