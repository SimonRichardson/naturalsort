package natural

import (
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// Sort sorts input strings into a more human representation, for example in
// natural sorting `z11` should go _after_ `z2`, because `2 < 11`
func Sort(input []string) {
	sort.Sort(slice(input))
}

type slice []string

func (s slice) Len() int {
	return len(s)
}

func (s slice) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

func (s slice) Less(a, b int) bool {
	return compare(s[a], s[b])
}

func compare(a, b string) bool {
	// Quick check to see if the length of a is empty and b has a value or the
	// inverse.
	if aLen, bLen := len(a), len(b); aLen == 0 && bLen > 0 {
		return true
	} else if bLen == 0 && aLen > 0 {
		return false
	}

	// Make sure we don't call mutations on the original
	x, y := a[:], b[:]

	// Strategy, walk through each segment and check against the other source.
	// Note: that a segments are greedy, so `001` is a segment and will be
	// converted to `1` using `strconv.Atoi`.
	for {
		// Check to see if the string contains digits at the start of it
		xPos, yPos := indexOfNumber(x), indexOfNumber(y)
		if xPos == -1 && yPos == -1 {
			return x < y
		} else if xPos == -1 && yPos >= 0 {
			return false
		} else if yPos == -1 {
			return true
		}

		// Compare actual segments (seg)
		if xSeg, ySeg := x[:xPos], y[:yPos]; xSeg != ySeg {
			return xSeg < ySeg
		}

		// Move on past the non-digit
		x, y = x[xPos:], y[yPos:]
		xPos, yPos = indexOfNonNumber(x), indexOfNonNumber(y)
		if xPos == -1 {
			xPos = len(x)
		}
		if yPos == -1 {
			yPos = len(y)
		}

		// Attempt to grab the integer from the string
		// Note: Decimal positioning, because `.` are treated above, then we can use
		// the position of matching values to check for decimal precision.
		if xNum, yNum := coerceToInt(x[:xPos]), coerceToInt(y[:yPos]); xNum != yNum {
			return xNum < yNum
		}

		// Sometimes numbers are not the same `001` vs `1` so rank them
		// accordingly. Larger values (positions) will get put lastly.
		if xPos != yPos {
			return yPos > xPos
		}

		// Continue onwards
		x, y = x[xPos:], y[yPos:]
	}
}

func indexOfNumber(s string) int {
	return strings.IndexFunc(s, unicode.IsDigit)
}

func indexOfNonNumber(s string) int {
	// We just want the inverse of the `IsDigit` code, shame `unicode` doesn't
	// offer one.
	return strings.IndexFunc(s, not(unicode.IsDigit))
}

// coerceToInt converts a string value to a integer
// Note: this panic's because sorting doesn't understand failures.
func coerceToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(errors.Wrapf(err, "invalid number %q", s))
	}
	return n
}

func not(fn func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !fn(r)
	}
}
