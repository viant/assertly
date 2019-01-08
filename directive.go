package assertly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
)

const (
	KeyExistsDirective             = "@exists@"
	KeyDoesNotExistsDirective      = "@!exists@"
	TimeFormatDirective            = "@timeFormat@"
	TimeLayoutDirective            = "@timeLayout@"
	SwitchByDirective              = "@switchCaseBy@"
	CastDataTypeDirective          = "@cast@"
	IndexByDirective               = "@indexBy@"
	KeyCaseSensitiveDirective      = "@keyCaseSensitive@"
	CaseSensitiveDirective         = "@caseSensitive@"
	SourceDirective                = "@source@"
	SortTextDirective              = "@sortText@"
	NumericPrecisionPointDirective = "@numericPrecisionPoint@"
	CoalesceWithZeroDirective      = "@coalesceWithZero@"
	AssertPathDirective            = "@assertPath@"
	LengthDirective                = "@length@"
)

type AssertPath struct {
	SubPath  string
	Expected interface{}
}

//Match represents a validation TestDirective
type Directive struct {
	DataPath
	KeyExists             map[string]bool
	KeyDoesNotExist       map[string]bool
	TimeLayout            string
	KeyCaseSensitive      bool
	CaseSensitive         bool
	TimeLayouts           map[string]string
	DataType              map[string]string
	Lengths               map[string]int
	SwitchBy              []string
	CoalesceWithZero      bool
	NumericPrecisionPoint int
	IndexBy               []string
	Source                string
	SortText              bool
	AssertPaths           []*AssertPath
}

func (d *Directive) mergeFrom(source *Directive) {
	mergeTextMap(source.DataType, &d.DataType)
	mergeTextMap(source.TimeLayouts, &d.TimeLayouts)
	mergeBoolMap(source.KeyExists, &d.KeyExists)
	mergeBoolMap(source.KeyDoesNotExist, &d.KeyDoesNotExist)
	if d.MatchingPath() == "" && len(d.IndexBy) == 0 {
		d.IndexBy = source.IndexBy
	}

	if d.NumericPrecisionPoint == 0 {
		d.NumericPrecisionPoint = source.NumericPrecisionPoint
	}
	if d.TimeLayout == "" {
		d.TimeLayout = source.TimeLayout
	}
}

//AddKeyExists adds key exists TestDirective
func (d *Directive) AddSort(key string) {
	if key == SortTextDirective {
		d.SortText = true
	}
}

//AddKeyExists adds key exists TestDirective
func (d *Directive) AddKeyExists(key string) {
	if len(d.KeyExists) == 0 {
		d.KeyExists = make(map[string]bool)
	}
	d.KeyExists[key] = true
}

//AddKeyDoesNotExist adds key does exist TestDirective
func (d *Directive) AddKeyDoesNotExist(key string) {
	if len(d.KeyDoesNotExist) == 0 {
		d.KeyDoesNotExist = make(map[string]bool)
	}
	d.KeyDoesNotExist[key] = true
}

//AddTimeLayout adds time layout TestDirective
func (d *Directive) AddTimeLayout(key, value string) {
	if len(d.TimeLayouts) == 0 {
		d.TimeLayouts = make(map[string]string)
	}
	d.TimeLayouts[key] = value
}

//AddDataType adds data type TestDirective
func (d *Directive) AddDataType(key, value string) {
	if len(d.DataType) == 0 {
		d.DataType = make(map[string]string)
	}
	d.DataType[key] = value
}

//ExtractDataTypes extracts data from from supplied map
func (d *Directive) ExtractDataTypes(aMap map[string]interface{}) {
	for k, v := range aMap {
		if toolbox.IsInt(v) {
			d.AddDataType(k, "int")
		} else if toolbox.IsFloat(v) {
			d.AddDataType(k, "float")
		} else if toolbox.IsBool(v) {
			d.AddDataType(k, "bool")
		} else if toolbox.IsTime(v) {
			if _, has := d.TimeLayouts[k]; !has {
				var dateFormat = "yyyy-MM-dd HH:mm:ss.SSSZ"
				layout := toolbox.DateFormatToLayout(dateFormat)
				d.AddTimeLayout(k, layout)
			}
		}
	}
}

