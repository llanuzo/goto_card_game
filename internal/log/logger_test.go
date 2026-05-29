package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	t.Run("all default fields are set when logger", func(t *testing.T) {

		loggerName := "test-logger"
		logFormat := "test-log: %s"
		logFormatInput := uuid.NewString()

		var buf bytes.Buffer
		logger := NewLogger(loggerName, WithWriter(&buf))

		logger.Infof(logFormat, logFormatInput)

		got := buf.Bytes()

		jsonLog := struct {
			Time       time.Time      `json:"time"`
			Level      string         `json:"level"`
			Source     map[string]any `json:"source"`
			Msg        string         `json:"msg"`
			App        string         `json:"app"`
			LoggerName string         `json:"logger_name"`
		}{}

		err := json.Unmarshal(got, &jsonLog)
		require.NoError(t, err)

		assert.WithinDuration(t, time.Now(), jsonLog.Time, 100*time.Millisecond)
		assert.Equal(t, "INFO", jsonLog.Level)
		assert.NotEmpty(t, jsonLog.Source)
		assert.Equal(t, "card-game", jsonLog.App)
		assert.Equal(t, loggerName, jsonLog.LoggerName)
		assert.Equal(t, fmt.Sprintf(logFormat, logFormatInput), jsonLog.Msg)
	})

	t.Run("default log level is info", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger("test-logger", WithWriter(&buf))

		logger.Debugf("DEBUG")
		logger.Infof("INFO")
		logger.Warnf("WARN")
		logger.Errorf("ERROR")
		got := buf.String()

		assert.NotContains(t, got, "DEBUG")
		assert.Contains(t, got, "INFO")
		assert.Contains(t, got, "WARN")
		assert.Contains(t, got, "ERROR")
	})
}
