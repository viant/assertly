package assertly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailure_Index(t *testing.T) {

	{
		var failure = NewFailure("[/]:ad[12].we", "", nil, nil)
		assert.EqualValues(t, 12, failure.Index())
	}
	{
		var failure = NewFailure("ad[12].we", "", nil, nil)
		assert.EqualValues(t, -1, failure.Index())
	}
	{
		var failure = NewFailure(":ad[a].we", "", nil, nil)
		assert.EqualValues(t, -1, failure.Index())
	}
	{
		var failure = NewFailure(":ad[].we", "", nil, nil)
		assert.EqualValues(t, -1, failure.Index())
	}
}

func TestFailure_MergeFrom(t *testing.T) {
	var failure = NewFailure("[/]:ad[12].we", "", nil, nil)
	source := NewValidation()
	source.PassedCount = 2
	source.AddFailure(failure)
	target := NewValidation()
	target.AddFailure(failure)
	target.PassedCount = 2
	target.MergeFrom(source)
	assert.EqualValues(t, 4, target.PassedCount)
	assert.EqualValues(t, 2, target.FailedCount)
	assert.EqualValues(t, 2, len(target.Failures))
}

func TestValidation_Report(t *testing.T) {

	source := NewValidation()
	source.AddFailure(NewFailure(":ad[].we", "test", nil, nil))
	source.PassedCount++
	assert.EqualValues(t, "test\nPassed: 1\nFailed: 1", source.Report())

}
