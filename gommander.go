package gommander

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	settings         AppSettings
	global_settings  *AppSettings
	theme            Theme
	version          string
	usage_str        string
	custom_usage_str string
	sub_cmd_groups   map[string][]*Command
	app_ref          *Command
}

func App() *Command {
	return NewCommand("")._set_is_root().AddFlag(
		NewFlag("version").
			Short('v').
			Help("Print out version information"),
	).Theme(DefaultTheme())
}

func NewCommand(name string) *Command {
	return &Command{
		name:            name,
		arguments:       []*Argument{},
		flags:           []*Flag{NewFlag("help").Short('h').Help("Print out help information")},
		options:         []*Option{},
		sub_commands:    []*Command{},
		parent:          nil,
		settings:        AppSettings{},
		is_root:         false,
		emitter:         new_emitter(),
		global_settings: &AppSettings{},
		usage_str:       name,
		sub_cmd_groups:  make(map[string][]*Command),
	}
}

/****************************** Value Getters ****************************/

// Simply returns the alias of the command or an empty string
func (c *Command) GetAlias() string { return c.alias }

// Returns the author of the program if one is set
func (c *Command) GetAuthor() string { return c.author }

// Returns a slice of the configured arguments for a command
func (c *Command) GetArguments() []*Argument { return c.arguments }

// Returns a slice of the configured flags
func (c *Command) GetFlags() []*Flag { return c.flags }

// Returns the help string / description that gets printed out on help
func (c *Command) GetHelp() string { return c.help }

// Returns the configured name of a command
func (c *Command) GetName() string { return c.name }

// Returns the slice of options belonging to a command
func (c *Command) GetOptions() []*Option { return c.options }

// Returns the parent of a command or nil if none is found
func (c *Command) GetParent() *Command { return c.parent }

// Returns a slice of subcommands chained to the command instance
func (c *Command) GetSubCommands() []*Command { return c.sub_commands }

// Returns the version of the program
func (c *Command) GetVersion() string { return c.version }

// Returns the default usage string or a custom_usage_str if one exists
func (c *Command) GetUsageStr() string {
	if len(c.custom_usage_str) > 0 {
		return c.custom_usage_str
	} else {
		return c.usage_str
	}
}

/****************************** Command Metadata Setters ****************************/

// This method set the callback to be excuted when a command is matched
func (c *Command) Action(cb CommandCallback) *Command {
	c.callback = cb
	return c
}

// A method for adding a flag to a command. It is similar to the `.Flag()` method except this method receives an instance of an already created flag while `.Flag()` receives a string, creates a flag from it and call this method internally
func (c *Command) AddFlag(flag *Flag) *Command {
	for _, f := range c.flags {
		if f.short == flag.short {
			return c
		}
	}
	c.flags = append(c.flags, flag)
	return c
}

// A method for adding a new option to a command. The `.Option()` method invokes this one internally. Identical to the `.AddFlag()` method except this one is for options instead of flags
func (c *Command) AddOption(opt *Option) *Command {
	for _, o := range c.options {
		if o.short == opt.short {
			return c
		}
	}
	c.options = append(c.options, opt)
	return c
}

// Simply sets the alias of a command
func (c *Command) Alias(alias string) *Command {
	c.alias = alias
	return c
}

func (c *Command) AddArgument(arg *Argument) *Command {
	c.arguments = append(c.arguments, arg)
	return c
}

// A method for setting any expected arguments for a command, it takes in the value of the argument e.g. `<image-name>` and the help string for said argument
func (c *Command) Argument(val string, help string) *Command {
	argument := new_argument(val, help)
	c.AddArgument(argument)
	return c
}

// Simply sets the author of the program, usually invoked on the root command
func (c *Command) Author(val string) *Command {
	c.author = val
	return c
}

// Receives a string representing the flag structure and the flag help string and creates a new flag from it. Acceptable values include:
// ("-h --help", "A help flag")
// You could also omit the short or long version of the flag
func (c *Command) Flag(val string, help string) *Command {
	flag := new_flag(val, help)
	return c.AddFlag(&flag)
}

// Used to set more information or the command discussion which gets printed out when help is invoked, at the bottom most section
func (c *Command) Info(info string) *Command {
	c.discussion = info
	return c
}

// Simply sets the help string, otherwise known as description of a command
func (c *Command) Help(help string) *Command {
	c.help = help
	return c
}