func (d *Directive) asCaseInsensitveMap(aMap map[string]string) map[string]string {
	if len(aMap) == 0 {
		return aMap
	}
	var result = make(map[string]string)
	for k, v := range aMap {
		result[strings.ToUpper(k)] = v
	}
	return result
}

func (d *Directive) ApplyKeyCaseInsensitive() {
	if len(d.IndexBy) > 0 {
		d.IndexBy = strings.Split(strings.ToUpper(strings.Join(d.IndexBy, ",")), ",")
	}
	d.TimeLayouts = d.asCaseInsensitveMap(d.TimeLayouts)
	d.DataType = d.asCaseInsensitveMap(d.DataType)
}

//Add adds by to supplied target
func (d *Directive) Add(target map[string]interface{}) {
	if len(d.SwitchBy) > 0 {
		target[SwitchByDirective] = d.SwitchBy
	}
	if len(d.IndexBy) > 0 {
		target[IndexByDirective] = d.IndexBy
	}

	if d.NumericPrecisionPoint > 0 {
		target[NumericPrecisionPointDirective] = d.NumericPrecisionPoint
	}

	if d.CoalesceWithZero {
		target[CoalesceWithZeroDirective] = d.CoalesceWithZero
	}

	if len(d.DataType) > 0 {
		for k, v := range d.DataType {
			target[CastDataTypeDirective+k] = v
		}
	}
	if len(d.TimeLayouts) > 0 {
		for k, v := range d.TimeLayouts {
			target[TimeLayoutDirective+k] = v
		}
	}
	if d.TimeLayout != "" {
		target[TimeLayoutDirective] = d.TimeLayout
	}
}

func (d *Directive) addAssertPath(subpath string, expected interface{}) {
	d.AssertPaths = append(d.AssertPaths, &AssertPath{
		SubPath:  subpath,
		Expected: expected,
	})
}

