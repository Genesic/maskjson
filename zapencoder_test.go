package maskjson

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	err := zap.RegisterEncoder("mine", func(config zapcore.EncoderConfig) (zapcore.Encoder, error) {
		jsonEncoder := zapcore.NewJSONEncoder(config)
		maskEncoder := NewZapEncoder(jsonEncoder, false, 3)
		return maskEncoder, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	config := zap.NewDevelopmentConfig()
	// Set the encoder to use our custom encoder
	config.Encoding = "mine"
	log, err := config.Build()
	if err != nil {
		t.Fatal(err)
	}

	logger := log.Sugar()

	msg := message{
		ID: "sdkfjnwkejf",
		Event: secretConfig{
			Host:      "127.0.0.1",
			SecretKey: "abcd1234",
		},
	}

	logger = logger.Named("test").With("message", msg)
	logger.Infow("should be mask", "message2", msg)
}
