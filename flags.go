package gommander

import (
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
		LongVal: "--" + name,
	}
}

// A method that simply sets the short version of a flag. It takes in a rune and appends a `-` to it then sets that as the short value for the flag
func (f *Flag) Short(val rune) *Flag {
	f.ShortVal = "-" + string(val)
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

func helpFlag() *Flag {
	return &Flag{
		Name:     "help",
		LongVal:  "--help",
		ShortVal: "-h",
		HelpStr:  "Print out help information",
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
	flag := Flag{HelpStr: help}
	values := strings.Split(val, " ")

	for _, v := range values {
		if strings.HasPrefix(v, "--") {
			flag.LongVal = v
		} else if strings.HasPrefix(v, "-") {
			flag.ShortVal = v
		}
	}
	flag.Name = strings.TrimPrefix(flag.LongVal, "--")
	return flag
}

func (f *Flag) generate(app *Command) (string, string) {
	var leading strings.Builder

	if len(f.ShortVal) > 0 {
		leading.WriteString(f.ShortVal + ",")
	} else {
		leading.WriteString("   ")
	}

	if len(f.LongVal) > 0 {
		leading.WriteString(" " + f.LongVal)
	}

	return leading.String(), f.HelpStr
}