// Sets the name of a command, and updates the usage str as well
func (c *Command) Name(name string) *Command {
	c.name = name
	c.usage_str = name
	return c
}

// Sets the version of a command, usually the entry point command(App)
func (c *Command) Version(version string) *Command {
	c.version = version
	return c
}

// An identical method to the `.Flag()` method but for options. Expected syntax: "-p --port <port-number>"
func (c *Command) Option(val string, help string) *Command {
	option := new_option(val, help, false)
	return c.AddOption(&option)
}

// This method is used to mark an option as required for a given command. Another way of achieving this is using the `.AddOption()` method and using the `NewOption()` builder interface to define option parameters
func (c *Command) RequiredOption(val string, help string) *Command {
	opt := new_option(val, help, true)
	return c.AddOption(&opt)
}

// Used to define a custom usage string. If one is present, it will be used instead of the default one
func (c *Command) UsageStr(val string) *Command {
	c.custom_usage_str = val
	return c
}

/****************************** Subcommand related methods ****************************/

// When chained on a command, this method adds said command to the provided sub_cmd group in the parent of the command.
func (c *Command) AddToSubCommandGroup(name string) *Command {
	c.parent.sub_cmd_groups[name] = append(c.parent.sub_cmd_groups[name], c)
	return c
}

// Receives a reference to a command, sets the command parent and usage string then adds its to the slice of subcommands. This method is called internally by the `.SubCommand()` method but users can also invoke it directly
func (c *Command) AddSubCommand(sub_cmd *Command) *Command {
	sub_cmd._set_parent(c)
	cmd_path := []string{c.usage_str, sub_cmd.usage_str}
	c.sub_commands = append(c.sub_commands, sub_cmd)
	sub_cmd.usage_str = strings.Join(cmd_path, " ")

	// propagate theme
	sub_cmd.theme = c.theme

	if c.is_root {
		sub_cmd.app_ref = c
	} else {
		sub_cmd.app_ref = c.app_ref
	}

	return c
}

// An easier method for creating sub_cmds while avoiding too much function paramets nesting. It accepts the name of the new sub_cmd and returns the newly created sub_cmd
func (c *Command) SubCommand(name string) *Command {
	sub_cmd := NewCommand(name)
	c.AddSubCommand(sub_cmd)
	return sub_cmd
}

// A manual way of creating a new subcommand group and adding the desired commands to it
func (c *Command) SubCommandGroup(name string, vals []*Command) {
	c.sub_cmd_groups[name] = append(c.sub_cmd_groups[name], vals...)
}

/****************************** Settings ****************************/

func (c *Command) _init() {
	// TODO: Check if override default listeners, add help subcmd
	if c.settings[IncludeHelpSubcommand] && len(c.sub_commands) > 0 {
		valid_subcmds := []string{}

		for _, c := range c.sub_commands {
			valid_subcmds = append(valid_subcmds, c.name)
		}

		c.SubCommand("help").
			Help("Print out help information for the passed command").
			AddArgument(
				NewArgument("<COMMAND>").
					Help("The name of the command to output help for").
					ValidateWith(valid_subcmds),
			).
			Action(func(pm *ParserMatches) {
				val, _, _ := pm.GetArgValue("<COMMAND>")
				parent := pm.matched_cmd.parent

				if parent != nil {
					cmd, _ := parent.find_subcommand(val)
					cmd.PrintHelp()
				}
			})
	}

	// Default help listener cannot be overridden
	c.emitter.on(OutputHelp, func(ec *EventConfig) {
		cmd := ec.matched_cmd
		cmd.PrintHelp()
	}, -4)

	if !c.settings[OverrideAllDefaultListeners] {
		c.emitter.on_errors(func(ec *EventConfig) {
			err := ec.err
			// TODO: Match theme in better way
			err.Display(c.theme)
		})

		c.emitter.on(OutputVersion, func(ec *EventConfig) {
			// TODO: Print version in a better way
			app := ec.app_ref

			fmt.Println(app.GetName(), app.GetVersion())
			fmt.Println(app.GetAuthor())
			fmt.Println(app.GetHelp())
		}, -4)

		for _, event := range c.emitter.events_to_override {
			c.emitter.rm_default_lstnr(event)
		}
	}
}

// A method for configuring the settings of a command
func (c *Command) Set(s Setting, value bool) *Command {
	c.settings[s] = value
	return c
}

