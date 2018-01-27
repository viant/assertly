package assertly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataPath_Key(t *testing.T) {
	var path = NewDataPath("root")
	subPath := path.Key("f1").Index(1).Key("s1")
	assert.Equal(t, "f1/*/s1", subPath.MatchingPath())
}

func TestDataPath_Path(t *testing.T) {
	var path = NewDataPath("root")
	subPath := path.Key("f1").Index(1).Key("s1")
	assert.Equal(t, "root:f1[1].s1", subPath.Path())
}
