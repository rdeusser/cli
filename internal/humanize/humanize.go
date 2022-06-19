package humanize

import "strconv"

// Copied from https://github.com/dustin/go-humanize/blob/v1.0.0/ordinals.go#L8

// Ordinal gives you the input number in a rank/ordinal format.
//
// Ordinal(3) -> 3rd
func Ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}
