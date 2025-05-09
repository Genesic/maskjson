package maskjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Mask struct {
	fullMask    bool
	atLeastStar int
}

func NewMask(fullMask bool, atLeastStar uint) *Mask {
	return &Mask{fullMask, int(atLeastStar)}
}

func (m *Mask) Marshal(v interface{}) ([]byte, error) {
	masked := m.maskStruct(reflect.ValueOf(v))
	return json.Marshal(masked)
}

func (m *Mask) maskStruct(val reflect.Value) interface{} {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if !val.IsValid() {
		return nil
	}

	if val.Kind() != reflect.Struct {
		return val.Interface()
	}

	typ := val.Type()
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// define json key from tag
		parts := strings.Split(jsonTag, ",")
		fieldName := parts[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		omitEmpty := contains(parts[1:], "omitempty")

		if field.Tag.Get("mask") == "true" {
			result[fieldName] = m.maskValue(fieldVal)
			continue
		}

		kind := fieldVal.Kind()
		switch kind {
		case reflect.Struct, reflect.Ptr:
			masked := m.maskStruct(fieldVal)
			if omitEmpty && isEmptyValue(fieldVal) {
				continue
			}
			result[fieldName] = masked

		case reflect.Slice, reflect.Array:
			if omitEmpty && fieldVal.Len() == 0 {
				continue
			}
			slice := make([]interface{}, 0, fieldVal.Len())
			for i := 0; i < fieldVal.Len(); i++ {
				slice = append(slice, m.maskStruct(fieldVal.Index(i)))
			}
			result[fieldName] = slice

		case reflect.Map:
			if omitEmpty && fieldVal.Len() == 0 {
				continue
			}
			mapped := make(map[string]interface{})
			for _, key := range fieldVal.MapKeys() {
				v := fieldVal.MapIndex(key)
				mapped[fmt.Sprint(key.Interface())] = m.maskStruct(v)
			}
			result[fieldName] = mapped

		case reflect.Interface:
			if omitEmpty && fieldVal.IsNil() {
				continue
			}
			inner := fieldVal.Elem()
			if !inner.IsValid() {
				result[fieldName] = nil
				continue
			}

			if omitEmpty && isEmptyValue(inner) {
				continue
			}
			result[fieldName] = m.maskStruct(inner)
		default:
			if omitEmpty && isEmptyValue(fieldVal) {
				continue
			}
			result[fieldName] = fieldVal.Interface()
		}
	}

	return result
}

func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (m *Mask) maskValue(val reflect.Value) interface{} {
	if !val.IsValid() {
		return strings.Repeat("*", m.atLeastStar)
	}

	if m.fullMask {
		return strings.Repeat("*", m.atLeastStar)
	}

	if val.Kind() == reflect.String {
		str := val.String()
		n := len(str)

		if n == 0 {
			return ""
		}

		if n <= m.atLeastStar {
			return strings.Repeat("*", m.atLeastStar)
		}

		visible := (n + m.atLeastStar - 1) / m.atLeastStar
		maskedLen := n - visible
		if maskedLen < m.atLeastStar {
			maskedLen = m.atLeastStar
		}
		return str[:visible] + strings.Repeat("*", maskedLen)
	}

	return strings.Repeat("*", m.atLeastStar)
}
