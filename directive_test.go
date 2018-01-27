package assertly

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirective_Extract(t *testing.T) {

	directive := NewDirectives()

	{ //test switch by directive
		var source = map[string]interface{}{
			"k1":              1,
			SwitchByDirective: "k1",
		}
		directive.Extract(source)
		assert.Equal(t, []string{"k1"}, directive.SwitchBy)

	}

	{ //test index by directive
		var source = map[string]interface{}{
			"k1":             1,
			"k2":             3,
			IndexByDirective: []string{"k1", "k2"},
		}
		directive.Extract(source)
		assert.Equal(t, []string{"k1", "k2"}, directive.IndexBy)

	}
	{ //test index by directive
		var source = map[string]interface{}{
			"k1":             1,
			"k2":             3,
			IndexByDirective: []string{"k1", "k2"},
		}
		directive.Extract(source)
		assert.Equal(t, []string{"k1", "k2"}, directive.IndexBy)

	}
	{ //test key exists, does not exist directive
		var source = map[string]interface{}{
			"k1": KeyDoesNotExistsDirective,
			"k2": KeyExistsDirective,
		}
		directive.Extract(source)
		assert.True(t, directive.KeyExists["k2"])
		assert.True(t, directive.KeyDoesNotExist["k1"])

	}
	{ //test key time format
		var source = map[string]interface{}{
			TimeFormatDirective + "k2": "yyyy-MM-dd",
			TimeFormatDirective:        "yyyy-MM-dd",

			"k2": "2011-02-11",
		}
		directive.Extract(source)
		assert.EqualValues(t, "2006-01-02", directive.TimeLayouts["k2"])
		assert.EqualValues(t, "2006-01-02", directive.DefaultTimeLayout())

	}
	{ //test key time format
		var source = map[string]interface{}{
			CastDataTypeDirective + "k2": "float",
			"k2": "3.7",
		}
		directive.Extract(source)
		assert.EqualValues(t, "float", directive.DataType["k2"])
	}

}
