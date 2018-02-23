package assertly

import (
	"fmt"
	"strings"
)

//DataPath represents a dat path
type DataPath interface {
	//MatchingPath returns matching path
	MatchingPath() string

	//Path data path
	Path() string

	//Index creates subpath for supplied index
	Index(index int) DataPath

	//Index creates subpath for supplied key
	Key(key string) DataPath

	//Set source for this path, source may represent detail location of data point
	SetSource(string)

	//Get source from this path
	Source() string

	//Match returns a matched directive for this path
	Match(context *Context) *Directive


	//Match returns a directive for this path
	Directive() *Directive

	//Each traverse each data path node upto parent
	Each(callback func(path DataPath) bool)
}

type dataPath struct {
	root      string
	source    string
	index     int
	key       string
	parent    *dataPath
	directive *Directive
}

func (p *dataPath) Index(index int) DataPath {
	return &dataPath{
		index:  index,
		parent: p,
	}
}

func (p *dataPath) Key(field string) DataPath {
	return &dataPath{
		key:    field,
		parent: p,
	}
}


func (p *dataPath) SetSource(source string)  {
	p.source = source;
}


func (p *dataPath) Source() string {
	if p.source != "" {
		return p.source
	}
	var result = ""
	p.each(func(node *dataPath) bool {
		if node.source != "" {
			result = node.source
			return false
		}
		return true
	})
	return result
}


func (p *dataPath) Directive() *Directive {
	return p.directive
}


func (p *dataPath) Match(context *Context) *Directive {
	if p.directive != nil {
		return p.directive
	}

	directive := context.Directives.Match(p)
	p.each(func(node *dataPath) bool {
		if node.directive != nil {
			directive.mergeFrom(node.directive)
			return false
		}
		return true
	})

	p.directive = directive
	return directive
}


func (p *dataPath) Each(callback func(path DataPath) bool) {
	var node = p
	for node != nil {
		if !callback(node) {
			break
		}
		node = node.parent
	}
}


func (p *dataPath) each(callback func(path *dataPath) bool) {
	var node = p
	for node != nil {
		if !callback(node) {
			break
		}
		node = node.parent
	}
}

func (p *dataPath) Path() string {
	var result = make([]string, 0)
	p.each(func(node *dataPath) bool {
		if node.root != "" {
			result = append(result, "["+node.root+"]:")
		} else if node.key != "" {
			var dot = "."
			if node.parent != nil && node.parent.root != "" {
				dot = ""
			}
			result = append(result, dot+node.key)
		} else {
			result = append(result, fmt.Sprintf("[%d]", node.index))
		}
		return true
	})
	reverseSlice(result)
	return strings.Join(result, "")
}

func (p *dataPath) MatchingPath() string {
	var result = make([]string, 0)
	p.each(func(node *dataPath) bool {
		if node.root != "" {
			return false
		}
		if node.key != "" {
			result = append(result, node.key)
		} else {
			result = append(result, "*")
		}
		return true
	})
	reverseSlice(result)
	return strings.Join(result, "/")
}

//NewDataPath returns a new data path.
func NewDataPath(root string) DataPath {
	if root == "" {
		root = " "
	}
	return &dataPath{
		root: root,
	}
}
