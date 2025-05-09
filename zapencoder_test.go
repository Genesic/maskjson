package maskjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

type message struct {
	ID    string
	Event interface{}
}

type secretConfig struct {
	Host      string `json:"host"`
	SecretKey string `json:"secret_key" mask:"true"`
}

func TestZapEncoder_AddReflected(t *testing.T) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	maskEncoder := NewZapEncoder(jsonEncoder, false, 3)

	var buf bytes.Buffer
	core := zapcore.NewCore(
		maskEncoder,
		zapcore.AddSync(&buf), // Write to buffer
		zapcore.InfoLevel,
	)

	log := zap.New(core, zap.AddCaller())
	logger := log.Sugar()

	// Create test data
	msg := message{
		ID: "sdkfjnwkejf",
		Event: secretConfig{
			Host:      "127.0.0.1",
			SecretKey: "abcd1234",
		},
	}

	// Log with the message attached
	bs, _ := json.Marshal(msg)
	logger = logger.Named("test").With("message", msg)
	logger.Infow("should be mask", "message2", msg)
	fmt.Println("original message:", string(bs))

	// Get the log output as a string
	output := buf.String()
	t.Log("Log output:", output)

	// Verify the output contains masked values
	if !strings.Contains(output, `"secret_key":"abc*****"`) {
		t.Errorf("Expected masked secret_key in output, got: %s", output)
	}

	// For more detailed verification, parse the JSON and check specific fields
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	// Check that message field is properly masked
	messageObj, ok := logEntry["message"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected message field to be a map, got: %T", logEntry["message"])
	}

	event, ok := messageObj["Event"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected message.Event field to be a map, got: %T", messageObj["Event"])
	}

	secretKey, ok := event["secret_key"].(string)
	if !ok {
		t.Fatalf("Expected message.Event.secret_key field to be a string, got: %T", event["secret_key"])
	}

	if secretKey != "abc*****" {
		t.Errorf("Expected message.Event.secret_key to be masked as 'abc*****', got: %s", secretKey)
	}

	// Check that message2 field is also properly masked
	message2, ok := logEntry["message2"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected message2 field to be a map, got: %T", logEntry["message2"])
	}

	event2, ok := message2["Event"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected message2.Event field to be a map, got: %T", message2["Event"])
	}

	secretKey2, ok := event2["secret_key"].(string)
	if !ok {
		t.Fatalf("Expected message2.Event.secret_key field to be a string, got: %T", event2["secret_key"])
	}

	if secretKey2 != "abc*****" {
		t.Errorf("Expected message2.Event.secret_key to be masked as 'abc*****', got: %s", secretKey2)
	}

}
