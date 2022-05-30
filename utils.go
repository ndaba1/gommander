package gommander

import (
	"fmt"
	"strings"
)

type HelpWriter struct{}

func (HelpWriter) Write(c *Command) {
	app := c.app_ref

	if c.is_root {
		app = c
	}
	// TODO: Check settings

	fmter := NewFormatter(app.theme)

	has_args := len(c.arguments) > 0
	has_info := len(c.discussion) > 0
	has_flags := len(c.flags) > 0
	has_options := len(c.options) > 0
	has_subcmds := len(c.sub_commands) > 0
	has_custom_usage := len(c.custom_usage_str) > 0
	has_subcmd_groups := len(c.sub_cmd_groups) > 0

	fmter.add(Description, fmt.Sprintf("\n%v\n", c.help))
	fmter.section("USAGE")

	if has_custom_usage {
		fmter.add(Keyword, fmt.Sprintf("    %v", c.custom_usage_str))
	} else {
		fmter.add(Keyword, fmt.Sprintf("    %v", c._get_usage_str()))

		if has_flags {
			fmter.add(Other, " [FLAGS]")
		}

		if has_options {
			fmter.add(Other, " [OPTIONS]")
		}

		if has_args {
			fmter.add(Other, " <ARGS>")
		}

		if has_subcmds {
			fmter.add(Other, " <SUBCOMMAND>")
		}
	}
	fmter.close()

	if has_args {
		fmter.section("ARGS")
		fmter.format(standardize(c.arguments))
	}

	if has_flags {
		fmter.section("FLAGS")
		fmter.format(standardize(c.flags))
	}

	if has_options {
		fmter.section("OPTIONS")
		fmter.format(standardize(c.options))
	}

	if has_subcmds && !has_subcmd_groups {
		fmter.section("SUBCOMMANDS")
		fmter.format(standardize(c.sub_commands))
	}

	if has_subcmds && has_subcmd_groups {
		for k, v := range c.sub_cmd_groups {
			fmter.section(k)
			fmter.format(standardize(v))
		}
		// TODO: Check for `other subcommands`
	}

	if has_info {
		// TODO: Format discussion here
		fmter.section(strings.ToUpper("discussion"))
		// fmter.discussion(app.discussion, 80)
	}

	fmter.print()
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

func suggest_sub_cmd(c *Command, val string) []string {
	var MIN_MATCH_SIZE = 3
	var matches []string

	cmd_map := make(map[string]int, 0)

	for _, v := range c.sub_commands {
		cmd_map[v.name] = 0
	}

	for i, v := range strings.Split(val, "") {
		for _, sc := range c.sub_commands {
			if len(sc.name) > i && string(sc.name[i]) == v {
				cmd_map[sc.name] = cmd_map[sc.name] + 1
			}
		}
	}

	for k, v := range cmd_map {
		if v >= MIN_MATCH_SIZE {
			matches = append(matches, k)
		}
	}

	return matches
}
