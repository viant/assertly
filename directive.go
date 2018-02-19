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
	SwitchByDirective         = "@switchCaseBy@"
	CastDataTypeDirective     = "@cast@"
	IndexByDirective          = "@indexBy@"
	SourceDirective           = "@source@"
)

//Directive represents a validation directive
type Directive struct {
	DataPath
	KeyExists       map[string]bool
	KeyDoesNotExist map[string]bool
	TimeLayout      string
	TimeLayouts     map[string]string
	DataType        map[string]string
	SwitchBy        []string
	IndexBy         []string
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
			timeValue := toolbox.AsTime(v, toolbox.DefaultDateLayout)
			var dateFormat = "yyyy-MM-dd"
			if timeValue.Hour() > 0 {
				dateFormat += " hh"
				if timeValue.Minute() > 0 {
					dateFormat += ":mm"
					if timeValue.Second() > 0 {
						dateFormat += ":ss"
					}
				}
			}
			layout := toolbox.DateFormatToLayout(dateFormat)
			d.AddTimeLayout(k, layout)
		}
	}
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
			target[TimeFormatDirective+k] = v
		}
	}
	if d.TimeLayout != "" {
		target[TimeFormatDirective] = d.TimeLayout
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

		if k == IndexByDirective {
			d.IndexBy = toStringSlice(v)
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
		if !ok {
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
		if !ok {
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
	return strings.HasPrefix(key, TimeFormatDirective) ||
		strings.HasPrefix(key, CastDataTypeDirective) ||
		key == IndexByDirective ||
		key == SwitchByDirective ||
		key == SourceDirective
}

//IsDirectiveKey returns true if value is directive
func (d *Directive) IsDirectiveValue(value string) bool {
	return value == KeyExistsDirective ||
		value == KeyDoesNotExistsDirective
}

//NewDirective creates a new directive for supplied path
func NewDirective(dataPath DataPath) *Directive {
	var result = &Directive{
		DataPath: dataPath,
	}
	return result
}
