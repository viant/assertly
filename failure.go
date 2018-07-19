package assertly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
	"unicode"
)

//Failure represents a validation failre
type Failure struct {
	Source   string
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

func (f *Failure) LeafKey() string {
	leafIndex := strings.LastIndex(f.Path, ".")
	if leafIndex == -1 {
		return ""
	}
	return string(f.Path[leafIndex+1:])
}

func removeDirectives(aMap map[string]interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	for k, v := range aMap {
		if strings.HasPrefix(k, "@") {
			continue
		}
		result[k] = v
	}
	return result
}

//NewFailure creates a new failure
func NewFailure(source, path string, reason string, expected, actual interface{}, args ...interface{}) *Failure {

	if expected != nil && toolbox.IsMap(expected) {
		expected = removeDirectives(toolbox.AsMap(expected))
	}
	if actual != nil && toolbox.IsMap(actual) {
		actual = removeDirectives(toolbox.AsMap(actual))
	}

	var result = &Failure{
		Source:   source,
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
		return fmt.Sprintf("entry for %v was missing, expected: %v, actual keys: %v", failure.Args[0], failure.Expected, failure.Actual)
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
		return fmt.Sprintf("actual length %v  was not equal: %v", failure.Actual, failure.Expected)
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
