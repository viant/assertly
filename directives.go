package assertly


//Directives represent a directive
type Directives struct {
	*Directive
	PathDirectives map[string]*Directive
}

func (d *Directives) Match(path DataPath) *Directive {
	var result = NewDirective(path)
	result.mergeFrom(d.Directive)

	if matched, ok := d.PathDirectives[path.MatchingPath()]; ok {
		result.mergeFrom(matched)
	}
	return result
}

//NewDirectives returns new directives
func NewDirectives(directives ...*Directive) *Directives {
	var result = &Directives{
		Directive:      NewDirective(NewDataPath("")),
		PathDirectives: make(map[string]*Directive),
	}
	for i, directive := range directives {
		if directive.MatchingPath() == "" {
			result.Directive = directive
			continue
		}
		result.PathDirectives[directive.MatchingPath()] = directives[i]
	}
	return result
}