//ExtractDirective extract TestDirective from supplied map
func (d *Directive) ExtractDirectives(aMap map[string]interface{}) bool {
	var keyCount = len(aMap)
	var directiveCount = 0

	if len(d.Lengths) == 0 {
		d.Lengths = make(map[string]int)
	}
	for k, v := range aMap {
		if d.IsDirectiveKey(k) {
			directiveCount++
		}

		if k == SwitchByDirective {
			d.SwitchBy = toStringSlice(v)
			continue
		}
		if k == SortTextDirective {
			d.SortText = toolbox.AsBoolean(v)
			continue
		}

		if k == IndexByDirective {
			d.IndexBy = toStringSlice(v)
			continue
		}

		if k == IndexByDirective {
			d.IndexBy = toStringSlice(v)
			continue
		}

		if k == KeyCaseSensitiveDirective {
			d.KeyCaseSensitive = toolbox.AsBoolean(v)
			continue
		}
		if k == CaseSensitiveDirective {
			d.CaseSensitive = toolbox.AsBoolean(v)
			continue
		}

		if k == NumericPrecisionPointDirective {
			d.NumericPrecisionPoint = toolbox.AsInt(v)
			continue
		}

		if k == CoalesceWithZeroDirective {
			d.CoalesceWithZero = toolbox.AsBoolean(v)
			continue
		}

		if k == SourceDirective {
			d.Source = toolbox.AsString(v)
			continue
		}

		if strings.HasPrefix(k, AssertPathDirective) {
			var subPath = strings.Replace(k, AssertPathDirective, "", 1)
			if subPath != "" {
				d.addAssertPath(subPath, v)
			} else if toolbox.IsSlice(v) {
				for _, item := range toolbox.AsSlice(v) {
					if toolbox.IsMap(item) {
						for subPath, expcted := range toolbox.AsMap(item) {
							d.addAssertPath(subPath, expcted)
						}
					}
				}
			} else if toolbox.IsMap(v) {
				for subPath, expcted := range toolbox.AsMap(v) {
					d.addAssertPath(subPath, expcted)
				}
			}
			continue
		}

		if strings.HasPrefix(k, LengthDirective) {
			var key = strings.Replace(k, LengthDirective, "", 1)
			d.Lengths[key] = toolbox.AsInt(v)
			continue
		} else if strings.HasPrefix(k, KeyExistsDirective) {
			var key = strings.Replace(k, KeyExistsDirective, "", 1)
			if toolbox.AsBoolean(v) {
				d.AddKeyExists(key)
			} else {
				d.AddKeyDoesNotExist(key)
			}
			continue
		} else if strings.HasPrefix(k, KeyDoesNotExistsDirective) {
			var key = strings.Replace(k, KeyDoesNotExistsDirective, "", 1)
			if toolbox.AsBoolean(v) {
				d.AddKeyDoesNotExist(key)
			} else {
				d.AddKeyExists(key)
			}
			continue
		}

		if text, ok := v.(string); ok {
			if text == KeyExistsDirective {
				d.AddKeyExists(k)
				continue
			}
			if text == KeyDoesNotExistsDirective {
				d.AddKeyDoesNotExist(k)
				continue
			}

			if strings.HasPrefix(k, TimeFormatDirective) {
				var key = strings.Replace(k, TimeFormatDirective, "", 1)
				if key == "" {
					d.TimeLayout = toolbox.DateFormatToLayout(text)
				} else {
					d.AddTimeLayout(key, toolbox.DateFormatToLayout(text))
				}
				continue
			}

			if strings.HasPrefix(k, TimeLayoutDirective) {
				var key = strings.Replace(k, TimeLayoutDirective, "", 1)
				if key == "" {
					d.TimeLayout = text
				} else {
					d.AddTimeLayout(key, text)
				}
				continue
			}
			if strings.HasPrefix(k, CastDataTypeDirective) {
				var key = strings.Replace(k, CastDataTypeDirective, "", 1)
				d.AddDataType(key, text)
			}
		}
	}
	return keyCount > 0 && keyCount == directiveCount
}

//Apply applies TestDirective to supplied map
func (d *Directive) Apply(aMap map[string]interface{}) error {
	if err := d.applyTimeFormat(aMap); err != nil {
		return err
	}
	if d.NumericPrecisionPoint != 0 {
		aMap[NumericPrecisionPointDirective] = d.NumericPrecisionPoint
	}
	if d.CoalesceWithZero {
		aMap[CoalesceWithZeroDirective] = d.CoalesceWithZero
	}
	if err := d.castData(aMap); err != nil {
		return err
	}
	return nil
}

//DefaultTimeLayout returns default time layout
func (d *Directive) DefaultTimeLayout() string {
	if d.TimeLayout == "" {
		d.TimeLayout = toolbox.DefaultDateLayout
	}
	return d.TimeLayout
}

func (d *Directive) applyTimeFormat(aMap map[string]interface{}) error {
	if len(d.TimeLayouts) == 0 {
		return nil
	}
	for key, layout := range d.TimeLayouts {
		val, ok := aMap[key]
		if !ok || val == nil || getPredicate(val) != nil || toolbox.IsFunc(val) {
			continue
		}
		timeValue, err := toolbox.ToTime(val, layout)
		if err != nil {
			return err
		}
		aMap[key] = timeValue
	}
	return nil
}

