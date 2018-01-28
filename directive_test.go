package assertly

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/viant/toolbox"
)

func TestDirective_ExtractExpected(t *testing.T) {

	directive := NewDirectives()

	{ //test switch by directive
		var source = map[string]interface{}{
			"k1":              1,
			SwitchByDirective: "k1",
		}
		directive.ExtractDirectives(source)
		assert.Equal(t, []string{"k1"}, directive.SwitchBy)

	}

	{ //test index by directive
		var source = map[string]interface{}{
			"k1":             1,
			"k2":             3,
			IndexByDirective: []string{"k1", "k2"},
		}
		directive.ExtractDirectives(source)
		assert.Equal(t, []string{"k1", "k2"}, directive.IndexBy)

	}
	{ //test index by directive
		var source = map[string]interface{}{
			"k1":             1,
			"k2":             3,
			IndexByDirective: []string{"k1", "k2"},
		}
		directive.ExtractDirectives(source)
		assert.Equal(t, []string{"k1", "k2"}, directive.IndexBy)

	}
	{ //test key exists, does not exist directive
		var source = map[string]interface{}{
			"k1": KeyDoesNotExistsDirective,
			"k2": KeyExistsDirective,
		}
		directive.ExtractDirectives(source)
		assert.True(t, directive.KeyExists["k2"])
		assert.True(t, directive.KeyDoesNotExist["k1"])

	}
	{ //test key time format
		var source = map[string]interface{}{
			TimeFormatDirective + "k2": "yyyy-MM-dd",
			TimeFormatDirective:        "yyyy-MM-dd",

			"k2": "2011-02-11",
		}
		directive.ExtractDirectives(source)
		assert.EqualValues(t, "2006-01-02", directive.TimeLayouts["k2"])
		assert.EqualValues(t, "2006-01-02", directive.DefaultTimeLayout())

	}
	{ //test key time format
		var source = map[string]interface{}{
			CastDataTypeDirective + "k2": "float",
			"k2":                         "3.7",
		}
		directive.ExtractDirectives(source)
		assert.EqualValues(t, "float", directive.DataType["k2"])
	}

}

func truncateDate(dateFormat string) *time.Time {
	var layout = toolbox.DateFormatToLayout(dateFormat)
	date := time.Now().Format(layout)
	result, _ := time.Parse(layout, date)
	return &result
}

func TestDirective_ExtractDataTypes(t *testing.T) {
	date := truncateDate("yyyy-MM-dd")
	dateHour := truncateDate("yyyy-MM-dd hh")
	dateHourMiniute := truncateDate("yyyy-MM-dd hh:mm")
	dateHourMiniuteSec := truncateDate("yyyy-MM-dd hh:mm:ss")

	{ //test index by directive
		var source = map[string]interface{}{
			"d1":             date,
			"d2":             dateHour,
			"d3":             dateHourMiniute,
			"d4":             dateHourMiniuteSec,
			"f":              3.2,
			"i":              213,
			"b":              true,
		}

		directive := NewDirectives()
		directive.ExtractDataTypes(source)
		assert.EqualValues(t, "float", directive.DataType["f"])
		assert.EqualValues(t, "int", directive.DataType["i"])
		assert.EqualValues(t, "bool", directive.DataType["b"])

		assert.EqualValues(t, "2006-01-02", directive.TimeLayouts["d1"])
		assert.EqualValues(t, "2006-01-02 03", directive.TimeLayouts["d2"])
		assert.EqualValues(t, "2006-01-02 03:04", directive.TimeLayouts["d3"])
		assert.EqualValues(t, "2006-01-02 03:04:05", directive.TimeLayouts["d4"])

	}

}
