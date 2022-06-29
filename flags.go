package gommander

import (
	"fmt"
	"strings"
)

type Flag struct {
	name     string
	long     string
	short    string
	help     string
	isGlobal bool
}

// A Builder method for creating a new flag. It sets the name of the flag and the long version of the flag by appending `--` to the name then returns the flag for further manipulation.
func NewFlag(name string) *Flag {
	return &Flag{
		name: name,
		long: fmt.Sprintf("--%v", name),
	}
}

// A method that simply sets the short version of a flag. It takes in a rune and appends a `-` to it then sets that as the short value for the flag
func (f *Flag) Short(val rune) *Flag {
	f.short = fmt.Sprintf("-%c", val)
	return f
}

// A method for setting the help string or description of the flag
func (f *Flag) Help(val string) *Flag {
	f.help = val
	return f
}

// A method for setting a flag as global. Global flags are propagated to all the subcommands of a given command
func (f *Flag) Global(val bool) *Flag {
	f.isGlobal = val
	return f
}

func (f *Flag) compare(b *Flag) bool {
	return f.name == b.name && f.short == b.short && f.long == b.long && f.help == b.help
}

func helpFlag() *Flag {
	return NewFlag("help").Short('h').Help("Print out help information")
}

func versionFlag() *Flag {
	return NewFlag("version").Short('v').Help("Print out version information")
}

func newFlag(val string, help string) Flag {
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
		name:  strings.ReplaceAll(long, "-", ""),
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
