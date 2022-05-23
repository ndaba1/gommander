package gommander

import (
	"fmt"
	"strings"
)

type Callback = func(ParserMatches)

type App struct {
	name         string
	alias        string
	arguments    []*Argument
	flags        []*Flag
	options      []*Option
	sub_commands []*App
	parent       *App
	help         string
	version      string
	settings     Settings
	theme        Theme
	discussion   string
	is_root      bool
	callback     Callback
}

func Program() *App {
	return Command("").set_is_root(true)
}

func Command(name string) *App {
	return &App{
		name:         name,
		alias:        "",
		arguments:    []*Argument{},
		flags:        []*Flag{},
		options:      []*Option{},
		sub_commands: []*App{},
		parent:       nil,
		help:         "",
		version:      "",
		settings:     Settings{},
		theme:        Theme{},
		discussion:   "",
		is_root:      false,
	}
}

// Value getters
func (app *App) GetName() string { return app.name }

func (app *App) GetAlias() string { return app.alias }

func (app *App) GetArguments() []*Argument { return app.arguments }

func (app *App) GetFlags() []*Flag { return app.flags }

func (app *App) GetOptions() []*Option { return app.options }

func (app *App) GetSubCommands() []*App { return app.sub_commands }

func (app *App) GetParent() *App { return app.parent }

func (app *App) GetHelp() string { return app.help }

func (app *App) GetVersion() string { return app.version }

// Value setters
func (app *App) Name(name string) *App {
	app.name = name
	return app
}

func (app *App) Alias(alias string) *App {
	app.alias = alias
	return app
}

func (app *App) Help(help string) *App {
	app.help = help
	return app
}

func (app *App) Version(version string) *App {
	app.version = version
	return app
}

func (app *App) Action(cb Callback) *App {
	app.callback = cb
	return app
}

func (app *App) Argument(val string, help string) *App {
	argument := NewArgument(val, help)
	app.arguments = append(app.arguments, argument)
	return app
}

func (app *App) Flag(val string) *App {
	values := strings.Split(val, ",")
	flag := Flag{
		short: values[0],
		long:  values[1],
		help:  values[2],
		name:  strings.Replace(values[1], "-", "", -1),
	}
	app.flags = append(app.flags, &flag)
	return app
}

func (app *App) Option(val string) *App {
	values := strings.Split(val, ",")

	var arg_slice []*Argument
	for idx, a := range values {
		if idx > 1 && idx < (len(values)-1) {
			arg := NewArgument(a, "")
			arg_slice = append(arg_slice, arg)
		}
	}

	option := Option{
		short: values[0],
		long:  values[1],
		help:  values[2],
		name:  strings.Replace(values[1], "-", "", -1),
		args:  arg_slice,
	}

	app.options = append(app.options, &option)
	return app
}

func (app *App) Subcommand(name string) *App {
	return Command(name).set_parent(app)
}

// Interior utility functions
func (app *App) set_parent(parent *App) *App {
	app.parent = parent
	return app
}

func (app *App) set_is_root(val bool) *App {
	app.is_root = val
	return app
}

func (app *App) PrintHelp() {
	fmt.Printf(app.help)

	fmt.Printf("\n USAGE: \n")
	fmt.Printf("\t .exe [OPTIONS] [COMMAND]")
}
