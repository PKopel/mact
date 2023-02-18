package utils

import (
	"reflect"
	"testing"
)

func TestUnsnoc(t *testing.T) {
	klmn := "klmn"
	cases := map[string]struct {
		slice []string
		heads []string
		last  *string
	}{
		"happy path": {
			slice: []string{"asdf", "ghij", klmn},
			heads: []string{"asdf", "ghij"},
			last:  &klmn,
		},
		"empty slice": {
			slice: []string{},
			heads: []string{},
			last:  nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			heads, last := Unsnoc(tc.slice)
			if !reflect.DeepEqual(heads, tc.heads) && *last != *tc.last {
				t.Errorf("Expected %v %v, got %v %v", tc.heads, tc.last, heads, last)
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := map[string]struct {
		slice  []string
		value  string
		result bool
	}{
		"happy path": {
			slice:  []string{"asdf", "ghij", "klmn"},
			value:  "asdf",
			result: true,
		},
		"no match": {
			slice:  []string{"asdf", "ghij", "klmn"},
			value:  "xyz",
			result: false,
		},
		"empty slice": {
			slice:  []string{},
			value:  "asdf",
			result: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := Contains(tc.slice, tc.value)
			if result != tc.result {
				t.Errorf("Expected %v, got %v for %v %v", tc.result, result, tc.slice, tc.value)
			}
		})
	}
}