func (d *Directive) castData(aMap map[string]interface{}) error {
	if len(d.DataType) == 0 {
		return nil
	}
	for key, dataType := range d.DataType {
		var err error
		var casted interface{}

		val, ok := aMap[key]
		if !ok || val == nil || getPredicate(val) != nil || toolbox.IsFunc(val) {
			continue
		}

		textVal := toolbox.AsString(val)
		if strings.HasPrefix(textVal, "<") || strings.HasSuffix(textVal, ">") {
			continue
		}

		if d.IsDirectiveValue(toolbox.AsString(val)) {
			continue
		}

		if text, ok := val.(string); ok {
			if strings.HasPrefix(text, "!") || strings.HasPrefix(text, "/") || strings.HasPrefix(text, "~") {
				continue
			}
			val = strings.TrimSpace(text)
		}

		switch dataType {
		case "float":
			casted, err = toolbox.ToFloat(val)
		case "int":
			casted, err = toolbox.ToInt(val)

		case "bool":
			casted = toolbox.AsBoolean(val)
		default:
			err = fmt.Errorf("unsupported cast type: %v", dataType)
		}
		if toolbox.IsNilPointerError(err) {
			casted = nil
		} else if err != nil {
			return err
		}
		aMap[key] = casted
	}
	return nil
}

//IsDirectiveKey returns true if key is TestDirective
func (d *Directive) IsDirectiveKey(key string) bool {
	return strings.HasPrefix(key, "@") && strings.Count(key, "@") > 1
}

//IsDirectiveKey returns true if value is TestDirective
func (d *Directive) IsDirectiveValue(value string) bool {
	return value == KeyExistsDirective ||
		value == KeyDoesNotExistsDirective
}

//NewDirective creates a new TestDirective for supplied path
func NewDirective(path DataPath) *Directive {

	dataPath, ok := path.(*dataPath)
	if ok {
		if dataPath.directive != nil {
			return dataPath.directive
		}
	}

	var result = &Directive{
		DataPath:         path,
		KeyCaseSensitive: true,
		CaseSensitive:    true,
		AssertPaths:      make([]*AssertPath, 0),
	}
	if dataPath != nil {
		dataPath.directive = result
	}
	//inherit default time from first ancestor
	path.Each(func(path DataPath) bool {
		directive := path.Directive()
		if directive != nil {
			if directive.TimeLayout != "" {
				result.TimeLayout = directive.TimeLayout
				return false
			}
		}
		return true
	})

	//inherit default numeric precision point
	path.Each(func(path DataPath) bool {
		directive := path.Directive()
		if directive != nil {
			if directive.NumericPrecisionPoint != 0 {
				result.NumericPrecisionPoint = directive.NumericPrecisionPoint
				return false
			}
		}
		return true
	})

	//inherit default numeric precision point
	path.Each(func(path DataPath) bool {
		directive := path.Directive()
		if directive != nil {
			if directive.CoalesceWithZero {
				result.CoalesceWithZero = directive.CoalesceWithZero
				return false
			}
		}
		return true
	})
	return result
}

//TestDirective represents TestDirective record
type TestDirective map[string]interface{}

func (r TestDirective) IndexBy(key string) TestDirective {
	r[IndexByDirective] = key
	return r
}

func IndexBy(key string) TestDirective {
	var result = TestDirective{}
	return result.IndexBy(key)
}

func (r TestDirective) TimeFormat(key, format string) TestDirective {
	r[TimeFormatDirective+key] = format
	return r
}

func TimeFormat(key, format string) TestDirective {
	var result = TestDirective{}
	return result.TimeFormat(key, format)
}

func (r TestDirective) TimeLayout(key, format string) TestDirective {
	r[TimeLayoutDirective+key] = format
	return r
}

func TimeLayout(key, format string) TestDirective {
	var result = TestDirective{}
	return result.TimeLayout(key, format)
}

func (r TestDirective) KeyCaseSensitive() TestDirective {
	r[KeyCaseSensitiveDirective] = true
	return r
}

func (r TestDirective) Cast(field, dataType string) TestDirective {
	r[CastDataTypeDirective+field] = dataType
	return r
}

func (r TestDirective) SortText() TestDirective {
	r[SortTextDirective] = true
	return r
}
