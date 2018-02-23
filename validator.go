package assertly

import (
	"fmt"
	"github.com/viant/toolbox"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	MissingEntryViolation         = "entry was missing"
	IncompatibleDataTypeViolation = "data type was incompatible"
	KeyExistsViolation            = "key should exist"
	KeyDoesNotExistViolation      = "key should not exist"
	EqualViolation                = "value should be equal"
	NotEqualViolation             = "value should not be equal"
	LengthViolation               = "should have the same length"
	MissingCaseViolation          = "missing switch/case value"
	RegExprMatchesViolation       = "should match regrexp"
	RegExprDoesNotMatchViolation  = "should mot match regrexp"
	RangeViolation                = "should be in range"
	RangeNotViolation             = "should not be in range"
	ContainsViolation             = "should contain fragment"
	DoesNotContainViolation       = "should not contain fragment"
	PredicateViolation            = "should pass predicate"
)

//Assert validates expected against actual data structure for supplied path
func Assert(expected, actual interface{}, path DataPath) (*Validation, error) {
	context := NewDefaultContext()
	return AssertWithContext(expected, actual, path, context)
}

//AssertWithContext validates expected against actual data structure for supplied path and context
func AssertWithContext(expected, actual interface{}, path DataPath, context *Context) (*Validation, error) {
	validation := NewValidation()
	err := assertValue(expected, actual, path, context, validation)
	return validation, err
}

func getPredicate(input interface{}) toolbox.Predicate {
	predicate, ok := input.(toolbox.Predicate)
	if !ok {
		if predicatePointer, ok := input.(*toolbox.Predicate); ok {
			predicate = *predicatePointer
		}
	}
	return predicate
}

func expandExpectedText(text string, path DataPath, context *Context) (interface{}, error) {
	if context.Evaluator.HasMacro(text) {
		evaluated, err := context.Evaluator.Expand(context.Context, text)
		if err != nil {
			return nil, fmt.Errorf("failed to expand macro %v, path:%v, %v", text, path.Path(), err)
		}
		if !toolbox.IsString(evaluated) {
			return evaluated, nil
		}
		text = toolbox.AsString(evaluated)
	}
	if toolbox.IsCompleteJSON(text) {
		return asDataStructure(text), nil
	}
	return text, nil
}

func assertTime(expected *time.Time, actual interface{}, path DataPath, context *Context, validation *Validation) (err error) {
	dateLayout := path.Directive(context).DefaultTimeLayout()
	expectedTime, _ := toolbox.ToTime(expected, dateLayout)
	if actualTime, err := toolbox.ToTime(actual, dateLayout); err == nil {
		actual = actualTime
		if expectedTime.Equal(*actualTime) {
			validation.PassedCount++
			return nil
		}
	}
	validation.AddFailure(NewFailure(path.Source(), path.Path(), NotEqualViolation, expected, actual))
	return nil
}

func assertValue(expected, actual interface{}, path DataPath, context *Context, validation *Validation) (err error) {

	if expected == nil {
		if actual == nil {
			validation.PassedCount++
			return nil
		}
		validation.AddFailure(NewFailure(path.Source(), path.Path(), NotEqualViolation, expected, actual))
		return
	}

	if text, ok := expected.(string); ok {
		if expected, err = expandExpectedText(text, path, context); err != nil {
			return err
		}
	}
	predicate := getPredicate(expected)
	if predicate != nil {
		if !predicate.Apply(actual) {
			validation.AddFailure(NewFailure(path.Source(), path.Path(), PredicateViolation, predicate, actual))
		} else {
			validation.PassedCount++
		}
		return nil
	}
	if text, ok := actual.(string); ok {
		if toolbox.IsCompleteJSON(text) {
			actual = asDataStructure(text)
		}
	}

	dateLayout := path.Directive(context).DefaultTimeLayout()
	if toolbox.IsTime(expected) {
		expectedTime, _ := toolbox.ToTime(expected, dateLayout)
		return assertTime(expectedTime, actual, path, context, validation)
	} else if toolbox.IsStruct(expected) {
		var converter = toolbox.NewColumnConverter(dateLayout)
		var expectedMap = make(map[string]interface{})
		converter.AssignConverted(&expectedMap, expected)
		expected = expectedMap
		if toolbox.IsStruct(actual) {
			var actualMap = make(map[string]interface{})
			converter.AssignConverted(&actualMap, actual)
			actual = actualMap
		}
	}

	if toolbox.IsMap(expected) {
		return assertMap(toolbox.AsMap(expected), actual, path, context, validation)
	} else if toolbox.IsSlice(expected) {
		return assertSlice(toolbox.AsSlice(expected), actual, path, context, validation)
	} else if expected == actual || reflect.DeepEqual(expected, actual) {
		validation.PassedCount++
		return nil
	}

	directive := NewDirective(path)
	expectedText := toolbox.AsString(expected)

	if ! context.StrictDatTypeCheck {

		if expectedTime, err := toolbox.ToTime(expectedText, directive.DefaultTimeLayout()); err == nil && expectedTime != nil {
			if actualTime, err := toolbox.ToTime(actual, directive.DefaultTimeLayout());err == nil {
				if actualTime.Equal(*expectedTime) {
					validation.PassedCount++
					return nil
				}
			}
		}
	}
	return assertText(toolbox.AsString(expected), toolbox.AsString(actual), path, context, validation)
}

