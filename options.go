package gommander

import (
	"fmt"
	"strings"
)

type Option struct {
	Name       string
	HelpStr    string
	ShortVal   string
	LongVal    string
	Arg        *Argument
	IsRequired bool
}

// A builder method to generate a new option
func NewOption(name string) *Option {
	return &Option{
		Name:    name,
		LongVal: fmt.Sprintf("--%v", name),
	}
}

// Simply sets the shorthand version of the option
func (o *Option) Short(val rune) *Option {
	o.ShortVal = fmt.Sprintf("-%c", val)
	return o
}

// A method for setting the help string / description for an option
func (o *Option) Help(val string) *Option {
	o.HelpStr = val
	return o
}

// Sets whether or not the option is required
func (o *Option) Required(val bool) *Option {
	o.IsRequired = val
	return o
}

// A method for adding a new argument to an option. Takes as input the name of the argument
func (o *Option) Argument(val string) *Option {
	o.AddArgument(newArgument(val, ""))
	return o
}

// A builder method for adding an argument. Expects an instance of an argument as input
func (o *Option) AddArgument(arg *Argument) *Option {
	id := fmt.Sprintf("opt-%s-arg-%s", o.Name, arg.Name)
	if !cache[id] {
		cache[id] = true
		o.Arg = arg
	}
	return o
}

func (o *Option) compare(p *Option) bool {
	return o.Name == p.Name && o.HelpStr == p.HelpStr && o.LongVal == p.LongVal && o.ShortVal == p.ShortVal && o.IsRequired == p.IsRequired
}

func newOption(val string, help string, required bool) Option {
	opt := Option{HelpStr: help, IsRequired: required}
	values := strings.Split(val, " ")

	for _, v := range values {
		if strings.HasPrefix(v, "--") {
			opt.LongVal = v
		} else if strings.HasPrefix(v, "-") {
			opt.ShortVal = v
		} else {
			opt.Arg = newArgument(v, "")
		}
	}
	opt.Name = strings.TrimPrefix(opt.LongVal, "--")

	return opt
}

func (o *Option) generate() (string, string) {
	var leading, floating strings.Builder

	if len(o.ShortVal) > 0 {
		leading.WriteString(fmt.Sprintf("%v,", o.ShortVal))
	} else {
		leading.WriteString("   ")
	}

	if len(o.LongVal) > 0 {
		leading.WriteString(fmt.Sprintf(" %v ", o.LongVal))
	}

	if o.Arg != nil {
		leading.WriteString(fmt.Sprintf("%v ", o.Arg.getRawValue()))
	}

	floating.WriteString(o.HelpStr)
	if o.Arg != nil && o.Arg.hasDefaultValue() {
		floating.WriteString(fmt.Sprintf(" (default: %v)", o.Arg.DefaultValue))
	}

	return leading.String(), floating.String()
}
