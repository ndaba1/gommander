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

func NewOption(name string) *Option {
	return &Option{
		name: name,
		long: fmt.Sprintf("--%v", name),
	}
}

func (o *Option) Short(val rune) *Option {
	o.short = fmt.Sprintf("-%c", val)
	return o
}

func (o *Option) Help(val string) *Option {
	o.help = val
	return o
}

func (o *Option) Required(val bool) *Option {
	o.required = val
	return o
}

func (o *Option) Argument(val string) *Option {
	o.args = append(o.args, new_argument(val, ""))
	return o
}

func (o *Option) AddArgument(arg *Argument) *Option {
	o.args = append(o.args, arg)
	return o
}

func (o *Option) compare(p *Option) bool {
	return o.name == p.name && o.help == p.help && o.long == p.long && o.short == p.short && o.required == p.required
}

func new_option(val string, help string, required bool) Option {
	values := strings.Split(val, " ")

	short, long := "", ""
	raw_args := []string{}
	args := []*Argument{}

	for _, v := range values {
		if strings.HasPrefix(v, "--") {
			long = v
		} else if strings.HasPrefix(v, "-") {
			short = v
		} else {
			raw_args = append(raw_args, v)
		}
	}

	for _, a := range raw_args {
		arg := new_argument(a, "")
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
	var leading strings.Builder

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
			leading.WriteString(fmt.Sprintf("%v ", a.get_raw_value()))
		}
	}

	return leading.String(), o.help
}
