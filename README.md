# Data structure testing library (assertly)

[![Data structure testing library for Go.](https://goreportcard.com/badge/github.com/viant/assertly)](https://goreportcard.com/report/github.com/viant/assertly)
[![GoDoc](https://godoc.org/github.com/viant/assertly?status.svg)](https://godoc.org/github.com/viant/assertly)

This library is compatible with Go 1.8+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Introduction](#Introduction)
- [Motivation](#Motivation)
- [Usage](#Usage)
- [Validation](#Validation)
- [Directive](#Directive)
- [Macros](#Macros)
- [External Resource](#external)
- [License](#License)
- [Credits and Acknowledgements](#Credits-and-Acknowledgements)



<a name="Introduction"></a>
## Introduction

This library enables complex data structure testing, specifically: 
1. Realtime transformation or casting of incompatible data types with directives system.
2. Consistent way of testing of unordered structures. 
3. Contains, Range, RegExp support on any data structure deeph level.
4. Switch case directive to provide expected value alternatives based on actual switch/case input match.
5. Macro system enabling complex predicate and expression evaluation, and customization.


<a name="Motivation"></a>

## Motivation

This library has been created as a way to unify original testing approaches introduced 
to [dsunit](https://github.com/viant/dsunit) and [endly](https://github.com/viant/endly)






<a name="Usage"></a>
## Usage

**Complete data validation with concrete types**

```go

import(
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
)


func Test_XX(t *testing.T) {
    
   	
   	
   	var actualRecords []*User = //get actual
   	var expectedRecords []*User = //get expected
   	assertly.AssertValues(t, expectedRecords, actualRecords)
   	
   	//or with custom path and testing.T integration
   	validation, err := assertly.Assert(expected, actual, assertly.NewDataPath("/"))
   	assert.EqualValues(t, 0, validation.FailedCount, validation.Report())

   	
}

```
**Partial data validation with directive and reg expression**

```go

func Test_XX(t *testing.T) {
    
    var actualConfig = &Config{
        Endpoint: &Endpoint{
            Port: 8080,
            TimeoutMs: 2000,
        },
        LogTypes: map[string]*LogType{
            "type1": &LogType{
                Locations:[]*Location{
                    {
                        URL:"file:///data/log/type1",
                    },
                },
                MaxQueueSize: 2048,
                QueueFlashCount: 1024,
                FlushFrequencyInMs: 500,
            },
            "type2":  &LogType{
                Locations:[]*Location{
                    {
                        URL:"file:///data/log/type2",
                    },
                },
                MaxQueueSize: 4096,
                QueueFlashCount: 2048,
                FlushFrequencyInMs: 1000,
            },
        },
    }
                           
       
    var expectedConfig = expected: `{
      "Endpoint": {
        "Port": 8080,
        "TimeoutMs": 2000
      },
      "LogTypes": {
        "type1": {
          "Locations":[
            {
              "URL":"~/type1/"
            }
          ],
          "MaxQueueSize": 2048,
          "QueueFlashCount": 1024,
          "FlushFrequencyInMs": 500
        },
        "type2": "@exists@" 
      }
    }`,
    
    assertly.AssertValues(t, expectedConfig, actualConfig)
   	
}

```

- reg expression:  "URL":"~/type1/"
- directive: @exists@




**Validation with custom macro value provider**

```go

type fooProvider struct{}

func (*fooProvider) Get(context toolbox.Context, arguments ...interface{}) (interface{}, error) {
	var args = []string{}
	for _, arg := range arguments {
		args = append(args, toolbox.AsString(arg))
	}
	return fmt.Sprintf("foo{%v}", strings.Join(args, ",")), nil
}

func Test_XX(t *testing.T) {
	ctx := NewDefaultContext()
	var provider toolbox.ValueProvider = &fooProvider{}
	ctx.Evaluator.ValueProviderRegistry.Register("foo", provider)

	var actual = map[string]string{
		"k1":"v1",
		"k2":"Macro test: foo{1,abc} !",
	}
	
	var expected = map[string]string{
		"k1":"v1",
		"k2":"Macro test: <ds:foo[1,\"abc\"]> !",
	}


	AssertValuesWithContext(ctx, t, expected, actual)
}

```


**Validation with custom predicate**
```go
type rangePredicate struct {
	min int
	max int
	actual int
	err error
}

func (p *rangePredicate) String() string {
	return fmt.Sprintf("min: %v, max: %v, actual: %v, err: %v", p.min, p.max, p.actual, p.err)
}

func (p *rangePredicate) Apply(value interface{}) bool {
	p.actual, p.err = toolbox.ToInt(value)
	return p.actual >= p.min && p.actual <= p.max
}



type inRangePredicateProvider struct{}
func (*inRangePredicateProvider) Get(context toolbox.Context, arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 2 {
		return nil, fmt.Errorf("expected 2 arguments (min, max) but had: %v", len(arguments))
	}
	min, err := toolbox.ToInt(arguments[0])
	if err != nil {
		return nil, fmt.Errorf("invalid min %v", err)
	}
	max, err := toolbox.ToInt(arguments[1])
	if err != nil {
		return nil, fmt.Errorf("invalid min %v", err)
	}
	var predicate toolbox.Predicate =  &rangePredicate{min:min, max: max}
	return &predicate, nil
}



func Test_XX(t *testing.T) {
	ctx := NewDefaultContext()
	var provider toolbox.ValueProvider = &inRangePredicateProvider{}
	ctx.Evaluator.ValueProviderRegistry.Register("inRange", provider)


	var actual = map[string]int{
		"k1":1,
		"k2":3,
	}


	var expected = map[string]string{
		"k1":"1",
		"k2":"<ds:inRange[2,10]>",
	}


	AssertValuesWithContext(ctx, t, expected, actual)
}

```



<a name="Validation"></a>
## Validation


Validation rules:
1) JSON textual data is converted into data structure
2) New Line Delimited JSON is converted into data structure collection.
3) Object/Struct is converted into data structure
4) Only existing keys/fields in expected data structure are validated  
5) Only existing items in the array/slice are validated
6) Directive and macros/predicate provide validation extension
7) The following expression can be used on any data structure level:

| Assertion Type |  input | expected expression | example | 
| --- | --- | --- | --- | 
| equal |  actual | expected | a:a |
| not equal |  actual | !expected | a:!b |
| contains | actual | /expected/| abcd:/bc/|
| not contains | actual | !/expected/| abcd:!/xc/ |
| regExpr | actual | ~/expected/ | 1234a:/\d+/ |
| not regExpr | actual | !~/expected/ | 1234:!/\w/ |
| between | actual | /[minExpected..maxExpected]/ | 12:/[1..13]/ |
| exists | n/a | { "key": "@exists@" } | |
| not exists | n/a | { "key": "@!exists@" } | |

**example**:

```go

func Test_XX(t *testing.T) {
    
var expected = `
{
  "Meta": "abc",
  "Table": "/table_/",
  "Rows": [
    {
      "id": 1,
      "name": "~/name (\\d+)/",
      "@exists@":"dob"
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
var actual = `
{
  "Table": "table_xx",
  "Rows": [
    {
      "id": 1,
      "name": "name 12",
      "dob":"2018-01-01"
    },
    {
      "id": 2,
      "name": "name 2",
      "settings": {
        "k1": "v20"
      }
    },
    {
      "id": 4,
      "name": "name 2"
    }
  ]
}`,
	
    validation, err := assertly.Assert(expected, actual, assertly.NewDataPath("/"))
   	assert.EqualValues(t, 0, validation.FailedCount, validation.Report())
}


```


<a name="Directive"></a>
## Directive

Directive is piece of information instructing validator to either convert data just before validation takes place or to validate a date according to provided rules. 
-	KeyExistsDirective          = "@exists@"
-	KeyDoesNotExistsDirective   = "@!exists@"
-	TimeFormatDirective         = "@timeFormat@"
-	TimeLayoutDirective         = "@timeLayout@"
-	SwitchByDirective           = "@switchCaseBy@"
-	CastDataTypeDirective       = "@cast@"
-	IndexByDirective            = "@indexBy@"
-	CaseSensitiveDirective      = "@caseSensitive@"
-	KeyCaseSensitiveDirective = "@CaseSensitive@"
-   NumericPrecisionPoint       = "@numericPrecisionPoint@"
-   CoalesceWithZero            = "@coalesceWithZero@"

### Index by

**@indexBy@** - index by directive indexes a slice for validation process, specifically.

1) Two unordered array/slice/collection that can be index by a unique fields 
2) A map with a actual array/slice/collection that can be ordered by unique fields




**Example 1**


\#expected
```json
{
"@indexBy@":"id",
"1" :{"id":1, "name":"name1"},
"2" :{"id":2, "name":"name2"}
}

```
	
\#actual
```json
[
{"id":1, "name":"name1"},
{"id":2, "name":"name2"}
]

```

**Example 2**

\#expected
```json
{"@indexBy@":"id"}
{"id":1, "name":"name1"}
{"id":2, "name":"name2"}
```
	
\#actual
```json
{"id":1, "name":"name1"}
{"id":2, "name":"name2"}
```




**Example 3**

\#expected
```json
{"@indexBy@":"request.id"}
{"request":{"id":1111, "name":"name1"}, "ts":189321233}
{"request":{"id":2222, "name":"name2"}, "ts":189321235}
```
	
\#actual
```json
{"request":{"id":2222, "name":"name2"}, "ts":189321235}
{"request":{"id":1111, "name":"name1"}, "ts":189321233}
```


## Switch/case 

**@switchCaseBy@** - switch directive instructs a validator to select matching expected subset based on some actual value.
.
For non deterministic system there could be various alternative output for the same input.

**Example 1**

\#expected 
 ```json
 [
   {
     "@switchCaseBy@":["experimentID"]
   },
   {
     "1":{"experimentID":1, "seq":1, "outcome":[1.53,7.42,6.34]},
     "2":{"experimentID":2, "seq":1, "outcome":[3.53,6.32,3.34]}
   },
   {
     "1":{"experimentID":1, "seq":2, "outcome":[5.63,4.3]},
     "2":{"experimentID":1, "seq":2, "outcome":[3.65,3.2]}
   }
 ]
```

\#actual
```json
{"experimentID":1, "seq":1, "outcome":[1.53,7.42,6.34]}
{"experimentID":1, "seq":2, "outcome":[5.63,4.3]}
``` 


**Example 2**

\#expected 
 ```json
 [
   {
     "@switchCaseBy@":["experimentID"]
   },
   {
     "1":{"experimentID":1, "seq":1, "outcome":[1.53,7.42,6.34]},
     "2":{"experimentID":2, "seq":1, "outcome":[3.53,6.32,3.34]},
     "shared": {"k1":"v1", "k2":"v2"}
   },
   {
     "1":{"experimentID":1, "seq":2, "outcome":[5.63,4.3]},
     "2":{"experimentID":1, "seq":2, "outcome":[3.65,3.2]},
     "shared": {"k1":"v10", "k2":"v20"}
   }
 ]
```

\#actual
```json
{"experimentID":1, "seq":1, "outcome":[1.53,7.42,6.34], "k1":"v1", "k2":"v2"}
{"experimentID":1, "seq":2, "outcome":[5.63,4.3], "k1":"v10", "k2":"v20"}
```


## Time format

@timeFormat@ - time format directive instructs a validator to convert data into time with specified time format  before actual validation takes place.

Time format is expressed in java style date format.

**Example**

\#expected 

```go
expected := map[string]interface{}{
    "@timeFormat@date": "yyyy-MM-dd",
    "@timeFormat@ts": "yyyy-MM-dd hh:mm:ss"
    "@timeFormat@" "yyyy-MM-dd hh:mm:ss" //default time format       
    "id":123,
    "date": "2019-01-01",
    "ts": "2019-01-01 12:00:01",
}
```

\#actual 

```go
expected := map[string]interface{}{
	"id":123,
    "date": "2019-01-01 12:00:01",,
    "ts": "2019-01-01 12:00:01",
}
```





## Time layout

@timeLayout@ - time format directive instructs a validator to convert data into time with specified time format  before actual validation takes place.

Time layout uses golang time layout.


**Example**

\#expected 

```go
expected := map[string]interface{}{
    "@timeFormat@date": "yyyy-MM-dd",
    "@timeFormat@ts": "yyyy-MM-dd hh:mm:ss"
    "@timeFormat@" "yyyy-MM-dd hh:mm:ss" //default time format       
    "id":123,
    "date": "2019-01-01",
    "ts": "2019-01-01 12:00:01",
}
```

\#actual 

```go
expected := map[string]interface{}{
	"id":123,
    "date": "2019-01-01 12:00:01",,
    "ts": "2019-01-01 12:00:01",
}
```


## Cast data type

@cast@ - instruct a validator to convert data to the specified data type before actual validation takes place.

Supported data type casting:
* int
* float
* boolean

**Example**


\#expected 
 ```json
 [
   {
     "@cast@field1":"float","@cast@field2":"int"
   },
   {
        "field1":2.3,
        "field2":123
   },
   {
      "field1":6.3,
      "field2":551
   }
 ]
```

\#actual
```json
{"field1":"2.3","field2":"123"}
{"field1":"6.3","field2":"551"}
```


## KeyCaseSensitiveDirective

By default map key match is case sensitive, directive allows to disable that behaviours.

## CaseSensitiveDirective

By default text value match is case sensitive, directive allows to disable that behaviours.



## NumericPrecisionPoint

NumericPrecisionPoint controls numeric precision validation comparision

**Example**

\#expected 
 ```json
 [
   {
     "@numericPrecisionPoint@":"7"
   },
   {
       "field1":0.006521405,
        "field2":123
   },
   {
      "field1":0.006521408,
      "field2":551
   }
 ]
```

\#actual
```json
[
   {
       "field1":0.0065214,
        "field2":123
   },
   {
      "field1":0.0065214,
      "field2":551
   }
]
```

## CoalesceWithZero

Coalesce with zero directive sets all nil numeric values to zero 

## Source directive

Source directive is helper directive providing additional information about data point source, i.e. file.json#L113

<a name="Macro"></a>
## Macro and predicates


The macro is an expression with parameters that expands original text value. 
The general format of macro: &lt;ds:MACRO_NAME [json formated array of parameters]>

The following macro are build-in:


| Name | Parameters | Description | Example | 
| --- | --- | --- | --- |
| env | name env variable| Returns value env variable| &lt;ds:env["user"]> |
| nil |n/a| Returns nil value| &lt;ds:nil> |
| cast | type name| Returns value env variable| &lt;ds:cast["int", "123"]> |
| current_timestamp | n/a | Returns time.Now() | &lt;ds:current_timestamp> |
| dob | user age, month, day, format(yyyy-MM-dd as default)  | Returns Date Of Birth| &lt;ds:dob> |

## Predicates

Predicate allows expected value to be evaluated with actual data using custom predicate logic.


| Name | Parameters | Description | Example | 
| --- | --- | --- | --- |
| between | from, to values | Evaluate actual value with between predicate | &lt;ds:between[1.888889, 1.88889]> |
| within_sec | base time, delta, optional date format | Evaluate if actual time is within delta of the base time | &lt;ds:within_sec["now", 6, "yyyyMMdd HH:mm:ss"]> |


**Example**

```go
    expected := `<ds:between[1,10]>`
    actual := 3
```

```go
    expected := `1<ds:env["USER"]>3`,
    actual := fmt.Sprintf("1%v3", os.Getenv("USER"))
```

```go
    expected := `<ds:dob[3, 6, 3>`
    actual := 2015-06-03
```

```go
    expected := `<ds:dob[3, 6, 3,"yyyy-MM-dd"]>`
    actual := 2015-06-03
```


```go
    expected := `<ds:dob[3, 6, 3,"yyyy"]>`
    actual := 2015
```

```go
    expected := `<ds:dob[3, 9, 2,"yyyy-MM"]>`
    actual := 2015-09
```

```go
    expected := `<ds:dob[5, 12, 25,"-MM-dd"]>`
    actual := 12-25
```


<a name="external"></a>
## External resource

- [Advanced Data Testing with Go](https://www.youtube.com/watch?v=ESUJkSMS47k) 
- [Advanced Data Testing with Go Presentation](https://github.com/adrianwit/golang-data-testing/raw/master/advanced_data_testing.pptx)


## GoCover

[![GoCover](https://gocover.io/github.com/viant/assertly)](https://gocover.io/github.com/viant/assertly)


<a name="License"></a>
## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.


<a name="Credits-and-Acknowledgements"></a>

##  Credits and Acknowledgements

**Library Author:** Adrian Witas
