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
}

func Program() *Command {
	return NewCommand("").set_is_root(true)
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
	return App(name).set_parent(c)
}

func (c *Command) On(event Event, cb EventCallback) {
	c.emitter.on(event, cb, 0)
}

func (c *Command) BeforeAll(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, -5)
}

func (c *Command) Emit(cfg EventConfig) {
	c.emitter.emit(cfg)
}

func (c *Command) TestEmit() {
	c.Emit(EventConfig{[]string{""}, OutputHelp, c, int(1), c.alias})

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

func (app *Command) PrintHelp() {
	fmt.Printf(app.help)

	fmt.Printf("\n USAGE: \n")
	fmt.Printf("\t .exe [OPTIONS] [COMMAND]")
}
