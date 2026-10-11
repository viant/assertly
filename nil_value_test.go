package assertly_test

import (
	"testing"
	"time"

	"github.com/viant/assertly"
)

func TestAssertRejectsNilExpectedAgainstPresentValues(t *testing.T) {
	for _, actual := range []interface{}{true, false, "2026-10-11T20:00:00Z", time.Date(2026, 10, 11, 20, 0, 0, 0, time.UTC)} {
		t.Run("present value", func(t *testing.T) {
			validation, err := assertly.Assert(map[string]interface{}{"@strictMapCheck@": true, "value": nil}, map[string]interface{}{"value": actual}, assertly.NewDataPath("/"))
			if err != nil || validation.FailedCount != 1 {
				t.Fatalf("expected one value failure, got %v: %#v", err, validation)
			}
		})
	}
}

func TestAssertRejectsExtraObjectWithBoolean(t *testing.T) {
	expected := map[string]interface{}{"@strictMapCheck@": true, "@length@records": 1, "records": []interface{}{map[string]interface{}{"id": 1, "enabled": true}}}
	actual := map[string]interface{}{"records": []interface{}{map[string]interface{}{"id": 1, "enabled": true}, map[string]interface{}{"id": 1, "enabled": true}}}
	validation, err := assertly.Assert(expected, actual, assertly.NewDataPath("/"))
	if err != nil || validation.FailedCount == 0 {
		t.Fatalf("extra object was not rejected: %v: %#v", err, validation)
	}
}

func TestAssertEmptyExpectedTimeDoesNotPanic(t *testing.T) {
	validation, err := assertly.Assert("", time.Date(2026, 10, 11, 20, 0, 0, 0, time.UTC), assertly.NewDataPath("/"))
	if err != nil || validation.FailedCount != 1 {
		t.Fatalf("expected one time failure, got %v: %#v", err, validation)
	}
}
