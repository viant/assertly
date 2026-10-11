package assertly_test

import (
	"testing"

	"github.com/viant/assertly"
)

func TestSortTextValidatesActualValues(t *testing.T) {
	cases := []struct {
		name   string
		actual []interface{}
		valid  bool
	}{
		{"reordered duplicates", []interface{}{"b", "a", "a"}, true},
		{"changed value", []interface{}{"a", "a", "c"}, false},
		{"changed duplicate count", []interface{}{"a", "b", "b"}, false},
		{"missing value", []interface{}{"a", "b"}, false},
		{"extra duplicate", []interface{}{"a", "a", "b", "b"}, false},
		{"empty actual", []interface{}{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected := map[string]interface{}{
				"@length@segments": 3,
				"segments":         []interface{}{map[string]interface{}{"@sortText@": true}, "a", "a", "b"},
			}
			actual := map[string]interface{}{"segments": tc.actual}
			validation, err := assertly.Assert(expected, actual, assertly.NewDataPath("/"))
			if err != nil {
				t.Fatal(err)
			}
			if valid := validation.FailedCount == 0; valid != tc.valid {
				t.Fatalf("valid = %v, want %v: %s", valid, tc.valid, validation.Report())
			}
		})
	}
}
