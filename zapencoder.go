package maskjson

import (
	"encoding/json"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type ZapEncoder struct {
	zapcore.Encoder

	maskEncoder MaskEncoder
}

func NewZapEncoder(encoder zapcore.Encoder, fullMask bool, atLeastStar uint) zapcore.Encoder {
	maskEncoder := NewMask(fullMask, atLeastStar)
	return ZapEncoder{encoder, maskEncoder}
}

func (m ZapEncoder) Clone() zapcore.Encoder {
	return ZapEncoder{m.Encoder.Clone(), m.maskEncoder}
}

func (m ZapEncoder) AddReflected(key string, val interface{}) error {
	data, err := m.maskEncoder.Marshal(val)
	if err != nil {
		return m.Encoder.AddReflected(key, val)
	}

	var masked interface{}
	if err = json.Unmarshal(data, &masked); err != nil {
		return m.Encoder.AddReflected(key, val)
	}

	return m.Encoder.AddReflected(key, masked)
}

func (m ZapEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	maskedFields := make([]zapcore.Field, len(fields))

	for i, field := range fields {
		maskedFields[i] = field

		// Only process fields with complex values that might need masking
		if field.Type == zapcore.ReflectType || field.Type == zapcore.ObjectMarshalerType {
			var val interface{}
			if field.Type == zapcore.ReflectType {
				val = field.Interface
			} else {
				// For object marshalers, we need to convert to a map first
				enc := zapcore.NewMapObjectEncoder()
				if om, ok := field.Interface.(zapcore.ObjectMarshaler); ok {
					if err := om.MarshalLogObject(enc); err == nil {
						val = enc.Fields
					} else {
						continue
					}
				} else {
					continue
				}
			}

			// Try to mask the value
			data, err := m.maskEncoder.Marshal(val)
			if err != nil {
				continue
			}

			var masked interface{}
			if err = json.Unmarshal(data, &masked); err != nil {
				continue
			}

			// Create a new field with the masked value
			maskedFields[i] = zapcore.Field{
				Key:       field.Key,
				Type:      zapcore.ReflectType,
				Interface: masked,
			}
		}
	}

	// Use the base encoder to encode the entry with our masked fields
	return m.Encoder.EncodeEntry(ent, maskedFields)
}
