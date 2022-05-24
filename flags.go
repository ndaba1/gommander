package gommander

import (
	"strings"
)

type Flag struct {
	name  string
	long  string
	short string
	help  string
}

func NewFlag(name string) *Flag {
	// TODO(vndaba): append `--` more efficiently
	temp := []string{"--"}
	temp = append(temp, name)

	return &Flag{
		name: name,
		long: strings.Join(temp, ""),
	}
}

func (f *Flag) Short(val rune) *Flag {
	temp := []string{"-"}
	temp = append(temp, string(val))

	f.short = strings.Join(temp, "")
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
