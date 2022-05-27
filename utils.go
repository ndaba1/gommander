package gommander

import (
	"fmt"
)

type HelpWriter struct{}

func (HelpWriter) Write(c *Command) {
	fmter := NewFormatter(c.theme)

	has_args := len(c.arguments) > 0
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
		fmter.add(Keyword, fmt.Sprintf("    %v", c.usage_str))

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
