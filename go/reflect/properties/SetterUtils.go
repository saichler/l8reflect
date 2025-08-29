package properties

import (
	"reflect"
	"strconv"
)

func ConvertValue(target, source reflect.Value) reflect.Value {
	targetKind := target.Kind()
	sourceKind := source.Kind()

	// Handle numeric conversions
	if IsNumeric(targetKind) && IsNumeric(sourceKind) {
		return convertNumeric(target, source)
	}

	// Handle string conversion
	if targetKind == reflect.String {
		return reflect.ValueOf(convertToString(source))
	}

	// Return source unchanged if no conversion needed
	return source
}

func IsNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

func convertNumeric(target, source reflect.Value) reflect.Value {
	targetKind := target.Kind()
	sourceKind := source.Kind()

	// Handle complex numbers separately
	if sourceKind == reflect.Complex64 || sourceKind == reflect.Complex128 {
		if targetKind == reflect.Complex64 || targetKind == reflect.Complex128 {
			c := source.Complex()
			if targetKind == reflect.Complex64 {
				return reflect.ValueOf(complex64(c))
			}
			return reflect.ValueOf(complex128(c))
		}
		// Convert complex to real (use real part)
		numValue := real(source.Complex())
		switch targetKind {
		case reflect.Int:
			return reflect.ValueOf(int(numValue))
		case reflect.Int8:
			return reflect.ValueOf(int8(numValue))
		case reflect.Int16:
			return reflect.ValueOf(int16(numValue))
		case reflect.Int32:
			return reflect.ValueOf(int32(numValue))
		case reflect.Int64:
			return reflect.ValueOf(int64(numValue))
		case reflect.Uint:
			return reflect.ValueOf(uint(numValue))
		case reflect.Uint8:
			return reflect.ValueOf(uint8(numValue))
		case reflect.Uint16:
			return reflect.ValueOf(uint16(numValue))
		case reflect.Uint32:
			return reflect.ValueOf(uint32(numValue))
		case reflect.Uint64:
			return reflect.ValueOf(uint64(numValue))
		case reflect.Float32:
			return reflect.ValueOf(float32(numValue))
		case reflect.Float64:
			return reflect.ValueOf(numValue)
		}
		return source
	}

	// Handle conversion to complex
	if targetKind == reflect.Complex64 || targetKind == reflect.Complex128 {
		var realPart float64
		switch sourceKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			realPart = float64(source.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			realPart = float64(source.Uint())
		case reflect.Float32, reflect.Float64:
			realPart = source.Float()
		default:
			return source
		}
		if targetKind == reflect.Complex64 {
			return reflect.ValueOf(complex64(complex(realPart, 0)))
		}
		return reflect.ValueOf(complex128(complex(realPart, 0)))
	}

	// Get the numeric value as float64 for conversion
	var numValue float64
	switch sourceKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numValue = float64(source.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		numValue = float64(source.Uint())
	case reflect.Float32, reflect.Float64:
		numValue = source.Float()
	default:
		return source
	}

	// Convert to target type with overflow protection
	switch targetKind {
	case reflect.Int:
		// Check for overflow
		if numValue > float64(^uint(0)>>1) || numValue < float64(-1<<63) {
			return reflect.ValueOf(int(0)) // Return zero on overflow
		}
		return reflect.ValueOf(int(numValue))
	case reflect.Int8:
		return reflect.ValueOf(int8(numValue))
	case reflect.Int16:
		return reflect.ValueOf(int16(numValue))
	case reflect.Int32:
		return reflect.ValueOf(int32(numValue))
	case reflect.Int64:
		// Check for overflow
		if numValue > 9223372036854775807 || numValue < -9223372036854775808 {
			return reflect.ValueOf(int64(0)) // Return zero on overflow
		}
		return reflect.ValueOf(int64(numValue))
	case reflect.Uint:
		return reflect.ValueOf(uint(numValue))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(numValue))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(numValue))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(numValue))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(numValue))
	case reflect.Float32:
		return reflect.ValueOf(float32(numValue))
	case reflect.Float64:
		return reflect.ValueOf(numValue)
	}

	return source
}

func convertToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Slice:
		// Handle []byte specially
		if value.Type().Elem().Kind() == reflect.Uint8 {
			return string(value.Interface().([]byte))
		}
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		// Use proper precision for float32
		return strconv.FormatFloat(value.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.Complex64:
		c := value.Complex()
		return strconv.FormatComplex(c, 'g', -1, 64)
	case reflect.Complex128:
		c := value.Complex()
		return strconv.FormatComplex(c, 'g', -1, 128)
	default:
		return value.String()
	}
}
