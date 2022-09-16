package cli

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rdeusser/cli/constraints"
	"github.com/rdeusser/cli/internal/errors"
)

//go:generate go run tools/gen-type-converter/main.go

// Value represents all possible value types that can be passed to a flag,
// argument, or environment variable.
//
// This is a good demonstration of how generics can be used badly. You can't
// just throw a bunch of random unrelated types as constraints and expect the
// language to handle it well. For a better example of how constraints should be
// used/grouped, see here:
// https://pkg.go.dev/golang.org/x/exp/constraints. Using generics like this
// here is fine because this is limited to a single file in a small
// project. Alternatively, I could've written a generator that could have
// generated all these flag/arg types. Is that a lot more code? Yes, but it's
// generated. For the parsing parts of it you _could_ use generics as intended
// and do something like this for parsing signed types:
//
//	func ParseSigned[T constraints.Signed](s string) (T, error) {
//	    var t T
//
//	    typ := fmt.Sprintf("%T", t)
//	    bitSize := strings.TrimFunc(typ, func(r rune) bool {
//	        return !unicode.IsNumber(r)
//	    })
//
//	    if bitSize == "" {
//	        bitSize = "32"
//	    }
//
//	    i, err := strconv.Atoi(bitSize)
//	    if err != nil {
//	        return t, err
//	    }
//
//	    v, err := strconv.ParseInt(s, 10, i)
//	    if err != nil {
//	        return t, err
//	    }
//
//	    t = T(v)
//
//	    return t, nil
//	}
//
// And then use that in your generated signed types. There's probably a better
// way to write that, but it's a function in a comment so I'm not putting a lot
// of effort into it.
type Value interface {
	constraints.Bool | constraints.Signed | constraints.Unsigned | constraints.Float | constraints.Complex | constraints.Bytes | constraints.String | constraints.Time | constraints.URL | constraints.IP | Path
}

var _ fmt.Stringer = (*Path)(nil)

// Path is a path on the filesystem.
type Path struct {
	Path   string
	Ext    string
	IsDir  bool
	Exists bool
}

// ParsePath takes a path as input, cleans it, and returns Path.
func ParsePath(path string) (Path, error) {
	p := Path{
		Path:   filepath.Clean(path),
		Ext:    filepath.Ext(path),
		IsDir:  false,
		Exists: false,
	}

	file, err := os.Open(p.Path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return p, err
	}
	defer file.Close()

	// File doesn't exist, so return early.
	if file == nil {
		return p, nil
	}

	info, err := os.Stat(file.Name())
	if err != nil {
		return p, err
	}

	p.IsDir = info.IsDir()
	p.Exists = true

	return p, nil
}

// String returns the path provided to Path.
func (p Path) String() string {
	return p.Path
}

