// Package gommander is an easily-extensible commander package for easily creating Command Line Interfaces.
// View the README for getting started documentation:
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
	aliases        []string
	arguments      []*Argument
	author         string
	callback       CommandCallback
	discussion     string
	emitter        EventEmitter
	flags          []*Flag
	help           string
	isRoot         bool
	name           string
	options        []*Option
	parent         *Command
	subCommands    []*Command
	settings       AppSettings
	globalSettings *AppSettings
	theme          Theme
	version        string
	usageStr       string
	customUsageStr string
	subCmdGroups   map[string][]*Command
	appRef         *Command
}

func App() *Command {
	// return NewCommand("")._setRoot().AddFlag(versionFlag()).Theme(DefaultTheme())
	return &Command{
		isRoot:       true,
		flags:        []*Flag{helpFlag(), versionFlag()},
		parent:       nil,
		settings:     AppSettings{},
		emitter:      newEmitter(),
		subCmdGroups: make(map[string][]*Command),
		theme:        DefaultTheme(),
	}
}

func NewCommand(name string) *Command {
	return &Command{
		name:           name,
		flags:          []*Flag{helpFlag()},
		settings:       AppSettings{},
		isRoot:         false,
		emitter:        newEmitter(),
		globalSettings: &AppSettings{},
		usageStr:       name,
		subCmdGroups:   make(map[string][]*Command),
	}
}

/****************************** Value Getters ****************************/

// Simply returns the alias of the command or an empty string
func (c *Command) GetAliases() []string { return c.aliases }

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
func (c *Command) GetSubCommands() []*Command { return c.subCommands }

// Returns the version of the program
func (c *Command) GetVersion() string { return c.version }

// Returns the default usage string or a custom_usage_str if one exists
func (c *Command) GetUsageStr() string {
	if len(c.customUsageStr) > 0 {
		return c.customUsageStr
	}
	return c.usageStr
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
	c.aliases = append(c.aliases, alias)
	return c
}

func (c *Command) AddArgument(arg *Argument) *Command {
	for _, a := range c.arguments {
		if a.name == arg.name {
			return c
		}
	}
	c.arguments = append(c.arguments, arg)
	return c
}

// A method for setting any expected arguments for a command, it takes in the value of the argument e.g. `<image-name>` and the help string for said argument
func (c *Command) Argument(val string, help string) *Command {
	argument := newArgument(val, help)
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
	flag := newFlag(val, help)
	return c.AddFlag(&flag)
}

// Used to set more information or the command discussion which gets printed out when help is invoked, at the bottom most section
func (c *Command) Discussion(info string) *Command {
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
	c.usageStr = name
	return c
}

// Sets the version of a command, usually the entry point command(App)
func (c *Command) Version(version string) *Command {
	c.version = version
	return c
}

// An identical method to the `.Flag()` method but for options. Expected syntax: "-p --port <port-number>"
func (c *Command) Option(val string, help string) *Command {
	option := newOption(val, help, false)
	return c.AddOption(&option)
}

// This method is used to mark an option as required for a given command. Another way of achieving this is using the `.AddOption()` method and using the `NewOption()` builder interface to define option parameters
func (c *Command) RequiredOption(val string, help string) *Command {
	opt := newOption(val, help, true)
	return c.AddOption(&opt)
}

// Used to define a custom usage string. If one is present, it will be used instead of the default one
func (c *Command) UsageStr(val string) *Command {
	c.customUsageStr = val
	return c
}

/****************************** Subcommand related methods ****************************/

// When chained on a command, this method adds said command to the provided sub_cmd group in the parent of the command.
func (c *Command) AddToGroup(name string) *Command {
	c.parent.subCmdGroups[name] = append(c.parent.subCmdGroups[name], c)
	return c
}

// Receives a reference to a command, sets the command parent and usage string then adds its to the slice of subcommands. This method is called internally by the `.SubCommand()` method but users can also invoke it directly
func (c *Command) AddSubCommand(subCmd *Command) *Command {
	subCmd.parent = c
	c.subCommands = append(c.subCommands, subCmd)

	cmdPath := []string{c.usageStr, subCmd.usageStr}
	subCmd.usageStr = strings.Join(cmdPath, " ")

	// Propagate global flags to children
	for _, f := range c.GetFlags() {
		if f.isGlobal {
			subCmd.AddFlag(f)
		}
	}

	// propagate theme
	subCmd.theme = c.theme

	if c.isRoot {
		subCmd.appRef = c
	} else {
		subCmd.appRef = c.appRef
	}

	return c
}

// An easier method for creating sub_cmds while avoiding too much function paramets nesting. It accepts the name of the new sub_cmd and returns the newly created sub_cmd
func (c *Command) SubCommand(name string) *Command {
	subCmd := NewCommand(name)
	c.AddSubCommand(subCmd)
	return subCmd
}

// A manual way of creating a new subcommand group and adding the desired commands to it
func (c *Command) SubCommandGroup(name string, vals []*Command) {
	c.subCmdGroups[name] = append(c.subCmdGroups[name], vals...)
}

/****************************** Settings ****************************/

