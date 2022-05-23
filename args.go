package gommander

import (
	"strings"
)

type Argument struct {
	name        string
	help        string
	raw         string
	variadic    bool
	is_required bool
}

func NewArgument(name string) *Argument {
	return &Argument{
		name:        name,
		is_required: false,
	}
}

func (a *Argument) Help(val string) *Argument {
	a.help = val
	return a
}

func (a *Argument) Variadic(val bool) *Argument {
	a.variadic = val
	return a
}

func (a *Argument) Required(val bool) *Argument {
	a.is_required = val
	// TODO: Depending on whether required or not, infer the correct raw value
	return a
}

func new_argument(val string, help string) *Argument {
	var delimiters []string
	var required bool
	var variadic bool

	// FIXME: Find more robust way for checking
	if strings.ContainsAny(val, "<") {
		delimiters = []string{"<", ">"}
		required = true
	} else if strings.ContainsAny(val, "[") {
		delimiters = []string{"[", "]"}
		required = false
	}

	name := strings.Replace(val, delimiters[0], "", -1)
	name = strings.Replace(name, delimiters[1], "", -1)
	name = strings.Replace(name, "-", "_", -1)

	if strings.HasSuffix(val, "...") {
		variadic = true
		name = strings.Replace(name, "...", "", -1)
	}

	return &Argument{
		name:        name,
		help:        help,
		raw:         val,
		variadic:    variadic,
		is_required: required,
	}
}
