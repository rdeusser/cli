package cli

// Value mimics flag.Value to avoid a dependency on the standard flag
// library. Also helps with name collision.
type Value interface {
	String() string
	Set(string) error
}