func (c *Command) _init() {
	if c.settings[DisableVersionFlag] {
		c.removeFlag("--version")
	}

	if c.settings[IncludeHelpSubcommand] && len(c.subCommands) > 0 {
		validSubcmds := []string{}

		for _, c := range c.subCommands {
			validSubcmds = append(validSubcmds, c.name)
		}

		c.SubCommand("help").
			Help("Print out help information for the passed command").
			AddArgument(
				NewArgument("<COMMAND>").
					Help("The name of the command to output help for").
					ValidateWith(validSubcmds),
			).
			Action(func(pm *ParserMatches) {
				val, _ := pm.GetArgValue("<COMMAND>")
				parent := pm.matchedCmd.parent

				if parent != nil {
					cmd, _ := parent.findSubcommand(val)
					cmd.PrintHelp()
				}
			})
	}

	// Default help listener cannot be overridden
	c.emitter.on(OutputHelp, func(ec *EventConfig) {
		cmd := ec.matchedCmd
		cmd.PrintHelp()
	}, -4)

	if !c.settings[OverrideAllDefaultListeners] {
		c.emitter.onErrors(func(ec *EventConfig) {
			err := ec.err
			// TODO: Match theme in better way
			err.Display(c)
		})

		c.emitter.on(OutputVersion, func(ec *EventConfig) {
			// TODO: Print version in a better way
			app := ec.appRef

			fmt.Println(app.GetName(), app.GetVersion())
			fmt.Println(app.GetAuthor())
			fmt.Println(app.GetHelp())
		}, -4)

		for _, event := range c.emitter.eventsToOverride {
			c.emitter.rmDefaultLstnr(event)
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
	hasDefaults := func(list []*Argument) bool {
		for _, a := range c.arguments {
			if a.hasDefaultValue() {
			} else {
				return false
			}
		}
		return true
	}

	return len(c.subCommands) > 0 || (len(c.arguments) > 0 && !hasDefaults(c.arguments))
}

func (c *Command) _parse(vals []string) {
	// TODO: Init/build the commands- set default listeners, add help subcmd, sync settings
	c._init()
	c._setBinName(vals[0])

	rawArgs := vals[1:]
	parser := NewParser(c)
	matches, err := parser.parse(rawArgs)

	if err != nil {
		event := EventConfig{
			err:        *err,
			args:       err.args,
			event:      err.kind,
			exitCode:   err.exitCode,
			appRef:     c,
			matchedCmd: matches.matchedCmd,
		}
		c.emit(event)
	}

	// TODO: No errors, check special flags
	matchedCmd := matches.GetMatchedCommand()
	cmdIdx := matches.GetMatchedCommandIndex()

	// Check special flags
	// TODO: Sync with program settings
	if matches.ContainsFlag("help") {
		event := EventConfig{
			event:      OutputHelp,
			exitCode:   0,
			appRef:     c,
			matchedCmd: matchedCmd,
		}
		c.emit(event)
	} else if matches.ContainsFlag("version") {
		event := EventConfig{
			event:      OutputVersion,
			exitCode:   0,
			appRef:     c,
			matchedCmd: matchedCmd,
		}
		c.emit(event)
	}

	showHelp := func() {
		if !isTestMode() {
			matchedCmd.PrintHelp()
		}
	}

	if matchedCmd.callback != nil {
		// No args passed to the matched cmd
		if cmdIdx == -1 {
			cmdIdx++
		}
		if (len(rawArgs) == 0 || len(matches.rawArgs[cmdIdx:]) == 0) && matchedCmd._isExpectingValues() {
			showHelp()
			return
		}
		// Invoke callback
		matchedCmd.callback(matches)
	} else {
		showHelp()
	}
}

// A method for parsing the arguments passed to a program and invoking the callback on a command if one is found. This method also handles any errors encountered while parsing.
func (c *Command) Parse() {
	c._parse(os.Args)
}

func (c *Command) ParseFrom(args []string) {
	c._parse(args)
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

// This method is also used to add a new listener to a specific event but also overrides the default listener created by the package for said event
func (c *Command) Override(event Event, cb EventCallback) {
	c.emitter.override(event)
	c.emitter.on(event, cb, 0)
}

// A method for setting a listener that gets executed after all events encountered in the program
func (c *Command) AfterAll(cb EventCallback) {
	c.emitter.insertAfterAll(cb)
}

// Set a callback to be executed only after the help event
func (c *Command) AfterHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, 4)
}

// Set a callback to be executed before all events encountered
func (c *Command) BeforeAll(cb EventCallback) {
	c.emitter.insertBeforeAll(cb)
}

// Set a callback to be executed only before the help event
func (c *Command) BeforeHelp(cb EventCallback) {
	c.emitter.on(OutputHelp, cb, -4)
}

/****************************** Other Command Utilities ****************************/

func (c *Command) findSubcommand(val string) (*Command, error) {
	for _, sc := range c.subCommands {
		includes := func(val string) bool {
			for _, v := range sc.aliases {
				if v == val {
					return true
				}
			}
			return false
		}

		if sc.name == val || includes(val) {
			return sc, nil
		}
	}

	return NewCommand(""), errors.New("no such subcmd")
}

func (c *Command) removeFlag(val string) {
	newFlags := []*Flag{}
	for _, f := range c.flags {
		if !(f.short == val || f.long == val) {
			newFlags = append(newFlags, f)
		}
	}
	c.flags = newFlags
}

func (c *Command) _getAppRef() *Command {
	if c.isRoot {
		return c
	}
	return c.appRef
}

func (c *Command) _getUsageStr() string {
	var newUsage strings.Builder

	if len(c.customUsageStr) > 0 {
		if !strings.Contains(c.customUsageStr, c.parent.usageStr) {
			newUsage.WriteString(c.parent.usageStr)
			newUsage.WriteRune(' ')
		}
		newUsage.WriteString(c.customUsageStr)
	} else {
		if c.parent != nil && c.parent.isRoot && !strings.Contains(c.usageStr, c.parent.usageStr) {
			newUsage.WriteString(c.parent.usageStr)
		}
		newUsage.WriteString(c.usageStr)
	}

	return newUsage.String()
}

func (c *Command) _setBinName(val string) {
	if len(c.name) == 0 {
		binName := filepath.Base(val)

		// TODO: Validation
		c.Name(binName)
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