// parseValue parses any input value whose string form can be parsed as one of
// the above types.
//
// See comments on the `Value` type to understand this.
func parseValue[T Value](value any, separator byte, layout string) (T, error) {
	var result T

	s := fmt.Sprint(value)
	if s == "" {
		return result, nil
	}

	switch any(&result).(type) {
	case *bool, *[]bool:
		v, err := parseBool[T](s, separator)
		if err != nil {
			return result, err
		}
		result = v
	case *int, *[]int:
		v, err := parseSigned[T](s, separator, 32)
		if err != nil {
			return result, err
		}
		result = v
	case *int8, *[]int8:
		v, err := parseSigned[T](s, separator, 8)
		if err != nil {
			return result, err
		}
		result = v
	case *int16, *[]int16:
		v, err := parseSigned[T](s, separator, 16)
		if err != nil {
			return result, err
		}
		result = v
	case *int32, *[]int32:
		v, err := parseSigned[T](s, separator, 32)
		if err != nil {
			return result, err
		}
		result = v
	case *int64, *[]int64:
		v, err := parseSigned[T](s, separator, 64)
		if err != nil {
			return result, err
		}
		result = v
	case *uint, *[]uint:
		v, err := parseUnsigned[T](s, separator, 32)
		if err != nil {
			return result, err
		}
		result = v
	case *uint8, *[]uint8:
		v, err := parseUnsigned[T](s, separator, 8)
		if err != nil {
			return result, err
		}
		result = v
	case *uint16, *[]uint16:
		v, err := parseUnsigned[T](s, separator, 16)
		if err != nil {
			return result, err
		}
		result = v
	case *uint32, *[]uint32:
		v, err := parseUnsigned[T](s, separator, 32)
		if err != nil {
			return result, err
		}
		result = v
	case *uint64, *[]uint64:
		v, err := parseUnsigned[T](s, separator, 64)
		if err != nil {
			return result, err
		}
		result = v
	case *uintptr, *[]uintptr:
		v, err := parseUnsigned[T](s, separator, 64)
		if err != nil {
			return result, err
		}
		result = v
	case *float32, *[]float32:
		v, err := parseFloat[T](s, separator, 32)
		if err != nil {
			return result, err
		}
		result = v
	case *float64, *[]float64:
		v, err := parseFloat[T](s, separator, 64)
		if err != nil {
			return result, err
		}
		result = v
	case *complex64, *[]complex64:
		v, err := parseComplex[T](s, separator, 64)
		if err != nil {
			return result, err
		}
		result = v
	case *complex128, *[]complex128:
		v, err := parseComplex[T](s, separator, 128)
		if err != nil {
			return result, err
		}
		result = v
	case *[][]byte:
		v, err := parseBytes[T]([]byte(s), separator)
		if err != nil {
			return result, err
		}
		result = v
	case *string, *[]string:
		v, err := parseString[T](s, separator)
		if err != nil {
			return result, err
		}
		result = v
	case *time.Time, *[]time.Time:
		v, err := parseTime[T](s, separator, layout)
		if err != nil {
			return result, err
		}
		result = v
	case *url.URL, *[]url.URL:
		v, err := parseURL[T](s, separator)
		if err != nil {
			return result, err
		}
		result = v
	case *net.IP, *[]net.IP:
		v, err := parseIP[T](s, separator)
		if err != nil {
			return result, err
		}
		result = v
	case *Path, *[]Path:
		v, err := parsePath[T](s, separator)
		if err != nil {
			return result, err
		}
		result = v
	}

	return result, nil
}

func parseBool[T Value](s string, separator byte) (T, error) {
	var result T

	values := make([]bool, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return result, err
		}

		values = append(values, b)
	}

	switch v := any(&result).(type) {
	case *bool:
		*v = values[0]
	case *[]bool:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Bool, got %T", result)
	}

	return result, nil
}

func parseSigned[T Value](s string, separator byte, bitSize int) (T, error) {
	var result T

	values := make([]int64, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		i, err := strconv.ParseInt(v, 10, bitSize)
		if err != nil {
			return result, err
		}

		values = append(values, i)
	}

	switch v := any(&result).(type) {
	case *int:
		*v = convertInt64ToInt(values)[0]
	case *int8:
		*v = convertInt64ToInt8(values)[0]
	case *int16:
		*v = convertInt64ToInt16(values)[0]
	case *int32:
		*v = convertInt64ToInt32(values)[0]
	case *int64:
		*v = values[0]
	case *[]int:
		*v = convertInt64ToInt(values)
	case *[]int8:
		*v = convertInt64ToInt8(values)
	case *[]int16:
		*v = convertInt64ToInt16(values)
	case *[]int32:
		*v = convertInt64ToInt32(values)
	case *[]int64:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Signed, got %T", result)
	}

	return result, nil
}

func parseUnsigned[T Value](s string, separator byte, bitSize int) (T, error) {
	var result T

	values := make([]uint64, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		i, err := strconv.ParseUint(v, 10, bitSize)
		if err != nil {
			return result, err
		}

		values = append(values, i)
	}

	switch v := any(&result).(type) {
	case *uint:
		*v = convertUint64ToUint(values)[0]
	case *uint8:
		*v = convertUint64ToUint8(values)[0]
	case *uint16:
		*v = convertUint64ToUint16(values)[0]
	case *uint32:
		*v = convertUint64ToUint32(values)[0]
	case *uint64:
		*v = values[0]
	case *uintptr:
		*v = convertUint64ToUintptr(values)[0]
	case *[]uint:
		*v = convertUint64ToUint(values)
	case *[]uint8:
		*v = convertUint64ToUint8(values)
	case *[]uint16:
		*v = convertUint64ToUint16(values)
	case *[]uint32:
		*v = convertUint64ToUint32(values)
	case *[]uint64:
		*v = values
	case *[]uintptr:
		*v = convertUint64ToUintptr(values)
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Unsigned, got %T", result)
	}

	return result, nil
}

