package gommander

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type argumentType string

const (
	integer = "int"
	float   = "float"
	boolean = "bool"
	str     = "str"
)

type Argument struct {
	name         string
	help         string
	raw          string
	argType      argumentType
	isVariadic   bool
	isRequired   bool
	validValues  []string
	defaultValue string
	validatorFns [](func(string) error)
}

// A Builder method for creating a new argument. Valid values include <arg>, [arg] or simply the name of the arg
func NewArgument(name string) *Argument {
	required := false
	variadic := false
	_type := str
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

	if strings.ContainsRune(name, ':') {
		values := strings.Split(name, ":")

		_type = values[0]
		name = values[1]
	}

	if strings.HasSuffix(name, "...") {
		variadic = true
		name = strings.ReplaceAll(name, "...", "")
	}

	arg := Argument{
		name:       strings.ReplaceAll(name, "-", "_"),
		isRequired: required,
		isVariadic: variadic,
		argType:    argumentType(_type),
	}
	arg.addValidatorFns() // add default validator funcs
	return &arg
}

// A method for setting the default value on an argument to be used when no value is provided but the argument value is required
func (a *Argument) Default(val string) *Argument {
	// Check if value valid
	if len(a.validValues) > 0 {
		if !a.testValue(val) {
			// TODO: consider writing to stderr
			fmt.Printf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", a.name, val, a.validValues)
			if !isTestMode() {
				os.Exit(1)
			}
		}
	}
	// verify value against validator fn if any
	for _, fn := range a.validatorFns {
		err := fn(val)
		if err != nil {
			fmt.Printf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", a.name, val)
			if !isTestMode() {
				os.Exit(1)
			}
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

func (a *Argument) Type(val argumentType) *Argument {
	a.argType = val
	return a
}

// Configures the valid values for an argument
func (a *Argument) ValidateWith(vals []string) *Argument {
	a.validValues = vals
	return a
}

// A method to pass a custom validator function for arguments passed
func (a *Argument) ValidatorFunc(fn func(string) error) *Argument {
	a.validatorFns = append(a.validatorFns, fn)
	return a
}

// A method for setting what the argument should be displayed as when printing help
func (a *Argument) DisplayAs(val string) *Argument {
	a.raw = val
	return a
}

/****************************** Package utilities ********************************/

func (a *Argument) addValidatorFns() {
	switch a.argType {
	case integer:
		{
			a.ValidatorFunc(func(s string) error {
				_, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("%v is not a valid integer", s)
				}
				return nil
			})
		}
	case float:
		{
			a.ValidatorFunc(func(s string) error {
				_, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return fmt.Errorf("%v is not a valid float", s)
				}
				return nil
			})

		}
	case boolean:
		{
			a.ValidatorFunc(func(s string) error {
				if s != "true" && s != "false" {
					return fmt.Errorf("%v is not a valid boolean", s)
				}
				return nil
			})
		}
	case str:
		{
			a.ValidatorFunc(func(s string) error {
				return nil
			})
		}
	default:
		{
			fmt.Println(fmt.Errorf("found unknown argument type: `%v` for argument: `%v`", a.argType, a.getRawValue()))
			if !isTestMode() {
				os.Exit(1)
			}
		}
	}
}

func (a *Argument) testValue(val string) bool {
	valueMatch := false
	matchCount := 0

	if len(a.validValues) == 0 {
		valueMatch = true
	}

	for _, v := range a.validValues {
		if strings.EqualFold(v, val) {
			valueMatch = true
			break
		}
	}

	for _, fn := range a.validatorFns {
		err := fn(val)
		if err == nil {
			matchCount++
		}
	}

	return valueMatch && matchCount == len(a.validatorFns)
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
