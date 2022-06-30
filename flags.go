package gommander

import (
	"fmt"
	"strings"
)

type Flag struct {
	Name     string
	LongVal  string
	ShortVal string
	HelpStr  string
	IsGlobal bool
}

// A Builder method for creating a new flag. It sets the name of the flag and the long version of the flag by appending `--` to the name then returns the flag for further manipulation.
func NewFlag(name string) *Flag {
	return &Flag{
		Name:    name,
		LongVal: fmt.Sprintf("--%v", name),
	}
}

// A method that simply sets the short version of a flag. It takes in a rune and appends a `-` to it then sets that as the short value for the flag
func (f *Flag) Short(val rune) *Flag {
	f.ShortVal = fmt.Sprintf("-%c", val)
	return f
}

// A method for setting the help string or description of the flag
func (f *Flag) Help(val string) *Flag {
	f.HelpStr = val
	return f
}

// A method for setting a flag as global. Global flags are propagated to all the subcommands of a given command
func (f *Flag) Global(val bool) *Flag {
	f.IsGlobal = val
	return f
}

func (f *Flag) compare(b *Flag) bool {
	return f.Name == b.Name && f.ShortVal == b.ShortVal && f.LongVal == b.LongVal && f.HelpStr == b.HelpStr
}

func helpFlag() *Flag {
	return &Flag{
		Name:     "help",
		LongVal:  "--help",
		ShortVal: "-h",
		HelpStr:  "Print out help information",
		IsGlobal: true,
	}
}

func versionFlag() *Flag {
	return &Flag{
		Name:     "version",
		LongVal:  "--version",
		ShortVal: "-v",
		HelpStr:  "Print out version information",
	}
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
		Name:     strings.ReplaceAll(long, "-", ""),
		LongVal:  long,
		ShortVal: short,
		HelpStr:  help,
	}
}

func (f *Flag) generate() (string, string) {
	var leading strings.Builder

	if len(f.ShortVal) > 0 {
		leading.WriteString(fmt.Sprintf("%v,", f.ShortVal))
	} else {
		leading.WriteString("   ")
	}

	if len(f.LongVal) > 0 {
		leading.WriteString(fmt.Sprintf(" %v", f.LongVal))
	}

	return leading.String(), f.HelpStr
}
