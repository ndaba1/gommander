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

func NewArgument(val string, help string) *Argument {
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
