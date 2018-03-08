package assertly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)



func TestValidation_Report(t *testing.T) {

	source := NewValidation()
	source.AddFailure(NewFailure("",":ad[].we", "test", nil, nil))
	source.PassedCount++
	assert.EqualValues(t, ":ad[].we: test\nPassed: 1\nFailed: 1", source.Report())

}
