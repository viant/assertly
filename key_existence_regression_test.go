package assertly_test

import (
	"fmt"
	"testing"

	"github.com/viant/assertly"
)

func TestBareKeyExistenceDirectivesUseTheirValues(t *testing.T) {
	present := map[string]interface{}{"entity": nil, "status": "ok"}
	absent := map[string]interface{}{"status": "ok"}
	assertExistenceCounts(t, `{"@exists@":"entity"}`, present, 1, 0)
	assertExistenceCounts(t, `{"@exists@":"entity"}`, absent, 0, 1)
	assertExistenceCounts(t, `{"@!exists@":"entity"}`, present, 0, 1)
	assertExistenceCounts(t, `{"@!exists@":"entity"}`, absent, 1, 0)
	assertExistenceCounts(t, `{"@exists@":["entity","status"]}`, absent, 1, 1)
	assertExistenceCounts(t, `{"@!exists@":["entity","status"]}`, absent, 1, 1)
}

func TestKeyExistenceDirectiveWithoutPlaceholder(t *testing.T) {
	forms := []struct {
		key    string
		value  bool
		exists bool
	}{
		{"@exists@target", true, true},
		{"@exists@target", false, false},
		{"@!exists@target", true, false},
		{"@!exists@target", false, true},
	}
	actuals := []struct {
		name    string
		present bool
		value   interface{}
	}{
		{"missing", false, nil},
		{"value", true, "present"},
		{"null", true, nil},
	}
	for _, form := range forms {
		for _, actual := range actuals {
			for _, strict := range []bool{false, true} {
				for _, placeholder := range []bool{false, true} {
					name := fmt.Sprintf("%s=%t/%s/strict=%t/placeholder=%t", form.key, form.value, actual.name, strict, placeholder)
					t.Run(name, func(t *testing.T) {
						expected := map[string]interface{}{form.key: form.value, "@strictMapCheck@": strict}
						if placeholder {
							expected["target"] = nil
						}
						observed := map[string]interface{}{}
						if actual.present {
							observed["target"] = actual.value
						}
						passed, failed := 0, 1
						if form.exists == actual.present {
							passed, failed = 1, 0
						}
						assertExistenceCounts(t, expected, observed, passed, failed)
					})
				}
			}
		}
	}
}

func TestKeyExistenceDirectivePreservesStrictMapCheck(t *testing.T) {
	for _, strict := range []bool{false, true} {
		t.Run(fmt.Sprintf("strict=%t", strict), func(t *testing.T) {
			expected := map[string]interface{}{"@exists@target": true, "@strictMapCheck@": strict}
			actual := map[string]interface{}{"target": nil, "undeclared": 42}
			failed := 0
			if strict {
				failed = 1
			}
			assertExistenceCounts(t, expected, actual, 1, failed)
		})
	}
}

func TestKeyExistenceDirectiveNestedResponse(t *testing.T) {
	// Absence is membership, not a comparison with null: even a null messages
	// property must fail this response expectation.
	for _, present := range []bool{false, true} {
		t.Run(fmt.Sprintf("messagesPresent=%t", present), func(t *testing.T) {
			expected := map[string]interface{}{"data": map[string]interface{}{
				"@exists@messages": false, "id": 7,
			}}
			row := map[string]interface{}{"id": 7}
			if present {
				row["messages"] = nil
			}
			passed, failed := 2, 0
			if present {
				passed, failed = 1, 1
			}
			assertExistenceCounts(t, expected, map[string]interface{}{"data": row}, passed, failed)
		})
	}
}

func TestKeyExistenceValueDirectivesRemainSupported(t *testing.T) {
	assertExistenceCounts(t, map[string]interface{}{"target": "@exists@"}, map[string]interface{}{"target": nil}, 1, 0)
	assertExistenceCounts(t, map[string]interface{}{"target": "@exists@"}, map[string]interface{}{}, 0, 1)
	assertExistenceCounts(t, map[string]interface{}{"target": "@!exists@"}, map[string]interface{}{}, 1, 0)
	assertExistenceCounts(t, map[string]interface{}{"target": "@!exists@"}, map[string]interface{}{"target": nil}, 0, 1)
}

func assertExistenceCounts(t *testing.T, expected, actual interface{}, passed, failed int) {
	t.Helper()
	validation, err := assertly.Assert(expected, actual, assertly.NewDataPath("/"))
	if err != nil {
		t.Fatalf("unexpected assertion error: %v", err)
	}
	if validation.PassedCount != passed || validation.FailedCount != failed {
		t.Fatalf("counts passed=%d failed=%d; want passed=%d failed=%d: %s", validation.PassedCount, validation.FailedCount, passed, failed, validation.Report())
	}
}

func TestKeyExistenceDirectiveDoesNotExpandIntoChild(t *testing.T) {
	expected := map[string]interface{}{"@exists@landing_page_url": true, "size": map[string]interface{}{"width": 300}}
	actual := map[string]interface{}{"landing_page_url": nil, "size": map[string]interface{}{"width": 300}}
	assertExistenceCounts(t, expected, actual, 2, 0)
}

func TestKeyExistenceDirectiveNestedLocalTargets(t *testing.T) {
	for _, parentPresent := range []bool{false, true} {
		for _, childPresent := range []bool{false, true} {
			expected := map[string]interface{}{"@exists@parent": true, "child": map[string]interface{}{"@exists@own": true, "width": 300}}
			child := map[string]interface{}{"width": 300}
			actual := map[string]interface{}{"child": child}
			passed, failed := 1, 2
			if parentPresent {
				actual["parent"] = nil
				passed++
				failed--
			}
			if childPresent {
				child["own"] = nil
				passed++
				failed--
			}
			assertExistenceCounts(t, expected, actual, passed, failed)
		}
	}
	assertExistenceCounts(t, map[string]interface{}{"@!exists@parent": true, "child": map[string]interface{}{"width": 300}}, map[string]interface{}{"child": map[string]interface{}{"width": 300}}, 2, 0)
}
