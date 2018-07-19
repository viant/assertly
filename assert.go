package assertly

import (
	"fmt"
	"testing"
)

//AssertValues validates expected against actual data structure
func AssertValues(t *testing.T, expected, actual interface{}, arguments ...interface{}) bool {
	validation, err := Assert(expected, actual, NewDataPath("/"))
	if err != nil {
		if len(arguments) > 0 {
			handleFailure(t, fmt.Sprint(arguments...))
		}
		handleFailure(t, err)
		return false
	}

	if validation.FailedCount != 0 {
		if len(arguments) > 0 {
			handleFailure(t, arguments...)
		}
		for _, failure := range validation.Failures {
			handleFailure(t, fmt.Sprintf("%v: %v", failure.Path, failure.Message))
		}
		return false
	}
	return true
}
