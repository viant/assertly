package assertly

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"reflect"
	"testing"
)

func TestAsDataStructure(t *testing.T) {

	var useCases = []struct {
		Description string
		Input       string
		Kind        reflect.Kind
		Length      int
	}{
		{
			Description: "multiline textual value",
			Input:       "abc\nxyz",
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "multiline textual with incomplete JSON like struct at the end",
			Input:       "abc\nxyz\n{123}",
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "multiline textual with incomplete JSON like struct at the begining and at the end",
			Input:       "{123}\nabc\nxyz\n{123}",
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "multiline textual with incomplete JSON like struct at the begining",
			Input:       "{123}\nabc\nxyz\n{123}",
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "valida JSON only at the first and last line",
			Input:       "[1,2,3]\nabc\nxyz\n[1,2,3]",
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "valida JSON only at the first and last line",
			Input:       "[1,2,3]\n[2,3,4]\n[1,2,3]",
			Kind:        reflect.Slice,
			Length:      3,
		},
		{
			Description: "multiline valid JSON object",
			Input: `{
		"k1":"v1",
		"k2": "v2"
}`,
			Kind:   reflect.Map,
			Length: 2,
		},
		{
			Description: "multiline valid JSON array",
			Input: `[1,
2,
3
]`,
			Kind:   reflect.Slice,
			Length: 3,
		},
		{
			Description: "new line delimited valid JSON array with empty line",
			Input: `[1,2,3]

[2,3,4]`,
			Kind:   reflect.Slice,
			Length: 2,
		},

		{
			Description: "multiline invalid JSON object",
			Input: `{
		"k1":"v1", z
		"k2", "v2"
}`,
			Kind:   reflect.String,
			Length: -1,
		},
		{
			Description: "empty string",
			Input:       ` `,
			Kind:        reflect.String,
			Length:      -1,
		},
		{
			Description: "single line valid JSON object",
			Input:       `{"a":1}`,
			Kind:        reflect.Map,
			Length:      1,
		},
	}

	for _, useCase := range useCases {
		output := asDataStructure(useCase.Input)
		actualKind := reflect.TypeOf(output).Kind()

		if assert.EqualValues(t, useCase.Kind, actualKind, useCase.Description) {
			if actualKind == reflect.String { //check for no string modification
				assert.EqualValues(t, useCase.Input, output, useCase.Description)
			}
			if useCase.Length >= 0 {
				var actualLength = 0
				if actualKind == reflect.Map {
					actualLength = len(toolbox.AsMap(output))
				} else {
					actualLength = len(toolbox.AsSlice(output))
				}
				assert.EqualValues(t, useCase.Length, actualLength, useCase.Description)
			}
		}
	}

}

func TestReverseSlice(t *testing.T) {
	var aSlice = []string{"1", "10", "3"}
	reverseSlice(aSlice)
	assert.EqualValues(t, []string{"3", "10", "1"}, aSlice)
}

func TestMergeTextMap(t *testing.T) {

	{
		var source = map[string]string{}
		var target map[string]string
		mergeTextMap(source, &target)
		assert.EqualValues(t, 0, len(target))
	}
	{
		var source = map[string]string{
			"k1": "v1",
		}
		var target map[string]string
		mergeTextMap(source, &target)
		assert.EqualValues(t, 1, len(target))
	}
	{
		var source = map[string]string{
			"k1": "v1",
		}
		var target = make(map[string]string)
		mergeTextMap(source, &target)
		assert.EqualValues(t, 1, len(target))
	}

}

func TestMergeBoolMap(t *testing.T) {

	{
		var source = map[string]bool{}
		var target map[string]bool
		mergeBoolMap(source, &target)
		assert.EqualValues(t, 0, len(target))
	}
	{
		var source = map[string]bool{
			"k1": true,
		}
		var target map[string]bool
		mergeBoolMap(source, &target)
		assert.EqualValues(t, 1, len(target))
	}
	{
		var source = map[string]bool{
			"k1": true,
		}
		var target map[string]bool
		mergeBoolMap(source, &target)
		assert.EqualValues(t, 1, len(target))
	}

}

func Test_ToStringSlice(t *testing.T) {
	{
		aSlice := toStringSlice([]interface{}{"1", 2})
		assert.EqualValues(t, []string{"1", "2"}, aSlice)
	}
	{
		aSlice := toStringSlice(1)
		assert.EqualValues(t, []string{"1"}, aSlice)
	}
}
