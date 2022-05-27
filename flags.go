package gommander

import (
	"fmt"
	"strings"
)

type Flag struct {
	name  string
	long  string
	short string
	help  string
}

func NewFlag(name string) *Flag {
	return &Flag{
		name: name,
		long: fmt.Sprintf("--%v", name),
	}
}

func (f *Flag) Short(val rune) *Flag {
	f.short = fmt.Sprintf("-%c", val)
	return f
}

func (f *Flag) Help(val string) *Flag {
	f.help = val
	return f
}

func (f *Flag) compare(b *Flag) bool {
	return f.name == b.name && f.short == b.short && f.long == b.long && f.help == b.help
}

func new_flag(val string, help string) Flag {
	values := strings.Split(val, " ")

	short, long := "", ""

	for _, v := range values {
		if strings.HasPrefix(v, "--") {
			long = v
		} else if strings.HasPrefix(v, "-") {
			short = v
		}
	}

	return Flag{
		name:  strings.Replace(long, "-", "", -1),
		long:  long,
		short: short,
		help:  help,
	}
}

func (f *Flag) generate() (string, string) {
	var leading strings.Builder

	if len(f.short) > 0 {
		leading.WriteString(fmt.Sprintf("%v,", f.short))
	} else {
		leading.WriteString("   ")
	}

	if len(f.long) > 0 {
		leading.WriteString(fmt.Sprintf(" %v", f.long))
	}

	return leading.String(), f.help
}
