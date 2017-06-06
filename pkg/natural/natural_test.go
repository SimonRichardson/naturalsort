package natural

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		actual, expected []string
	}{
		{
			"empty right order",
			[]string{"", "a"},
			[]string{"", "a"},
		},
		{
			"empty left order",
			[]string{"a", ""},
			[]string{"", "a"},
		},
		{
			"differing lengths",
			[]string{"aac", "ac", "a"},
			[]string{"a", "aac", "ac"},
		},
		{
			"similar lengths",
			[]string{"aa", "cc", "bb"},
			[]string{"aa", "bb", "cc"},
		},
		{
			"similar digit lengths",
			[]string{"11", "33", "22"},
			[]string{"11", "22", "33"},
		},
		{
			"digit order",
			[]string{"b11", "a2"},
			[]string{"a2", "b11"},
		},
		{
			"human digit order",
			[]string{"z11", "z2"},
			[]string{"z2", "z11"},
		},
		{
			"alpha numeric",
			[]string{"a1", "a0", "a13", "a11", "a99", "a11", "a2"},
			[]string{"a0", "a1", "a2", "a11", "a11", "a13", "a99"},
		},
		{
			"numeric",
			[]string{"001", "2", "30", "22", "0", "00", "3", "1"},
			[]string{"0", "00", "1", "001", "2", "3", "22", "30"},
		},
		{
			"glyphs",
			[]string{"世界3", "世20"},
			[]string{"世20", "世界3"},
		},
		{
			"numeric padding",
			[]string{"001", "1"},
			[]string{"1", "001"},
		},
		{
			"decimal",
			[]string{"1.002", "1.001", "1.003"},
			[]string{"1.001", "1.002", "1.003"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Sort(tc.actual)

			if !reflect.DeepEqual(tc.actual, tc.expected) {
				t.Errorf("expected: %v, actual: %v", tc.expected, tc.actual)
			}
		})
	}
}
