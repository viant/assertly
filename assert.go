package assertly

import (
	"fmt"
	"testing"
)

//AssertValues validates expected against actual data structure
func AssertValues(t *testing.T, expected, actual interface{}, arguments ...interface{}) bool {
	validation, err := Assert(expected, actual, NewDataPath("/"))
	if err != nil {
		return handlerValidationError(t, err, arguments...)
	}
	if validation.FailedCount != 0 {
		return handlerValidationFailures(t, validation, arguments...)
	}
	return true
}

//AssertValuesWithContext validates expected against actual data structure with context
func AssertValuesWithContext(context *Context, t *testing.T, expected, actual interface{}, arguments ...interface{}) bool {
	validation, err := AssertWithContext(expected, actual, NewDataPath("/"), context)
	if err != nil {
		return handlerValidationError(t, err, arguments...)
	}
	if validation.FailedCount != 0 {
		return handlerValidationFailures(t, validation, arguments...)
	}
	return true
}

func handlerValidationError(t *testing.T, err error, arguments ...interface{}) bool {
	if len(arguments) > 0 {
		handleFailure(t, fmt.Sprint(arguments...))
	}
	handleFailure(t, err)
	return false
}

func handlerValidationFailures(t *testing.T, validation *Validation, arguments ...interface{}) bool {
	if len(arguments) > 0 {
		handleFailure(t, arguments...)
	}
	for _, failure := range validation.Failures {
		handleFailure(t, fmt.Sprintf("%v: %v", failure.Path, failure.Message))
	}
	return false
}