func parseFloat[T Value](s string, separator byte, bitSize int) (T, error) {
	var result T

	values := make([]float64, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		b, err := strconv.ParseFloat(v, bitSize)
		if err != nil {
			return result, err
		}

		values = append(values, b)
	}

	switch v := any(&result).(type) {
	case *float32:
		*v = convertFloat64ToFloat32(values)[0]
	case *float64:
		*v = values[0]
	case *[]float32:
		*v = convertFloat64ToFloat32(values)
	case *[]float64:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Float, got %T", result)
	}

	return result, nil
}

func parseComplex[T Value](s string, separator byte, bitSize int) (T, error) {
	var result T

	values := make([]complex128, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		b, err := strconv.ParseComplex(v, bitSize)
		if err != nil {
			return result, err
		}

		values = append(values, b)
	}

	switch v := any(&result).(type) {
	case *complex64:
		*v = convertComplex128ToComplex64(values)[0]
	case *complex128:
		*v = values[0]
	case *[]complex64:
		*v = convertComplex128ToComplex64(values)
	case *[]complex128:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Complex, got %T", result)
	}

	return result, nil
}

func parseBytes[T Value](s []byte, separator byte) (T, error) {
	var result T

	values := make([][]byte, 0)
	slice := splitBytes(s, separator)

	values = append(values, slice...)

	switch v := any(&result).(type) {
	case *[]byte:
		*v = values[0]
	case *[][]byte:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Bytes, got %T", result)
	}

	return result, nil
}

func parseString[T Value](s string, separator byte) (T, error) {
	var result T

	values := make([]string, 0)
	slice := splitString(s, separator)

	values = append(values, slice...)

	switch v := any(&result).(type) {
	case *string:
		*v = values[0]
	case *[]string:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.String, got %T", result)
	}

	return result, nil
}

func parseTime[T Value](s string, separator byte, layout string) (T, error) {
	var result T

	values := make([]time.Time, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		t, err := time.Parse(layout, v)
		if err != nil {
			return result, err
		}

		values = append(values, t)
	}

	switch v := any(&result).(type) {
	case *time.Time:
		*v = values[0]
	case *[]time.Time:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.Time, got %T", result)
	}

	return result, nil
}

func parseURL[T Value](s string, separator byte) (T, error) {
	var result T

	values := make([]url.URL, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		u, err := url.Parse(v)
		if err != nil {
			return result, err
		}

		values = append(values, *u)
	}

	switch v := any(&result).(type) {
	case *url.URL:
		*v = values[0]
	case *[]url.URL:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.URL, got %T", result)
	}

	return result, nil
}

func parseIP[T Value](s string, separator byte) (T, error) {
	var result T

	values := make([]net.IP, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		values = append(values, net.ParseIP(v))
	}

	switch v := any(&result).(type) {
	case *net.IP:
		*v = values[0]
	case *[]net.IP:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be constrained by constraints.IP, got %T", result)
	}

	return result, nil
}

func parsePath[T Value](s string, separator byte) (T, error) {
	var result T

	values := make([]Path, 0)
	slice := splitString(s, separator)

	for _, v := range slice {
		p, err := ParsePath(v)
		if err != nil {
			return result, err
		}

		values = append(values, p)
	}

	switch v := any(&result).(type) {
	case *Path:
		*v = values[0]
	case *[]Path:
		*v = values
	default:
		return result, fmt.Errorf("expected type to be Path, got %T", result)
	}

	return result, nil
}

func isZeroValue[T any](value T) bool {
	switch v := any(value).(type) {
	case bool:
		return v == *new(bool)
	case int:
		return v == *new(int)
	case int8:
		return v == *new(int8)
	case int16:
		return v == *new(int16)
	case int32:
		return v == *new(int32)
	case int64:
		return v == *new(int64)
	case uint:
		return v == *new(uint)
	case uint8:
		return v == *new(uint8)
	case uint16:
		return v == *new(uint16)
	case uint32:
		return v == *new(uint32)
	case uint64:
		return v == *new(uint64)
	case uintptr:
		return v == *new(uintptr)
	case float32:
		return v == *new(float32)
	case float64:
		return v == *new(float64)
	case complex64:
		return v == *new(complex64)
	case complex128:
		return v == *new(complex128)
	case string:
		return v == *new(string)
	case time.Time:
		return v.Equal(*new(time.Time))
	case url.URL:
		return v == *new(url.URL)
	case Path:
		return v == *new(Path)
	default:
		return false
	}
}
