package assertly

import (
	"fmt"
	"strings"
)


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
	for _, failure := range source.Failures {
		v.AddFailure(failure)
	}
}

//Report returns validation report
func (v *Validation) Report() string {
	var result = make([]string, 0)
	for _, failure := range v.Failures {
		result = append(result,  failure.Path + ": " +  failure.Message)
	}
	result = append(result, fmt.Sprintf("Passed: %v", v.PassedCount))
	result = append(result, fmt.Sprintf("Failed: %v", v.FailedCount))
	return strings.Join(result, "\n")
}

//NewValidation returns new validation
func NewValidation() *Validation {
	return &Validation{
		Failures: make([]*Failure, 0),
	}
}
