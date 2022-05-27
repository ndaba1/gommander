package gommander

import (
	"fmt"
	"os"
	"strings"
)

type CommandCallback = func(*ParserMatches)

type Command struct {
	alias            string
	arguments        []*Argument
	author           string
	callback         CommandCallback
	discussion       string
	emitter          EventEmitter
	flags            []*Flag
	help             string
	is_root          bool
	name             string
	options          []*Option
	parent           *Command
	sub_commands     []*Command
	settings         Settings
	global_settings  *Settings
	theme            Theme
	version          string
	usage_str        string
	custom_usage_str string
	sub_cmd_groups   map[string][]*Command
}

func App() *Command {
	return NewCommand("").set_is_root(true).AddFlag(
		NewFlag("version").
			Short('v').
			Help("Print out version information"),
	)
}

func NewCommand(name string) *Command {
	return &Command{
		name:            name,
		arguments:       []*Argument{},
		flags:           []*Flag{NewFlag("help").Short('h').Help("Print out help information")},
		options:         []*Option{},
		sub_commands:    []*Command{},
		parent:          nil,
		settings:        Settings{},
		theme:           DefaultTheme(),
		is_root:         false,
		emitter:         new_emitter(),
		global_settings: &Settings{},
		usage_str:       name,
	}
}

// Value getters
func (c *Command) GetAlias() string           { return c.alias }
func (c *Command) GetArguments() []*Argument  { return c.arguments }
func (c *Command) GetFlags() []*Flag          { return c.flags }
func (c *Command) GetHelp() string            { return c.help }
func (c *Command) GetName() string            { return c.name }
func (c *Command) GetOptions() []*Option      { return c.options }
func (c *Command) GetParent() *Command        { return c.parent }
func (c *Command) GetSubCommands() []*Command { return c.sub_commands }
func (c *Command) GetVersion() string         { return c.version }
func (c *Command) GetUsageStr() string        { return c.usage_str }

// Value setters

func (c *Command) Action(cb CommandCallback) *Command {
	c.callback = cb
	return c
}

func (c *Command) Alias(alias string) *Command {
	c.alias = alias
	return c
}

func (c *Command) Argument(val string, help string) *Command {
	argument := new_argument(val, help)
	c.arguments = append(c.arguments, argument)
	return c
}

func (c *Command) AddFlag(flag *Flag) *Command {
	for _, f := range c.flags {
		if f.short == flag.short {
			return c
		}
	}
	c.flags = append(c.flags, flag)
	return c
}

func (c *Command) AddOption(opt *Option) *Command {
	for _, o := range c.options {
		if o.short == opt.short {
			return c
		}
	}
	c.options = append(c.options, opt)
	return c
}

func (c *Command) Flag(val string, help string) *Command {
	flag := new_flag(val, help)
	return c.AddFlag(&flag)
}

func (c *Command) Info(info string) *Command {
	c.discussion = info
	return c
}

func (c *Command) Help(help string) *Command {
	c.help = help
	return c
}

func (c *Command) Name(name string) *Command {
	c.name = name
	return c
}

func (c *Command) Version(version string) *Command {
	c.version = version
	return c
}

func (c *Command) UsageStr(val string) *Command {
	c.custom_usage_str = val
	return c
}

func (c *Command) Option(val string, help string) *Command {
	option := new_option(val, help, false)
	return c.AddOption(&option)
}

func (c *Command) RequiredOption(val string, help string) *Command {
	opt := new_option(val, help, true)
	return c.AddOption(&opt)
}

func (c *Command) Subcommand(name string) *Command {
	sub_cmd := NewCommand(name).set_parent(c)

	cmd_path := []string{c.usage_str, sub_cmd.usage_str}
	c.sub_commands = append(c.sub_commands, sub_cmd)
	sub_cmd.usage_str = strings.Join(cmd_path, " ")

	return sub_cmd
}

func (c *Command) Parse() {
	parser := NewParser(c)
	matches, err := parser.parse(os.Args[1:])
	if err != nil {
		println(err.Error())
	}
	val := fmt.Sprintf("%v", matches)
	fmt.Print(val)
}

// Event emitters functionality

func (c *Command) On(event Event, cb EventCallback) {
	c.emitter.on(event, cb, 0)
}

func (c *Command) BeforeAll(cb EventCallback) {
	c.emitter.insert_before_all(cb)
}

func (c *Command) AfterAll(cb EventCallback) {
	c.emitter.insert_after_all(cb)
}

func (c *Command) BeforeHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, -4)
}

func (c *Command) AfterHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, 4)
}

func (c *Command) Emit(cfg EventConfig) {
	c.emitter.emit(cfg)
}

// Interior utility functions
func (app *Command) set_parent(parent *Command) *Command {
	app.parent = parent
	return app
}

func (app *Command) set_is_root(val bool) *Command {
	app.is_root = val
	return app
}

func (c *Command) PrintHelp() {
	fmter := NewFormatter()

	fmter.add(Description, fmt.Sprintf("%v\n", c.help))

	has_args := len(c.arguments) > 0
	has_flags := len(c.flags) > 0
	has_options := len(c.options) > 0
	has_subcmds := len(c.sub_commands) > 0
	has_custom_usage := len(c.custom_usage_str) > 0
	has_subcmd_groups := len(c.sub_cmd_groups) > 0

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
		args := []FormatGenerator{}
		for _, a := range c.arguments {
			args = append(args, a)
		}
		fmter.format(args)
	}

	if has_flags {
		fmter.section("FLAGS")
		flags := []FormatGenerator{}
		for _, f := range c.flags {
			flags = append(flags, f)
		}
		fmter.format(flags)
	}

	if has_options {
		fmter.section("OPTIONS")
		opts := []FormatGenerator{}
		for _, o := range c.options {
			opts = append(opts, o)
		}
		fmter.format(opts)
	}

	if has_subcmds && !has_subcmd_groups {
		fmter.section("SUBCOMMANDS")
		subcmds := []FormatGenerator{}
		for _, c := range c.sub_commands {
			subcmds = append(subcmds, c)
		}
		fmter.format(subcmds)
	}

	if has_subcmds && has_subcmd_groups {
		for k, v := range c.sub_cmd_groups {
			fmter.section(k)
			subcmds := []FormatGenerator{}
			for _, c := range v {
				subcmds = append(subcmds, c)
			}
			fmter.format(subcmds)
		}
	}

	fmter.print()
}

/****************************** Interface Implementations ****************************/

func (c *Command) generate() (string, string) {
	// TODO: Check if allow command aliases
	return c.GetName(), c.GetHelp()
}
