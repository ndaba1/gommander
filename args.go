package gommander

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type argumentType string

const (
	integer  argumentType = "int"
	uinteger argumentType = "uint"
	float    argumentType = "float"
	boolean  argumentType = "bool"
	str      argumentType = "str"
	filename argumentType = "file"
)

type Argument struct {
	Name         string
	HelpStr      string
	RawValue     string
	ArgType      argumentType
	IsVariadic   bool
	IsRequired   bool
	ValidValues  []string
	DefaultValue string
	ValidatorFns [](func(string) error)
	ValidatorRe  *regexp.Regexp
}

// A Builder method for creating a new argument. Valid values include <arg>, [arg] or simply the name of the arg
func NewArgument(name string) *Argument {
	arg := Argument{Name: name, ArgType: str}
	var delimiters []string

	if strings.HasPrefix(name, "<") {
		arg.IsRequired = true
		delimiters = []string{"<", ">"}
	} else if strings.HasPrefix(name, "[") {
		delimiters = []string{"[", "]"}
	}

	if len(delimiters) > 0 {
		arg.Name = strings.TrimPrefix(arg.Name, delimiters[0])
		arg.Name = strings.TrimSuffix(arg.Name, delimiters[1])
	}

	if strings.ContainsRune(arg.Name, ':') {
		values := strings.Split(arg.Name, ":")
		arg.Name = values[1]
		arg.ArgType = argumentType(values[0])
		arg.addValidatorFns() // add default validator funcs
	}

	if strings.HasSuffix(arg.Name, "...") {
		arg.IsVariadic = true
		arg.Name = strings.ReplaceAll(arg.Name, "...", "")
	}

	return &arg
}

// A method for setting the default value on an argument to be used when no value is provided but the argument value is required
func (a *Argument) Default(val string) *Argument {
	// Check if value valid
	if len(a.ValidValues) > 0 {
		if !a.testValue(val) {
			// TODO: consider writing to stderr
			fmt.Printf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", a.Name, val, a.ValidValues)
			if !isTestMode() {
				os.Exit(1)
			}
		}
	}
	// verify value against validator fn if any
	for _, fn := range a.ValidatorFns {
		err := fn(val)
		if err != nil {
			fmt.Printf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", a.Name, val)
			if !isTestMode() {
				os.Exit(1)
			}
		}
	}

	a.DefaultValue = val
	return a
}

// Simply sets the description or help string of the given argument
func (a *Argument) Help(val string) *Argument {
	a.HelpStr = val
	return a
}

// Sets whether an argument is variadic or not
func (a *Argument) Variadic(val bool) *Argument {
	a.IsVariadic = val
	return a
}

// Sets whether an argument is required or not
func (a *Argument) Required(val bool) *Argument {
	a.IsRequired = val
	return a
}

func (a *Argument) Type(val argumentType) *Argument {
	a.ArgType = val
	return a
}

// Configures the valid values for an argument
func (a *Argument) ValidateWith(vals []string) *Argument {
	a.ValidValues = vals
	return a
}

// A method to pass a custom validator function for arguments passed
func (a *Argument) ValidatorFunc(fn func(string) error) *Argument {
	a.ValidatorFns = append(a.ValidatorFns, fn)
	return a
}

func (a *Argument) ValidatorRegex(val string) *Argument {
	a.ValidatorRe = regexp.MustCompile(val)
	return a
}

// A method for setting what the argument should be displayed as when printing help
func (a *Argument) DisplayAs(val string) *Argument {
	a.RawValue = val
	return a
}

/****************************** Package utilities ********************************/

func (a *Argument) addValidatorFns() {
	switch a.ArgType {
	case str:
		{
			a.ValidatorFunc(func(s string) error {
				return nil
			})
		}
	case integer:
		{
			a.ValidatorFunc(func(s string) error {
				_, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("`%v` is not a valid integer", s)
				}
				return nil
			})
		}
	case uinteger:
		{
			a.ValidatorFunc(func(s string) error {
				_, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					return fmt.Errorf("`%v` is not a positive integer", s)
				}

				return nil
			})
		}
	case float:
		{
			a.ValidatorFunc(func(s string) error {
				_, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return fmt.Errorf("`%v` is not a valid float", s)
				}
				return nil
			})

		}
	case boolean:
		{
			a.ValidatorFunc(func(s string) error {
				if s != "true" && s != "false" {
					return fmt.Errorf("`%v` is not a valid boolean", s)
				}
				return nil
			})
		}
	case filename:
		{
			a.ValidatorFunc(func(s string) error {
				if _, e := os.Stat(s); e != nil {
					return fmt.Errorf("no such file or directory: `%v`", s)
				}
				return nil
			})
		}

	default:
		{
			fmt.Println(fmt.Errorf("found unknown argument type: `%v` for argument: `%v`", a.ArgType, a.getRawValue()))
			if !isTestMode() {
				os.Exit(1)
			}
		}
	}
}

func (a *Argument) testValue(val string) bool {
	valueMatch := false
	matchCount := 0

	if len(a.ValidValues) == 0 {
		valueMatch = true
	}

	for _, v := range a.ValidValues {
		if strings.EqualFold(v, val) {
			valueMatch = true
			break
		}
	}

	for _, fn := range a.ValidatorFns {
		err := fn(val)
		if err == nil {
			matchCount++
		}
	}

	if a.ValidatorRe != nil && !a.ValidatorRe.MatchString(val) {
		return false
	}

	return valueMatch && matchCount == len(a.ValidatorFns)
}

func (a *Argument) hasDefaultValue() bool {
	return len(a.DefaultValue) > 0
}

func (a *Argument) compare(b *Argument) bool {
	return a.HelpStr == b.HelpStr && a.Name == b.Name && a.getRawValue() == b.getRawValue()
}

func newArgument(val string, help string) *Argument {
	arg := NewArgument(val)
	arg.Help(help)
	return arg
}

func (a *Argument) getRawValue() string {
	if len(a.RawValue) == 0 {
		var value strings.Builder

		write := func(first rune, last rune) {
			value.WriteRune(first)
			value.WriteString(strings.ReplaceAll(a.Name, "_", "-"))
			if a.IsVariadic {
				value.WriteString("...")
			}
			value.WriteRune(last)
		}

		if a.IsRequired {
			write('<', '>')
		} else {
			write('[', ']')
		}
		return value.String()
	}
	return a.RawValue
}

/****************************** Interface implementations ********************************/

func (a *Argument) generate(app *Command) (string, string) {
	var leading strings.Builder
	var floating strings.Builder

	leading.WriteString(a.getRawValue())
	floating.WriteString(a.HelpStr)
	if a.hasDefaultValue() {
		floating.WriteString(fmt.Sprintf(" (default: %v)", a.DefaultValue))
	}

	return leading.String(), floating.String()
}