// A method for configuring the theme of a command
func (c *Command) Theme(value Theme) *Command {
	c.theme = value
	return c
}

// A method for configuring a command to use a package-predefined theme when printing output
func (c *Command) UsePredefinedTheme(value PredefinedTheme) *Command {
	c.Theme(GetPredefinedTheme(value))
	return c
}

/****************************** Parser Functionality ****************************/

func (c *Command) _isExpectingValues() bool {
	return len(c.sub_commands) > 0 || len(c.arguments) > 0
}

func (c *Command) Parse() {
	// TODO: Init/build the commands- set default listeners, add help subcmd, sync settings
	c._init()
	c._set_bin_name(os.Args[0])

	parser := NewParser(c)
	matches, err := parser.parse(os.Args[1:])

	if !err.is_nil {
		event := EventConfig{
			err:         err,
			args:        err.args,
			event:       err.kind,
			exit_code:   err.exit_code,
			app_ref:     c,
			matched_cmd: matches.matched_cmd,
		}
		c.emit(event)
	}

	// TODO: No errors, check special flags
	matched_cmd, cmd_idx := matches.GetMatchedCommand()

	// Check special flags
	// TODO: Sync with program settings
	if matches.ContainsFlag("help") {
		event := EventConfig{
			event:       OutputHelp,
			exit_code:   0,
			app_ref:     c,
			matched_cmd: matched_cmd,
		}
		c.emit(event)
	} else if matches.ContainsFlag("version") {
		// emit version event
		os.Exit(0)
	}

	if matched_cmd.callback != nil {
		// No args passed to the matched cmd
		if len(matches.raw_args[cmd_idx+1:]) == 0 && matched_cmd._isExpectingValues() {
			matched_cmd.PrintHelp()
			return
		} else {
			// Invoke callback
			matched_cmd.callback(&matches)
		}
	} else {
		matched_cmd.PrintHelp()
	}
}

/****************************** Event emitter functionality ****************************/

// Makes a call to the Command event emitter to `emit` a new event from the passed config
func (c *Command) emit(cfg EventConfig) {
	c.emitter.emit(cfg)
}

// Used to add a new listener for a specific event which gets triggered when the event occurs
func (c *Command) On(event Event, cb EventCallback) {
	c.emitter.on(event, cb, 0)
}

// This method is also used a new listener to a specific event but also overrides the default listener created by the package for said event
func (c *Command) Override(event Event, cb EventCallback) {
	c.emitter.override(event)
	c.emitter.on(event, cb, 0)
}

// A method for setting a listener that gets executed after all events encountered in the program
func (c *Command) AfterAll(cb EventCallback) {
	c.emitter.insert_after_all(cb)
}

// Set a callback to be executed only after the help event
func (c *Command) AfterHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, 4)
}

// Set a callback to be executed before all events encountered
func (c *Command) BeforeAll(cb EventCallback) {
	c.emitter.insert_before_all(cb)
}

// Set a callback to be executed only before the help event
func (c *Command) BeforeHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, -4)
}

/****************************** Other Command Utilities ****************************/

func (c *Command) find_subcommand(val string) (*Command, error) {
	for _, sc := range c.sub_commands {
		if sc.name == val || sc.alias == val {
			return sc, nil
		}
	}

	return NewCommand(""), errors.New("no such subcmd")
}

func (c *Command) _set_parent(parent *Command) *Command {
	c.parent = parent
	return c
}

func (c *Command) _set_is_root() *Command {
	c.is_root = true
	return c
}

func (c *Command) _get_usage_str() string {
	var new_usage strings.Builder

	if c.parent != nil && c.parent.is_root && !strings.Contains(c.usage_str, c.parent.usage_str) {
		new_usage.WriteString(c.parent.usage_str)
		new_usage.WriteString(c.usage_str)
	} else {
		new_usage.WriteString(c.usage_str)
	}

	return new_usage.String()
}

func (c *Command) _set_bin_name(val string) {
	if len(c.name) == 0 {
		bin_name := filepath.Base(val)

		// TODO: Validation
		c.Name(bin_name)
	}
}

func (c *Command) PrintHelp() {
	HelpWriter{}.Write(c)
}

/****************************** Interface Implementations ****************************/

func (c *Command) generate() (string, string) {
	// TODO: Check if allow command aliases
	return c.GetName(), c.GetHelp()
}
