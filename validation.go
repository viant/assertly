package assertly

//Failure represents a validation failre
type Failure struct {
	Path     string
	Expected interface{}
	Actual   interface{}
	Reason   string
}

func NewFailure(path string, reason string, expected, actual interface{}) *Failure {
	return &Failure{
		Path:     path,
		Reason:   reason,
		Expected: expected,
		Actual:   actual,
	}
}

//Validation validation
type Validation struct {
	PassedCount int
	FailedCount int
	Failures    []*Failure
}

func (v *Validation) AddFailure(failure *Failure) {
	if len(v.Failures) == 0 {
		v.Failures = make([]*Failure, 0)
	}
	v.Failures = append(v.Failures, failure)
	v.FailedCount++
}

//NewValidation returns new validation
func NewValidation() *Validation {
	return &Validation{
		Failures: make([]*Failure, 0),
	}
}
