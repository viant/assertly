package assertly_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"github.com/viant/toolbox"
	"testing"
	"fmt"
	"os"
)

type assertUseCase struct {
	Description string
	Expected    interface{}
	Actual      interface{}
	PassedCount int
	FailedCount int
	HasError    bool
}

func TestAssertMap(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "missing key test",
			Expected: map[string]interface{}{
				"k1": 1,
				"k2": 2.0,
				"k3": 11,
			},
			Actual: map[string]interface{}{
				"k1": 1,
				"k3": 11,
			},
			PassedCount: 2,
			FailedCount: 1,
		},
		{
			Description: "a map test",
			Expected: map[string]interface{}{
				"k1": 1,
				"k2": 2.0,
				"k3": 11,
			},
			Actual: map[string]interface{}{
				"k1": 1,
				"k2": 2.0,
				"k3": 11,
			},
			PassedCount: 3,
			FailedCount: 0,
		},
		{
			Description: "key does not exist violation test",
			Expected: map[string]interface{}{
				"k1": assertly.KeyDoesNotExistsDirective,
				"k2": 2.0,
			},
			Actual: map[string]interface{}{
				"k1": 1,
				"k2": 2.0,
			},
			PassedCount: 1,
			FailedCount: 1,
		},
		{
			Description: "key does not exist test",
			Expected: map[string]interface{}{
				"k1": assertly.KeyDoesNotExistsDirective,
				"k2": 2.0,
			},
			Actual: map[string]interface{}{
				"k2": 2.0,
			},
			PassedCount: 2,
			FailedCount: 0,
		},

		{
			Description: "key exists violation test",
			Expected: map[string]interface{}{
				"k1": assertly.KeyExistsDirective,
				"k2": 2.0,
			},
			Actual: map[string]interface{}{
				"k2": 2.0,
			},
			PassedCount: 1,
			FailedCount: 1,
		},
		{
			Description: "key exists test",
			Expected: map[string]interface{}{
				"k1": assertly.KeyExistsDirective,
				"k2": 2.0,
			},
			Actual: map[string]interface{}{
				"k1": 2.0,
				"k2": 2.0,
			},
			PassedCount: 2,
			FailedCount: 0,
		},

		{
			Description: "slice incompatible data type test",
			Expected: map[string]interface{}{
				"1": map[string]interface{}{
					"id":   1,
					"name": "name 1",
				},
				"2": map[string]interface{}{
					"id":   2,
					"name": "name 2",
				},
			},

			Actual: []interface{}{
				map[string]interface{}{
					"id":   1,
					"name": "name 1",
				},
				map[string]interface{}{
					"id":   2,
					"name": "name 2",
				},
			},
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "incompatible data type test",
			Expected: map[string]interface{}{
				assertly.IndexByDirective: "id",
				"1": map[string]interface{}{
					"id":   1,
					"name": "name 1",
				},
				"2": map[string]interface{}{
					"id":   2,
					"name": "name 2",
				},
			},
			Actual:      "123",
			PassedCount: 0,
			FailedCount: 1,
		},

		{
			Description: "slice index test",
			Expected: map[string]interface{}{
				assertly.IndexByDirective: "id",
				"1": map[string]interface{}{
					"id":   1,
					"name": "name 1",
				},
				"2": map[string]interface{}{
					"id":   2,
					"name": "name 2",
				},
			},
			Actual: []interface{}{
				map[string]interface{}{
					"id":   1,
					"name": "name 1",
				},
				map[string]interface{}{
					"id":   2,
					"name": "name 2",
				},
			},
			PassedCount: 4,
			FailedCount: 0,
		},

		{
			Description: "expected apply error",
			Expected: map[string]interface{}{
				assertly.CastDataTypeDirective + "k2": "abc",
				"k2":                                  2.0,
			},
			Actual:   map[string]interface{}{},
			HasError: true,
		},

		{
			Description: "time format directive test",
			Expected: map[string]interface{}{
				assertly.TimeFormatDirective + "k2": "yyyy-MM-dd",
				assertly.TimeFormatDirective + "k4": "yyyy-MM-dd",
				assertly.TimeFormatDirective:        "yyyy-MM-dd hh:mm:ss",

				"k2": "2019-01-01",
			},
			Actual: map[string]interface{}{
				"k2": "2019-01-01",
			},
			HasError:    false,
			PassedCount: 1,
		},
		{
			Description: "actual error",
			Expected: map[string]interface{}{
				assertly.TimeFormatDirective + "k2": "yyyy-MM-dd",
				"k2":                                "2019-01-01",
			},
			Actual: map[string]interface{}{
				"k2": "99-99-99",
			},
			HasError: true,
		},
	}
	runUseCases(t, useCases)
}

