package assertly

import (
	"fmt"
	"strings"
)

type DataPath interface {
	MatchingPath() string

	Path() string

	Index(index int) DataPath

	Key(key string) DataPath

	Directive(context *Context) *Directive
}

type dataPath struct {
	root      string
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

func (p *dataPath) Directive(context *Context) *Directive {
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

func NewDataPath(root string) DataPath {
	if root == "" {
		root = " "
	}
	return &dataPath{
		root: root,
	}
}
