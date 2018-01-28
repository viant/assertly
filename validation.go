package assertly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
	"unicode"
)

//Failure represents a validation failre
type Failure struct {
	Path     string
	Expected interface{}
	Actual   interface{}
	Args     []interface{}
	Reason   string
	Message  string
}

func (f *Failure) Index() int {
	pair := strings.SplitN(f.Path, ":", 2)
	if len(pair) != 2 {
		return -1
	}
	var index = ""
	var expectIndex = false
outer:
	for _, r := range pair[1] {

		switch r {
		case '[':
			expectIndex = true
		case ']':
			expectIndex = false
			if len(index) > 0 {
				break outer
			}
		default:
			if expectIndex && unicode.IsNumber(r) {
				index += string(r)
			}
		}
	}
	if len(index) == 0 {
		return -1
	}
	return toolbox.AsInt(index)
}



//NewFailure creates a new failure
func NewFailure(path string, reason string, expected, actual interface{}, args ... interface{}) *Failure {
	var result = &Failure{
		Path:     path,
		Reason:   reason,
		Expected: expected,
		Actual:   actual,
		Args:     args,
	}
	result.Message = FormatMessage(result)
	return result
}

func FormatMessage(failure *Failure) string {
	switch failure.Reason {
	case MissingEntryViolation:
		fmt.Sprintf("%v was missing, expected: %v", failure.Args[0], failure.Expected)
	case IncompatibleDataTypeViolation:
		return fmt.Sprintf("actual was %T, but expected %T(%v)", failure.Actual, failure.Expected, failure.Expected)
	case KeyExistsViolation:
		fmt.Sprintf("key '%v' should exists", failure.Expected)
	case KeyDoesNotExistViolation:
		fmt.Sprintf("'%v' should not exists", failure.Expected)
	case EqualViolation:
		return fmt.Sprintf("actual(%T): '%v' was not equal (%T) '%v'", failure.Actual, failure.Actual, failure.Expected, failure.Expected)
	case NotEqualViolation:
		return fmt.Sprintf("actual(%T): '%v' was equal (%T) '%v'", failure.Actual, failure.Actual, failure.Expected, failure.Expected)
	case LengthViolation:
		return fmt.Sprintf("actual length '%v'  was not equal: '%v'", failure.Actual, failure.Actual)
	case MissingCaseViolation:
		switchBy := failure.Args[0].([]string)
		caseValue := toolbox.AsString(failure.Args[1])
		var availableKeys = toolbox.MapKeysToStringSlice(failure.Expected)
		return fmt.Sprintf("actual case %v => %v, was missing in expected set: available keys: [%v]",
			strings.Join(switchBy, ","), caseValue, strings.Join(availableKeys, ","))
	case RegExprMatchesViolation:
		return fmt.Sprintf("actual: '%v' should matched %v", failure.Actual, failure.Expected)
	case RegExprDoesNotMatchViolation:
		return fmt.Sprintf("actual: '%v' should not be matched %v", failure.Actual, failure.Expected)
	case RangeViolation:
		return fmt.Sprintf("actual '%v' is not in: '%v'", failure.Actual, failure.Expected)
	case RangeNotViolation:
		return fmt.Sprintf("actual '%v' should not be in: '%v'", failure.Actual, failure.Expected)
	case ContainsViolation:
		return fmt.Sprintf("actual '%v' does not contain: '%v'", failure.Actual, failure.Expected)
	case DoesNotContainViolation:
		return fmt.Sprintf("actual '%v' should not not contain: '%v'", failure.Actual, failure.Expected)
	case PredicateViolation:
		return fmt.Sprintf("actual '%v' should pass predicate: '%v'", failure.Actual, failure.Expected)
	}
	return failure.Reason
}

//Validation validation
type Validation struct {
	TagID       string
	Description string
	PassedCount int
	FailedCount int
	Failures    []*Failure
}

//AddFailure add failure to current violation
func (v *Validation) AddFailure(failure *Failure) {
	if len(v.Failures) == 0 {
		v.Failures = make([]*Failure, 0)
	}
	v.Failures = append(v.Failures, failure)
	v.FailedCount++
}

//HasFailure returns true if validation has failures
func (v *Validation) HasFailure() bool {
	return v.FailedCount > 0
}

//MergeFrom merges failures and passes from source
func (v *Validation) MergeFrom(source *Validation) {
	v.PassedCount += source.PassedCount
	for _,  failure:= range source.Failures {
		v.AddFailure(failure)
	}
}

//NewValidation returns new validation
func NewValidation() *Validation {
	return &Validation{
		Failures: make([]*Failure, 0),
	}
}
