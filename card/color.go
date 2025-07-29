package card

import "slices"

type Colors []string

// NewColor returns a color object representing one of [ "white", "blue", "black", "red", "green" ]
func NewColor(color string) Colors {
	var c Colors
	switch color {
	case "white":
		return Colors{"W"}
	case "blue":
		return Colors{"U"}
	case "black":
		return Colors{"B"}
	case "red":
		return Colors{"R"}
	case "green":
		return Colors{"G"}
	}
	return c
}

// Add appends all elements of other not already contained in s
func (s *Colors) Add(other Colors) {
	for _, c := range other {
		if !slices.Contains(*s, c) {
			*s = append(*s, c)
		}
	}
}

// Equal returns whether s and other contain exactly the same set of colors
func (s Colors) Equal(other Colors) bool {
	return slices.Equal(s, other)
}

// IsSubset returns whether s is a subset of other
func (s Colors) IsSubset(other Colors) bool {
	for _, c := range s {
		if !slices.Contains(other, c) {
			return false
		}
	}
	return true
}
