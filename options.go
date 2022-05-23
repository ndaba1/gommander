package gommander

import "strings"

type Option struct {
	name     string
	help     string
	short    string
	long     string
	args     []*Argument
	required bool
}

func new_option(val string, help string, required bool) Option {
	values := strings.Split(val, " ")

	long := ""
	short := ""
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
