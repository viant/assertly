package assertly

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
)

//Context represent validation context
type Context struct {
	toolbox.Context
	state      data.Map
	Directives *Directives
}

func NewContext() *Context {
	return &Context{
		Context:    toolbox.NewContext(),
		state:      data.NewMap(),
		Directives: NewDirectives(),
	}
}
