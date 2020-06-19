package assertly

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"testing"
	"time"
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

	{ //test key exists, does not exist directive
		var source = map[string]interface{}{
			KeyDoesNotExistsDirective + "k001": true,
			KeyExistsDirective + "k002":        true,
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
			"d1": date,
			"d2": dateHour,
			"d3": dateHourMiniute,
			"d4": dateHourMiniuteSec,
			"f":  3.2,
			"i":  213,
			"b":  true,
		}

		directive := NewDirectives()
		directive.ExtractDataTypes(source)
		assert.EqualValues(t, "float", directive.DataType["f"])
		assert.EqualValues(t, "int", directive.DataType["i"])
		assert.EqualValues(t, "bool", directive.DataType["b"])
		assert.EqualValues(t, "2006-01-02 15:04:05.000-07", directive.TimeLayouts["d4"])
	}

}

func TestDirective_Add(t *testing.T) {

	directive := NewDirective(NewDataPath("/"))
	directive.DataType = map[string]string{"f1": "float"}
	directive.IndexBy = []string{"id"}
	directive.SwitchBy = []string{"case"}
	directive.DataType = map[string]string{"f1": "float"}
	directive.TimeLayouts = map[string]string{"t1": "2016-01"}
	directive.TimeLayout = "2016-01-02"

	var target = make(map[string]interface{})
	directive.Add(target)

	assert.EqualValues(t, "float", target[CastDataTypeDirective+"f1"])
	assert.EqualValues(t, "2016-01-02", target[TimeLayoutDirective])
	assert.EqualValues(t, "2016-01", target[TimeLayoutDirective+"t1"])
	assert.EqualValues(t, []string{"id"}, target[IndexByDirective])
	assert.EqualValues(t, []string{"case"}, target[SwitchByDirective])

}

func TestAssertPath_Directive(t *testing.T) {

	expected := `[
	{
		"@indexBy@":"SubPath"
	},
	{
		"SubPath": "group1.field1",
		"Expected": 1
	},
	{
		"SubPath": "group1.field2",
		"Expected": 2
	},
	{
		"SubPath": "group1.field3",
		"Expected": 3
	}
]
`

	{
		directive := NewDirective(NewDataPath("/"))
		var target = map[string]interface{}{
			AssertPathDirective + "group1.field1": 1,
			AssertPathDirective: []interface{}{
				map[string]interface{}{
					"group1.field2": 2,
				},
				map[string]interface{}{
					"group1.field3": 3,
				},
			},
		}
		directive.ExtractDirectives(target)
		AssertValues(t, expected, directive.AssertPaths)
	}

	{
		directive := NewDirective(NewDataPath("/"))
		var target = map[string]interface{}{
			AssertPathDirective + "group1.field1": 1,
			AssertPathDirective: map[string]interface{}{
				"group1.field2": 2,
				"group1.field3": 3,
			},
		}
		directive.ExtractDirectives(target)
		AssertValues(t, expected, directive.AssertPaths)
	}

}
