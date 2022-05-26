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
	var delimiters []string

	if strings.HasPrefix(name, "<") {
		required = true
		delimiters = []string{"<", ">"}
	} else if strings.HasPrefix(name, "[") {
		required = false
		delimiters = []string{"[", "]"}
	}

	if len(delimiters) > 0 {
		name = strings.ReplaceAll(name, delimiters[0], "")
		name = strings.ReplaceAll(name, delimiters[1], "")
	}

	if strings.HasSuffix(name, "...") {
		variadic = true
		name = strings.ReplaceAll(name, "...", "")
	}

	return &Argument{
		name:        strings.ReplaceAll(name, "-", "_"),
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

func (a *Argument) get_raw_value() string {
	if len(a.raw) == 0 {
		var value strings.Builder

		write := func(first rune, last rune) {
			value.WriteRune(first)
			value.WriteString(strings.ReplaceAll(a.name, "_", "-"))
			if a.variadic {
				value.WriteString("...")
			}
			value.WriteRune(last)
		}

		if a.is_required {
			write('<', '>')
		} else {
			write('[', ']')
		}
		return value.String()
	} else {
		return a.raw
	}
}

func (a *Argument) generate() (string, string) {
	leading := a.get_raw_value()

	return leading, a.help
}
