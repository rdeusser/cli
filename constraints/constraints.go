package constraints

import (
	"net"
	"net/url"
	"time"
)

// Bool is a constraint for boolean types.
type Bool interface {
	~bool | ~[]bool
}

// Signed is a constraint for signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~[]int | ~[]int8 | ~[]int16 | ~[]int32 | ~[]int64
}

// Unsigned is a constraint for unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 | ~[]uintptr
}

// Float is a constraint for floating-point types.
type Float interface {
	~float32 | ~float64 | ~[]float32 | ~[]float64
}

// Complex is a constraint for complex numeric types.
type Complex interface {
	~complex64 | ~complex128 | ~[]complex64 | ~[]complex128
}

// Bytes is a constraint for []byte types.
type Bytes interface {
	~[]byte | ~[][]byte
}

// String is a constraint for string types.
type String interface {
	~string | ~[]string
}

// Time is a constraint for time types.
type Time interface {
	time.Time | ~[]time.Time
}

// URL is a constraint for url types.
type URL interface {
	url.URL | ~[]url.URL
}

// IP is a constraint for ip types.
type IP interface {
	net.IP | ~[]net.IP
}
