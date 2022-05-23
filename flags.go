package gommander

import "strings"

type Flag struct {
	name  string
	long  string
	short string
	help  string
}

func NewFlag(name string) *Flag {
	return &Flag{
		name: name,
	}
}

func (f *Flag) Short(val string) *Flag {
	f.short = val
	return f
}

func (f *Flag) Long(val string) *Flag {
	f.long = val
	return f
}

func (f *Flag) Help(val string) *Flag {
	f.help = val
	return f
}

func new_flag(val string, help string) Flag {
	values := strings.Split(val, " ")

	short := ""
	long := ""

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
