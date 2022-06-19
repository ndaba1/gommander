package gommander

import (
	"fmt"
	"strings"
)

type Option struct {
	name     string
	help     string
	short    string
	long     string
	args     []*Argument
	required bool
}

// A builder method to generate a new option
func NewOption(name string) *Option {
	return &Option{
		name: name,
		long: fmt.Sprintf("--%v", name),
	}
}

// Simply sets the shorthand version of the option
func (o *Option) Short(val rune) *Option {
	o.short = fmt.Sprintf("-%c", val)
	return o
}

// A method for setting the help string / description for an option
func (o *Option) Help(val string) *Option {
	o.help = val
	return o
}

// Sets whether or not the option is required
func (o *Option) Required(val bool) *Option {
	o.required = val
	return o
}

// A method for adding a new argument to an option. Takes as input the name of the argument
func (o *Option) Argument(val string) *Option {
	o.args = append(o.args, newArgument(val, ""))
	return o
}

// A builder method for adding an argument. Expects an instance of an argument as input
func (o *Option) AddArgument(arg *Argument) *Option {
	o.args = append(o.args, arg)
	return o
}

func (o *Option) compare(p *Option) bool {
	return o.name == p.name && o.help == p.help && o.long == p.long && o.short == p.short && o.required == p.required
}

func newOption(val string, help string, required bool) Option {
	values := strings.Split(val, " ")

	short, long := "", ""
	rawArgs := []string{}
	args := []*Argument{}

	for _, v := range values {
		if strings.HasPrefix(v, "--") {
			long = v
		} else if strings.HasPrefix(v, "-") {
			short = v
		} else {
			rawArgs = append(rawArgs, v)
		}
	}

	for _, a := range rawArgs {
		arg := newArgument(a, "")
		args = append(args, arg)
	}

	return Option{
		name:     strings.Replace(long, "-", "", -1),
		help:     help,
		long:     long,
		short:    short,
		args:     args,
		required: required,
	}
}

func (o *Option) generate() (string, string) {
	var leading, floating strings.Builder

	if len(o.short) > 0 {
		leading.WriteString(fmt.Sprintf("%v,", o.short))
	} else {
		leading.WriteString("   ")
	}

	if len(o.long) > 0 {
		leading.WriteString(fmt.Sprintf(" %v ", o.long))
	}

	if len(o.args) > 0 {
		for _, a := range o.args {
			leading.WriteString(fmt.Sprintf("%v ", a.getRawValue()))
		}
	}

	floating.WriteString(o.help)
	if len(o.args) == 1 && o.args[0].hasDefaultValue() {
		floating.WriteString(fmt.Sprintf(" (default: %v)", o.args[0].defaultValue))
	}

	return leading.String(), floating.String()
}