func isNegated(candidate string) (string, bool) {
	isNot := strings.HasPrefix(candidate, "!")
	if isNot {
		candidate = string(candidate[1:])
	}
	return candidate, isNot
}

func assertRegExpr(isNegated bool, expected, actual string, path DataPath, context *Context, validation *Validation) error {
	expected = string(expected[2: len(expected)-1])
	useMultiLine := strings.Count(actual, "\n") > 0
	pattern := ""
	if useMultiLine {
		pattern = "(?m)"
	}
	pattern += expected
	compiled, err := regexp.Compile(pattern)

	if err != nil {
		return fmt.Errorf("failed to compile %v, path: %v, %v", expected, path, err)
	}
	var matches = compiled.Match(([]byte)(actual))
	if !matches && !isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), RegExprMatchesViolation, expected, actual))
	} else if matches && isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), RegExprDoesNotMatchViolation, expected, actual))
	} else {
		validation.PassedCount++
	}
	return nil
}

func assertRange(isNegated bool, expected, actual string, path DataPath, context *Context, validation *Validation) error {
	if strings.Count(expected, "..")+strings.Count(expected, ",") == 0 {
		return fmt.Errorf("invalid range format, expected /[min..max]/ or /[val1,val2,valN]/, but had:%v, path: %v", expected, path.Path())
	}
	actual = strings.TrimSpace(actual)
	expected = string(expected[2: len(expected)-2])
	var rangeValues = strings.Split(expected, "..")
	var withinRange bool
	if len(rangeValues) > 1 {
		var minExpected = toolbox.AsFloat(strings.TrimSpace(rangeValues[0]))
		var maxExpected = toolbox.AsFloat(strings.TrimSpace(rangeValues[1]))
		var actualNumber = toolbox.AsFloat(actual)
		withinRange = actualNumber >= minExpected && actualNumber <= maxExpected
	} else {
		rangeValues = strings.Split(expected, ",")
		for _, candidate := range rangeValues {
			if strings.TrimSpace(candidate) == actual {
				withinRange = true
				break
			}
		}
	}
	if !withinRange && !isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), RangeViolation, expected, actual))
	} else if withinRange && isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), RangeNotViolation, expected, actual))
	} else {
		validation.PassedCount++
	}
	return nil
}

func assertContains(isNegated bool, expected, actual string, path DataPath, context *Context, validation *Validation) {
	expected = string(expected[1: len(expected)-1])
	contains := strings.Contains(actual, expected)

	if !contains && !isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), ContainsViolation, expected, actual))
	} else if contains && isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), DoesNotContainViolation, expected, actual))
	} else {
		validation.PassedCount++
	}
}

func assertText(expected, actual string, path DataPath, context *Context, validation *Validation) error {
	expected = strings.TrimSpace(expected)
	if strings.HasSuffix(expected, "/") {
		expected, isNegated := isNegated(expected)
		isRegExpr := strings.HasPrefix(expected, "~/")
		if isRegExpr {
			return assertRegExpr(isNegated, expected, actual, path, context, validation)
		}
		isRangeExpr := (strings.HasPrefix(expected, "/[") || strings.HasPrefix(expected, "!/[")) && strings.HasSuffix(expected, "]/")
		if isRangeExpr {
			return assertRange(isNegated, expected, actual, path, context, validation)
		}
		isContains := strings.HasPrefix(expected, "/")
		if isContains {
			assertContains(isNegated, expected, actual, path, context, validation)
			return nil
		}
	}
	expected, isNegated := isNegated(expected)
	isEqual := expected == actual
	if !isEqual && !isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), EqualViolation, expected, actual))
	} else if isEqual && isNegated {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), NotEqualViolation, expected, actual))
	} else {
		validation.PassedCount++
	}
	return nil
}

func actualMap(expected, actualValue interface{}, path DataPath, directive *Directive, validation *Validation) map[string]interface{} {
	var actual map[string]interface{} = nil
	if toolbox.IsMap(actualValue) {
		actual = toolbox.AsMap(actualValue)
	} else if toolbox.IsSlice(actualValue) {
		if len(directive.IndexBy) == 0 {
			validation.AddFailure(NewFailure(path.Source(), path.Path(), IncompatibleDataTypeViolation, expected, actualValue))
			return nil
		}
		aSlice := toolbox.AsSlice(actualValue)
		actual = indexSliceBy(aSlice, directive.IndexBy...)
	} else {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), IncompatibleDataTypeViolation, expected, actualValue))
		return nil
	}
	return actual
}



