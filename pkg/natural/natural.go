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

	x, y := a[:], b[:]

	// Strategy, walk through each segment and check against the other source.
	// Note: that a segments are greedy, so `001` is a segment.
	for {
		// Check to see if the string contains digits at the start of it
		xPos, yPos := strings.IndexFunc(x, unicode.IsDigit), strings.IndexFunc(y, unicode.IsDigit)
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
		xPos, yPos = strings.IndexFunc(x, not(unicode.IsDigit)), strings.IndexFunc(y, not(unicode.IsDigit))
		if xPos == -1 {
			xPos = len(x)
		}
		if yPos == -1 {
			yPos = len(y)
		}

		// Attempt to grab the integer from the string
		// Note: Decimal positioning, because `.` are treated above, then we can use
		// the position of matching values to check for decimal precision.
		xNum, err := strconv.Atoi(x[:xPos])
		if err != nil {
			panic(errors.Wrapf(err, "invalid number %q", x[:xPos]))
		}
		yNum, err := strconv.Atoi(y[:yPos])
		if err != nil {
			panic(errors.Wrapf(err, "invalid number %q", y[:yPos]))
		}

		if xNum != yNum {
			return xNum < yNum
		}

		// Sometimes numbers are not the same `001` vs `1`
		if xPos != yPos {
			return yPos > xPos
		}

		// Continue onwards
		x, y = x[xPos:], y[yPos:]
	}
}

func not(fn func(rune) bool) func(rune) bool {
	return func(r rune) bool {
		return !fn(r)
	}
}
