package assertly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
)

const (
	KeyExistsDirective        = "@exists@"
	KeyDoesNotExistsDirective = "@!exists@"
	TimeFormatDirective       = "@timeFormat@"
	TimeLayoutDirective       = "@timeLayout@"
	SwitchByDirective         = "@switchCaseBy@"
	CastDataTypeDirective     = "@cast@"
	IndexByDirective          = "@indexBy@"
	CaseSensitiveDirective    = "@caseSensitive@"
	SourceDirective           = "@source@"
	SortTextDirective         = "@sortText@"
)

//Match represents a validation directive
type Directive struct {
	DataPath
	KeyExists       map[string]bool
	KeyDoesNotExist map[string]bool
	TimeLayout      string
	CaseSensitive   bool
	TimeLayouts     map[string]string
	DataType        map[string]string
	SwitchBy        []string
	IndexBy         []string
	Source          string
	SortText        bool
}

func (d *Directive) mergeFrom(source *Directive) {
	mergeTextMap(source.DataType, &d.DataType)
	mergeTextMap(source.TimeLayouts, &d.TimeLayouts)
	mergeBoolMap(source.KeyExists, &d.KeyExists)
	mergeBoolMap(source.KeyDoesNotExist, &d.KeyDoesNotExist)
	if d.MatchingPath() == "" && len(d.IndexBy) == 0 {
		d.IndexBy = source.IndexBy
	}
	if d.TimeLayout == "" {
		d.TimeLayout = source.TimeLayout
	}
}


//AddKeyExists adds key exists directive
func (d *Directive) AddSort(key string) {
	if key == SortTextDirective {
		d.SortText = true
	}
}

//AddKeyExists adds key exists directive
func (d *Directive) AddKeyExists(key string) {
	if len(d.KeyExists) == 0 {
		d.KeyExists = make(map[string]bool)
	}
	d.KeyExists[key] = true
}

//AddKeyDoesNotExist adds key does exist directive
func (d *Directive) AddKeyDoesNotExist(key string) {
	if len(d.KeyDoesNotExist) == 0 {
		d.KeyDoesNotExist = make(map[string]bool)
	}
	d.KeyDoesNotExist[key] = true
}

//AddTimeLayout adds time layout directive
func (d *Directive) AddTimeLayout(key, value string) {
	if len(d.TimeLayouts) == 0 {
		d.TimeLayouts = make(map[string]string)
	}
	d.TimeLayouts[key] = value
}

//AddDataType adds data type directive
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
			var dateFormat = "yyyy-MM-dd HH:mm:ss.SSSZ"
			layout := toolbox.DateFormatToLayout(dateFormat)
			d.AddTimeLayout(k, layout)
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

func (d *Directive) ApplyCaseInsensitive() {
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

//ExtractDirective extract directive from supplied map
func (d *Directive) ExtractDirectives(aMap map[string]interface{}) bool {
	var keyCount = len(aMap)
	var directiveCount = 0
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

		if k == CaseSensitiveDirective {
			d.CaseSensitive = toolbox.AsBoolean(v)
			continue
		}



		if k == SourceDirective {
			d.Source = toolbox.AsString(v)
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

//Apply applies directive to supplied map
func (d *Directive) Apply(aMap map[string]interface{}) error {
	if err := d.applyTimeFormat(aMap); err != nil {
		return err
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
		if !ok || val == nil {
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
		if !ok || val == nil {
			continue
		}

		if d.IsDirectiveValue(toolbox.AsString(val)) {
			continue
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
		if err != nil {
			return err
		}
		aMap[key] = casted
	}
	return nil
}

//IsDirectiveKey returns true if key is directive
func (d *Directive) IsDirectiveKey(key string) bool {
	return strings.HasPrefix(key, "@") && strings.Count(key, "@") > 1
}

//IsDirectiveKey returns true if value is directive
func (d *Directive) IsDirectiveValue(value string) bool {
	return value == KeyExistsDirective ||
		value == KeyDoesNotExistsDirective
}

//NewDirective creates a new directive for supplied path
func NewDirective(dataPath DataPath) *Directive {
	var result = &Directive{
		DataPath:      dataPath,
		CaseSensitive: true,
	}
	//inherit default time from first ancestor
	dataPath.Each(func(path DataPath) bool {
		directive := path.Directive()
		if directive != nil {
			if directive.TimeLayout != "" {
				result.TimeLayout = directive.TimeLayout
				return false
			}
		}
		return true
	})

	return result
}
