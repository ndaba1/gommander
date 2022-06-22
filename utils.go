package gommander

import (
	"fmt"
	"strings"
)

type HelpWriter struct{}

func (HelpWriter) Write(c *Command) {
	app := c._getAppRef()

	// TODO: Check settings

	fmter := NewFormatter(app.theme)

	hasArgs := len(c.arguments) > 0
	hasInfo := len(c.discussion) > 0
	hasFlags := len(c.flags) > 0
	hasOptions := len(c.options) > 0
	hasSubcmds := len(c.subCommands) > 0
	hasCustomUsage := len(c.customUsageStr) > 0
	hasSubcmdGroups := len(c.subCmdGroups) > 0

	fmter.Add(Description, fmt.Sprintf("\n%v\n", c.help))
	fmter.section("USAGE")

	if hasCustomUsage {
		fmter.Add(Keyword, fmt.Sprintf("    %v", c.customUsageStr))
	} else {
		fmter.Add(Keyword, fmt.Sprintf("    %v", c._getUsageStr()))

		if hasFlags {
			fmter.Add(Other, " [FLAGS]")
		}

		if hasOptions {
			fmter.Add(Other, " [OPTIONS]")
		}

		if hasArgs {
			fmter.Add(Other, " <ARGS>")
		}

		if hasSubcmds {
			fmter.Add(Other, " <SUBCOMMAND>")
		}
	}
	fmter.close()

	if app.settings[ShowCommandAliases] {
		fmter.section("ALIASES")
		fmter.Add(Description, fmt.Sprintf("    [%v]\n", strings.Join(c.aliases, ", ")))
	}

	if hasArgs {
		fmter.section("ARGS")
		fmter.format(standardize(c.arguments))
	}

	if hasFlags {
		fmter.section("FLAGS")
		fmter.format(standardize(c.flags))
	}

	if hasOptions {
		fmter.section("OPTIONS")
		fmter.format(standardize(c.options))
	}

	if hasSubcmds && !hasSubcmdGroups {
		fmter.section("SUBCOMMANDS")
		fmter.format(standardize(c.subCommands))
	}

	if hasSubcmds && hasSubcmdGroups {
		for k, v := range c.subCmdGroups {
			fmter.section(k)
			fmter.format(standardize(v))
		}
		// TODO: Simplify this logic
		groupContains := func(val *Command) bool {
			included, total := 0, 0
			for _, g := range c.subCmdGroups {
				total++
				if sliceContains(g, val) {
					included--
				} else {
					included++
				}
			}
			return included != total
		}

		otherCmds := []*Command{}
		for _, sc := range c.subCommands {
			if !groupContains(sc) {
				otherCmds = append(otherCmds, sc)
			}
		}

		if len(otherCmds) > 0 {
			fmter.section("Other Commands")
			fmter.format(standardize(otherCmds))
		}
	}

	if hasInfo {
		// TODO: Format discussion here
		fmter.section(strings.ToUpper("discussion"))
		// fmter.discussion(app.discussion, 80)
	}

	fmter.Print()
}

func sliceContains(slice []*Command, val *Command) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

type FormatterType interface {
	*Command | *Flag | *Option | *Argument
	FormatGenerator
}

func standardize[T FormatterType](vals []T) []FormatGenerator {
	values := []FormatGenerator{}
	for _, c := range vals {
		values = append(values, c)
	}
	return values
}

func suggestSubCmd(c *Command, val string) []string {
	var minMatchSize = 3
	var matches []string

	cmdMap := make(map[string]int, 0)

	for _, v := range c.subCommands {
		cmdMap[v.name] = 0
	}

	for i, v := range strings.Split(val, "") {
		for _, sc := range c.subCommands {
			if len(sc.name) > i && string(sc.name[i]) == v {
				cmdMap[sc.name] = cmdMap[sc.name] + 1
			}
		}
	}

	for k, v := range cmdMap {
		if v >= minMatchSize {
			matches = append(matches, k)
		}
	}

	return matches
}
