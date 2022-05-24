package gommander

import (
	"strings"
)

type Argument struct {
	name         string
	help         string
	raw          string
	variadic     bool
	is_required  bool
	valid_values []string
}

func NewArgument(name string) *Argument {
	var required bool
	var variadic bool

	if strings.HasPrefix(name, "<") {
		required = true
		name = strings.ReplaceAll(name, "<", "")
		name = strings.ReplaceAll(name, ">", "")
	} else if strings.HasPrefix(name, "[") {
		required = false
		name = strings.ReplaceAll(name, "[", "")
		name = strings.ReplaceAll(name, "]", "")
	}

	if strings.HasSuffix(name, "...") {
		variadic = true
		name = strings.ReplaceAll(name, "...", "")
	}

	return &Argument{
		name:        name,
		is_required: required,
		variadic:    variadic,
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

func (a *Argument) ValidateWith(vals []string) *Argument {
	a.valid_values = vals
	return a
}

func (a *Argument) ValueIsValid(val string) bool {
	for _, v := range a.valid_values {
		if strings.EqualFold(v, val) {
			return true
		}
	}

	return false
}

func new_argument(val string, help string) *Argument {
	var delimiters []string
	var required bool
	var variadic bool

	if strings.HasPrefix(val, "<") {
		delimiters = []string{"<", ">"}
		required = true
	} else if strings.HasPrefix(val, "[") {
		delimiters = []string{"[", "]"}
		required = false
	}

	name := strings.ReplaceAll(val, delimiters[0], "")
	name = strings.ReplaceAll(name, delimiters[1], "")
	name = strings.ReplaceAll(name, "-", "_")

	if strings.HasSuffix(val, "...") {
		variadic = true
		name = strings.ReplaceAll(name, "...", "")
	}

	return &Argument{
		name:        name,
		help:        help,
		raw:         val,
		variadic:    variadic,
		is_required: required,
	}
}
