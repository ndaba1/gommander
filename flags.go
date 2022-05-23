package gommander

import "strings"

type Flag struct {
	name  string
	long  string
	short string
	help  string
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
