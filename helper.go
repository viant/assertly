package assertly

import (
	"github.com/viant/toolbox"
	"strings"
)

func asDataStructure(candidate string) interface{} {

	if isMultiline(candidate) {
		lines := strings.Split(candidate, "\n")
		if toolbox.IsCompleteJSON(lines[0]) && toolbox.IsCompleteJSON(lines[len(lines)-1]) {
			var result = make([]interface{}, 0)
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				item, err := toolbox.JSONToInterface(line)
				if err != nil {
					result = []interface{}{}
					break
				}
				result = append(result, item)
			}
			if len(result) > 0 {
				return result
			}
		} else if toolbox.IsCompleteJSON(candidate) {
			if result, err := toolbox.JSONToInterface(candidate); err == nil {
				return result
			}
		}
	} else if result, err := toolbox.JSONToInterface(candidate); err == nil {
		return result
	}
	return candidate
}

func isMultiline(candidate string) bool {
	return strings.Count(candidate, "\n") > 0
}

func reverseSlice(stringSlice []string) {
	last := len(stringSlice) - 1
	for i := 0; i < len(stringSlice)/2; i++ {
		stringSlice[i], stringSlice[last-i] = stringSlice[last-i], stringSlice[i]
	}
}

func mergeTextMap(source map[string]string, target *map[string]string) {
	if len(source) == 0 {
		return
	}
	if target == nil || len(*target) == 0 {
		*target = make(map[string]string)
	}
	for k := range source {
		(*target)[k] = source[k]
	}
}

func mergeBoolMap(source map[string]bool, target *map[string]bool) {
	if len(source) == 0 {
		return
	}
	if target == nil || len(*target) == 0 {
		*target = make(map[string]bool)
	}
	for k := range source {
		(*target)[k] = source[k]
	}
}

func keysValue(aMap map[string]interface{}, keys ...string) string {
	var result = ""
	for _, key := range keys {
		value := aMap[key]
		result += toolbox.AsString(value)
	}
	return result
}

func keysPairValue(aMap map[string]interface{}, keys ...string) string {
	var result = ""
	for _, key := range keys {
		value := aMap[key]
		if len(result) > 0 {
			result += ""
		}
		result += key + "(" + toolbox.AsString(value) + ")"
	}
	return result
}

func indexSliceBy(aSlice []interface{}, indexFields ...string) map[string]interface{} {
	var result = make(map[string]interface{})
	for _, item := range aSlice {
		var value = keysValue(toolbox.AsMap(item), indexFields...)
		result[value] = item
	}
	return result
}

func toStringSlice(source interface{}) []string {
	if !toolbox.IsSlice(source) {
		return []string{toolbox.AsString(source)}
	}
	var result = make([]string, 0)
	for _, item := range toolbox.AsSlice(source) {
		result = append(result, toolbox.AsString(item))
	}
	return result
}

func isIndexable(source map[string]interface{}) bool {
	for _, v := range source {
		if v == nil {
			continue
		}
		if toolbox.IsMap(v) {
			return true
		}
	}
	return false
}
