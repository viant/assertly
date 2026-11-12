package assertly

import "testing"

func TestFloatNullRequiresExplicitCoalescing(t *testing.T) {
	var nilFloat *float64
	for _, coalesce := range []bool{false, true} {
		for _, expected := range []interface{}{float32(0), float64(0), int(0)} {
			for _, actual := range []interface{}{nil, nilFloat, float64(0)} {
				wantPass := coalesce || actual == float64(0)
				v, err := Assert(map[string]interface{}{"@coalesceWithZero@": coalesce, "value": expected}, map[string]interface{}{"value": actual}, NewDataPath("/"))
				if err != nil {
					t.Fatal(err)
				}
				if (!v.HasFailure()) != wantPass {
					t.Errorf("coalesce=%v expected=%T actual=%#v: %s", coalesce, expected, actual, v.Report())
				}
			}
		}
	}
	for _, coalesce := range []bool{false, true} {
		expected := `{"value":0}`
		if coalesce {
			expected = `{"@coalesceWithZero@":true,"value":0}`
		}
		v, err := Assert(expected, `{"value":null}`, NewDataPath("/"))
		if err != nil {
			t.Fatal(err)
		}
		if (!v.HasFailure()) != coalesce {
			t.Errorf("JSON coalesce=%v: %s", coalesce, v.Report())
		}
	}
}
