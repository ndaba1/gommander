package gommander

import (
	"fmt"
	"os"
	"strings"
)

type Argument struct {
	name          string
	help          string
	raw           string
	is_variadic   bool
	is_required   bool
	valid_values  []string
	default_value string
}

// A Builder method for creating a new argument. Valid values include <arg>, [arg] or simply the name of the arg
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
		is_variadic: variadic,
	}
}

// A method for setting the default value on an argument to be used when no value is provided but the argument value is required
func (a *Argument) Default(val string) *Argument {
	// Check if value valid
	if len(a.valid_values) > 0 {
		if !a.test_value(val) {
			fmt.Printf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", a.name, val, a.valid_values)
			os.Exit(10)
		}
	}
	a.default_value = val
	return a
}

// Simply sets the description or help string of the given argument
func (a *Argument) Help(val string) *Argument {
	a.help = val
	return a
}

// Sets whether an argument is variadic or not
func (a *Argument) Variadic(val bool) *Argument {
	a.is_variadic = val
	return a
}

// Sets whether an argument is required or not
func (a *Argument) Required(val bool) *Argument {
	a.is_required = val
	return a
}

// Configures the valid values for an argument
func (a *Argument) ValidateWith(vals []string) *Argument {
	a.valid_values = vals
	return a
}

/****************************** Package utilities ********************************/

func (a *Argument) test_value(val string) bool {
	for _, v := range a.valid_values {
		if strings.EqualFold(v, val) {
			return true
		}
	}

	return false
}

func (a *Argument) has_default_value() bool {
	return len(a.default_value) > 0
}

func (a *Argument) compare(b *Argument) bool {
	return a.help == b.help && a.name == b.name && a.get_raw_value() == b.get_raw_value()
}

func new_argument(val string, help string) *Argument {
	arg := NewArgument(val)
	arg.Help(help)
	return arg
}

func (a *Argument) get_raw_value() string {
	if len(a.raw) == 0 {
		var value strings.Builder

		write := func(first rune, last rune) {
			value.WriteRune(first)
			value.WriteString(strings.ReplaceAll(a.name, "_", "-"))
			if a.is_variadic {
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

/****************************** Interface implementations ********************************/

func (a *Argument) generate() (string, string) {
	leading := a.get_raw_value()

	return leading, a.help
}