func TestAssertSlice(t *testing.T) {
	var useCases = []*assertUseCase{
		{
			Description: "slice test",
			Expected:    []int{1, 2, 3},
			Actual:      []int{1, 2, 3},
			PassedCount: 3,
		},
		{
			Description: "slice not equal test",
			Expected:    []int{1, 2, 3},
			Actual:      []int{1, 2, 4},
			PassedCount: 2,
			FailedCount: 1,
		},
		{
			Description: "slice len not equal test",
			Expected:    []int{1, 2, 3},
			Actual:      []int{1, 2},
			PassedCount: 2,
			FailedCount: 1,
		},
		{
			Description: "expected slice shorter - no violation since only supplied element are expected to be validated",
			Expected:    []int{1, 2},
			Actual:      []int{1, 2, 3},
			PassedCount: 2,
			FailedCount: 0,
		},
		{
			Description: "indexed slice test",
			Expected: []map[string]interface{}{
				{
					assertly.IndexByDirective: "key",
				},
				{
					"key": 1,
					"x":   100,
					"y":   200,
				},
				{
					"key": 2,
					"x":   200,
					"y":   300,
				},
			},
			Actual: []map[string]interface{}{
				{
					"key": 2,
					"x":   200,
					"y":   300,
				},
				{
					"key": 1,
					"x":   100,
					"y":   200,
				},
			},
			PassedCount: 6,
		},
		{
			Description: "incompatible slice type test",
			Expected:    []int{1, 2},
			Actual:      "1,2,3",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "different len slice test",
			Expected:    []int{},
			Actual:      []int{1},
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "zero len slice test",
			Expected:    []int{},
			Actual:      []int{},
			PassedCount: 1,
		},

		{
			Description: "indexed slice cast error test",
			Expected: []map[string]interface{}{
				{
					assertly.IndexByDirective:            "key",
					assertly.CastDataTypeDirective + "x": "float",
				},
				{
					"key": 1,
					"x":   100,
					"y":   200,
				},
				{
					"key": 2,
					"x":   200,
					"y":   300,
				},
			},
			Actual: []map[string]interface{}{
				{
					"key": 2,
					"x":   "abc",
					"y":   300,
				},
				{
					"key": 1,
					"x":   "100",
					"y":   200,
				},
			},
			PassedCount: 3,
			HasError:    true,
		},

		{
			Description: "indexed slice test",
			Expected: []map[string]interface{}{
				{
					"key":                                1,
					"x":                                  100,
					"y":                                  200,
					assertly.CastDataTypeDirective + "x": "float",
				},
			},
			Actual: []map[string]interface{}{
				{
					"key": 1,
					"x":   "xyz",
					"y":   200,
				},
			},
			PassedCount: 3,
			HasError:    true,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertJSONSlice(t *testing.T) {
	var useCases = []*assertUseCase{
		{
			Description: "JSON slice test",
			Expected: `[1,2,3]


[3,4]`,
			Actual: `[1,2,3]
[3,5]`,
			PassedCount: 4,
			FailedCount: 1,
		},
		{
			Description: "broken JSON slice test",
			Expected: `[1,2,3]

	[2,]
[3,4]`,
			Actual: `[1,2,3]
[3,5]`,
			PassedCount: 0,
			FailedCount: 1,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertText(t *testing.T) {
	var useCases = []*assertUseCase{
		{
			Description: "text qual test",
			Expected:    "123",
			Actual:      "123",
			PassedCount: 1,
		},
		{
			Description: "text qual test",
			Expected:    "123",
			Actual:      "1234",
			FailedCount: 1,
		},
		{
			Description: "text qual test",
			Expected:    "!123",
			Actual:      "1234",
			PassedCount: 1,
		},
		{
			Description: "text qual test",
			Expected:    "!123",
			Actual:      "123",
			FailedCount: 1,
		},
	}
	runUseCases(t, useCases)

}

type TestStructA struct {
	K1 string
	K2 int
}

func TestAssertStruct(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "struct test",
			Expected:    &TestStructA{K1: "123", K2: 123},
			Actual:      &TestStructA{K1: "123", K2: 123},
			PassedCount: 2,
			FailedCount: 0,
		},
		{
			Description: "struct with JSON test",
			Expected:    &TestStructA{K1: "123", K2: 123},
			Actual:      `{"K1":"123", "K2":124}`,
			PassedCount: 1,
			FailedCount: 1,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertSwitchCase(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "switch/case test",
			Expected: `{
		"@switchCaseBy@":"alg",
		"1": {
			"alg":1,
			"value":100
		},
		"2":{
			"alg":2,
			"value":200
		}
}`,
			Actual: `{
			"alg":2,
			"value":200
		}
`,
			PassedCount: 2,
			FailedCount: 0,
		},
		{
			Description: "missing switch/case test",
			Expected: `{
		"@switchCaseBy@":"alg",
		"1": {
			"alg":1,
			"value":100
		},
		"2":{
			"alg":2,
			"value":200
		}
}`,
			Actual: `{
			"alg":3,
			"value":200
		}
`,
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "missing switch/case setup error test",
			Expected: `{
		"@switchCaseBy@":"alg",
		"1": 1,
		"2": 2
}`,
			Actual: `{
			"alg":1,
			"value":200
		}
`,
			HasError: true,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertWithGlobalDirective(t *testing.T) {
	context := assertly.NewDefaultContext()
	directivePath := assertly.NewDataPath("")

	{
		directive1 := assertly.NewDirective(directivePath)
		directive1.AddDataType("id", "int")
		directive1.AddDataType("isEnabled", "bool")

		directive2 := assertly.NewDirective(directivePath.Key("k1"))
		directive2.AddTimeLayout("date", toolbox.DateFormatToLayout("yyyy-MM-dd"))

		testPath := assertly.NewDataPath("root")
		context.Directives = assertly.NewDirectives(directive1, directive2)
		{
			validation, err := assertly.AssertWithContext(`{
	"id":"213",
	"isEnabled":false,
	"done":"true"
}
`, `{
	"id":213,
	"isEnabled":"false",
	"done":"true1"
}
`, testPath.Key("field"), context)
			assert.Nil(t, err)
			assert.Equal(t, 2, validation.PassedCount)
			assert.Equal(t, 1, validation.FailedCount)

		}

		{
			validation, err := assertly.AssertWithContext(`{
	"date":"2017-01-01",
	"id":1
}
`, `{
	"date":"2017-01-01",
	"id":1
}
`, testPath.Key("k1"), context)
			assert.Nil(t, err)
			assert.Equal(t, 2, validation.PassedCount)
			assert.Equal(t, 0, validation.FailedCount)

		}

	}

}

func TestAssertRegExpr(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "reg expr test",
			Expected:    "~/.+(\\d+).+/",
			Actual:      "avc1erwer",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "reg expr test",
			Expected:    "~/.+(\\d+).+/",
			Actual:      "avcerwer",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "reg expr not test",
			Expected:    "!~/.+(\\d+).+/",
			Actual:      "avc1erwer",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "multiline reg expr not test",
			Expected:    "!~/.+(\\d+).+/",
			Actual:      "avc\ner\nwer",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "multiline reg expr test",
			Expected:    "~/^1.+3$/",
			Actual:      "1avc\n1ass3\nwer4",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "reg expr compilation error test",
			Expected:    "~/m???:1/",
			Actual:      "123",
			HasError:    true,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertMacro(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "macro-predicate test",
			Expected:    "<ds:between[1,10]>",
			Actual:      "3",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "macro-predicate violation test",
			Expected:    "<ds:between[1,10]>",
			Actual:      "13",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "macro-predicate error test",
			Expected:    "<ds:between[1,10, 2]>",
			Actual:      "13",
			HasError:    true,
		},
		{
			Description: "macro expansion",
			Expected:    `1<ds:env["USER"]>3`,
			Actual:      fmt.Sprintf("1%v3", os.Getenv("USER")),
			PassedCount: 1,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertRange(t *testing.T) {

	var useCases = []*assertUseCase{
		{
			Description: "range min max test",
			Expected:    "/[1..10]/",
			Actual:      "3",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "range min max test",
			Expected:    "/[1..10]/",
			Actual:      "30",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "not in min max range test",
			Expected:    "!/[1..10]/",
			Actual:      "30",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "range min max test",
			Expected:    "!/[1..10]/",
			Actual:      "3",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "range test",
			Expected:    "/[1,3,10]/",
			Actual:      "3",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "range test",
			Expected:    "/[1,3,10]/",
			Actual:      "4",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "range error test",
			Expected:    "/[3]/",
			Actual:      "4",
			HasError:    true,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertContains(t *testing.T) {
	var useCases = []*assertUseCase{
		{
			Description: "contain test",
			Expected:    "/123/",
			Actual:      "123456",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "contain violation test",
			Expected:    "/123/",
			Actual:      "3456",
			PassedCount: 0,
			FailedCount: 1,
		},
		{
			Description: "does not contain test",
			Expected:    "!/123/",
			Actual:      "30",
			PassedCount: 1,
			FailedCount: 0,
		},
		{
			Description: "range min max test",
			Expected:    "!/123/",
			Actual:      "01234",
			PassedCount: 0,
			FailedCount: 1,
		},
	}
	runUseCases(t, useCases)
}

func runUseCases(t *testing.T, useCases []*assertUseCase) {
	for _, useCase := range useCases {
		path := assertly.NewDataPath("/")
		validation, err := assertly.Assert(useCase.Expected, useCase.Actual, path)
		if err != nil {
			if useCase.HasError {
				continue
			}
			assert.Nil(t, err, useCase.Description)
			continue
		} else if useCase.HasError {
			assert.NotNil(t, err, useCase.Description)
			continue
		}
		assert.EqualValues(t, useCase.PassedCount, validation.PassedCount, useCase.Description)
		assert.EqualValues(t, useCase.FailedCount, validation.FailedCount, useCase.Description)
	}
}

func runUseCasesWithContext(t *testing.T, useCases []*assertUseCase, context *assertly.Context) {
	for _, useCase := range useCases {
		path := assertly.NewDataPath("/")
		validation, err := assertly.AssertWithContext(useCase.Expected, useCase.Actual, path, context)
		if err != nil {
			if useCase.HasError {
				continue
			}
			assert.Nil(t, err, useCase.Description)
			continue
		} else if useCase.HasError {
			assert.NotNil(t, err, useCase.Description)
			continue
		}
		assert.EqualValues(t, useCase.PassedCount, validation.PassedCount, useCase.Description)
		assert.EqualValues(t, useCase.FailedCount, validation.FailedCount, useCase.Description)
		assert.EqualValues(t, useCase.FailedCount > 0, validation.HasFailure(), )
	}
}

func TestAssertStructure(t *testing.T) {
	var useCases = []*assertUseCase{

		{
			Description: "data structure test",
			Expected: `{
  "1": {
    "id":1,
    "name":"name 1"
  },
  "2": {
    "id":2,
    "name":"name 2"
  }
}`,
			Actual: `{
  "1": {
    "id":1,
    "name":"name 1"
  },
  "2": {
    "id":2,
    "name":"name 22"
  }
}`,
			PassedCount: 3,
			FailedCount: 1,
		},
		{
			Description: "data structure test",
			Expected: `{
  "Meta": "abc",
  "Table": "abc",
  "Rows": [
    {
      "id": 1,
      "name": "name 1"
    },
    {
      "id": 2,
      "name": "name 2",
      "settings": {
        "k1": "v2"
      }
    },
    {
      "id": 2,
      "name": "name 2"
    }
  ]
}`,
			Actual: `{
"Table":"abc",
"Rows":[
{
	"id":1,
	"name":"name 12"
},
{
	"id":2,
	"name":"name 2",
	"settings": {
		"k1":"v20"
	}
},
{
	"id":4,
	"name":"name 2"
}
	]
}`,
			PassedCount: 5,
			FailedCount: 4,
		},
	}
	runUseCases(t, useCases)

}

func TestAssertStructureWithIndexDirective(t *testing.T) {
	var useCases = []*assertUseCase{
		{
			Description: "data structure with index directive",
			Expected: `{
  "1": {
    "id":1,
	"seq":0,
    "name":"name 1"
  },
  "2": {
    "id":2,
	"seq":0,
    "name":"name 2"
  }
}`,
			Actual: `{
  "1": {
    "id":1,
	"seq":0,
    "name":"name 1"
  },
  "2": {
    "id":2,
	"seq":0,
    "name":"name 22"
  }
}`,
			PassedCount: 5,
			FailedCount: 1,
		},
	}
	defaultDirective := assertly.NewDirective(assertly.NewDataPath(""))
	defaultDirective.IndexBy = []string{"id", "seq"}
	context := assertly.NewContext(nil, nil, assertly.NewDirectives(defaultDirective), nil)
	runUseCasesWithContext(t, useCases, context)

}
