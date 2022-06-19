package gommander

import (
	"fmt"
	"os"
	"strings"
)

type Argument struct {
	name         string
	help         string
	raw          string
	isVariadic   bool
	isRequired   bool
	validValues  []string
	defaultValue string
	validatorFn  func(string) error
}

// A Builder method for creating a new argument. Valid values include <arg>, [arg] or simply the name of the arg
func NewArgument(name string) *Argument {
	required := false
	variadic := false
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
		name:       strings.ReplaceAll(name, "-", "_"),
		isRequired: required,
		isVariadic: variadic,
	}
}

// A method for setting the default value on an argument to be used when no value is provided but the argument value is required
func (a *Argument) Default(val string) *Argument {
	// Check if value valid
	if len(a.validValues) > 0 {
		if !a.testValue(val) {
			fmt.Printf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", a.name, val, a.validValues)
			os.Exit(10)
		}
	}
	// verify value against validator fn if any
	if a.validatorFn != nil {
		err := a.validatorFn(val)
		if err != nil {
			fmt.Printf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", a.name, val)
			os.Exit(10)
		}
	}

	a.defaultValue = val
	return a
}

// Simply sets the description or help string of the given argument
func (a *Argument) Help(val string) *Argument {
	a.help = val
	return a
}

// Sets whether an argument is variadic or not
func (a *Argument) Variadic(val bool) *Argument {
	a.isVariadic = val
	return a
}

// Sets whether an argument is required or not
func (a *Argument) Required(val bool) *Argument {
	a.isRequired = val
	return a
}

// Configures the valid values for an argument
func (a *Argument) ValidateWith(vals []string) *Argument {
	a.validValues = vals
	return a
}

// A method to pass a custom validator function for arguments passed
func (a *Argument) ValidatorFunc(fn func(string) error) *Argument {
	a.validatorFn = fn
	return a
}

// A method for setting what the argument should be displayed as when printing help
func (a *Argument) DisplayAs(val string) *Argument {
	a.raw = val
	return a
}

/****************************** Package utilities ********************************/

func (a *Argument) testValue(val string) bool {
	for _, v := range a.validValues {
		if strings.EqualFold(v, val) {
			return true
		}
	}

	return false
}

func (a *Argument) hasDefaultValue() bool {
	return len(a.defaultValue) > 0
}

func (a *Argument) compare(b *Argument) bool {
	return a.help == b.help && a.name == b.name && a.getRawValue() == b.getRawValue()
}

func newArgument(val string, help string) *Argument {
	arg := NewArgument(val)
	arg.Help(help)
	return arg
}

func (a *Argument) getRawValue() string {
	if len(a.raw) == 0 {
		var value strings.Builder

		write := func(first rune, last rune) {
			value.WriteRune(first)
			value.WriteString(strings.ReplaceAll(a.name, "_", "-"))
			if a.isVariadic {
				value.WriteString("...")
			}
			value.WriteRune(last)
		}

		if a.isRequired {
			write('<', '>')
		} else {
			write('[', ']')
		}
		return value.String()
	}
	return a.raw
}

/****************************** Interface implementations ********************************/

func (a *Argument) generate() (string, string) {
	leading := a.getRawValue()
	var floating strings.Builder
	floating.WriteString(a.help)

	if a.hasDefaultValue() {
		floating.WriteString(fmt.Sprintf(" (default: %v)", a.defaultValue))
	}

	return leading, floating.String()
}