func assertMap(expected map[string]interface{}, actualValue interface{}, path DataPath, context *Context, validation *Validation) error {

	directive := NewDirective(path)
	directive.mergeFrom(path.Directive(context))
	directive.ExtractDirectives(expected)
	path.SetSource(directive.Source)
	var actual = actualMap(expected, actualValue, path, directive, validation)
	if actual == nil {
		return nil
	}
	directive.ExtractDataTypes(actual)
	if err := directive.Apply(actual); err != nil {
		return fmt.Errorf("failed to apply directive to actual, path:%v, %v", path.Path(), err)
	}

	if len(directive.SwitchBy) > 0 {
		switchValue := keysValue(actual, directive.SwitchBy...)
		caseValue, ok := expected[switchValue]
		if !ok {
			validation.AddFailure(NewFailure(path.Source(), path.Path(), MissingCaseViolation, expected, actual, directive.SwitchBy, switchValue))
			return nil
		}
		if !toolbox.IsMap(caseValue) {
			return fmt.Errorf("case value should be map but was %T, path: %v", caseValue, path.Path())
		}
		expected = toolbox.AsMap(caseValue)
	}

	if err := directive.Apply(expected); err != nil {
		return fmt.Errorf("failed to apply directive to expected, path:%v %v", path.Path(), err)
	}

	indexable := isIndexable(expected)
	if len(directive.IndexBy) == 0 {
		indexable = false
	}


	for expectedKey, expectedValue := range expected {
		if directive.IsDirectiveKey(expectedKey) {
			continue
		}
		var keyPath DataPath
		if indexable {
			keyPath = path.Key(keysPairValue(toolbox.AsMap(expectedValue), directive.IndexBy...))
		} else {
			keyPath = path.Key(expectedKey)
		}
		actualValue, ok := actual[expectedKey]
		if directive.KeyDoesNotExist[expectedKey] {
			if ok {
				validation.AddFailure(NewFailure(keyPath.Source(), keyPath.Path(), KeyDoesNotExistViolation, expectedKey, expectedKey))
			} else {
				validation.PassedCount++
			}
			continue
		}


		if directive.KeyExists[expectedKey] {
			if !ok {
				availableKeys := toolbox.MapKeysToStringSlice(expected)
				validation.AddFailure(NewFailure(keyPath.Source(), keyPath.Path(), KeyExistsViolation, expectedKey, strings.Join(availableKeys, ",")))
			} else {
				validation.PassedCount++
			}
			continue
		}

		if !ok {
			key := "key:" + expectedKey
			validation.AddFailure(NewFailure(keyPath.Source(), keyPath.Path(), MissingEntryViolation, expectedValue, toolbox.MapKeysToStringSlice(actual), key))
			continue
		}
		if err := assertValue(expectedValue, actualValue, keyPath, context, validation); err != nil {
			return err
		}
	}
	return nil
}

func assertSlice(expected []interface{}, actualValue interface{}, path DataPath, context *Context, validation *Validation) error {
	if actualValue == nil {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), IncompatibleDataTypeViolation, expected, actualValue))
		return nil
	}
	if !toolbox.IsSlice(actualValue) {
		validation.AddFailure(NewFailure(path.Source(), path.Path(), IncompatibleDataTypeViolation, expected, actualValue))
		return nil
	}
	var actual = toolbox.AsSlice(actualValue)
	if len(expected) == 0 {
		if len(expected) == len(actual) {
			validation.PassedCount++
			return nil
		}
		validation.AddFailure(NewFailure(path.Source(), path.Path(), LengthViolation, len(expected), len(actual)))
		return nil
	}
	directive := path.Directive(context)

	if toolbox.IsMap(expected[0]) {
		first := toolbox.AsMap(expected[0])

		if directive.ExtractDirectives(first) {
			expected = expected[1:]
		}
		for i := 0; i < len(actual); i++ {
			var actualMap = toolbox.AsMap(actual[i])
			directive.ExtractDataTypes(actualMap)
		}

		//add directive to expected
		for i := 0; i < len(expected); i++ {
			var expectedMap = toolbox.AsMap(expected[i])
			directive.Add(expectedMap)
			directive.Apply(expectedMap)
		}



		shouldIndex := len(directive.IndexBy) > 0
		if shouldIndex {
			expectedMap := indexSliceBy(expected, directive.IndexBy...)
			actualMap := indexSliceBy(actual, directive.IndexBy...)
			return assertMap(expectedMap, actualMap, path, context, validation)
		}

	}
	for i := 0; i < len(expected); i++ {
		if i >= len(actual) {
			validation.AddFailure(NewFailure(path.Source(), path.Path(), LengthViolation, len(expected), len(actual)))
			return nil
		}
		indexPath := path.Index(i)
		if err := assertValue(expected[i], actual[i], indexPath, context, validation); err != nil {
			return err
		}
	}
	return nil
}
