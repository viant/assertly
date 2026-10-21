package assertly

import "testing"

func TestJSONNullRequiresNilIncludingTypedSQLPointers(t *testing.T) {
	var integer *int
	var number *float64
	var text *string
	var object *struct{ Value int }
	zeroInt, oneInt := 0, 1
	zeroFloat, oneFloat := 0., 1.
	emptyText, valueText := "", "value"
	cases := []struct {
		name  string
		value interface{}
		pass  bool
	}{
		{"untyped nil", nil, true}, {"SQL integer NULL", integer, true}, {"SQL float NULL", number, true}, {"SQL text NULL", text, true},
		{"zero integer pointer", &zeroInt, false}, {"nonzero integer pointer", &oneInt, false},
		{"zero float pointer", &zeroFloat, false}, {"nonzero float pointer", &oneFloat, false},
		{"nil object retains coalescing contract", object, false},
		{"empty text pointer", &emptyText, false}, {"nonempty text pointer", &valueText, false},
		{"zero integer", 0, false}, {"zero float", 0., false}, {"empty string", "", false},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			validation, err := Assert(map[string]interface{}{"FREQ_CAPPING": nil}, map[string]interface{}{"FREQ_CAPPING": test.value}, NewDataPath("SQL NULL"))
			if err != nil {
				t.Fatal(err)
			}
			if (len(validation.Failures) == 0) != test.pass {
				t.Fatalf("passes=%d failures=%v expected pass=%v", validation.PassedCount, validation.Failures, test.pass)
			}
		})
	}
}
