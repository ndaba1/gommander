package gommander

import (
	"strings"
)

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
	}
}

// Value getters
func (app *App) Name() string { return app.name }

func (app *App) Alias() string { return app.alias }

func (app *App) Arguments() []*Argument { return app.arguments }

func (app *App) Flags() []*Flag { return app.flags }

func (app *App) Options() []*Option { return app.options }

func (app *App) SubCommands() []*App { return app.sub_commands }

func (app *App) Parent() *App { return app.parent }

func (app *App) Help() string { return app.help }

func (app *App) Version() string { return app.version }

// Value setters
func (app *App) SetName(name string) *App {
	app.name = name
	return app
}

func (app *App) SetAlias(alias string) *App {
	app.alias = alias
	return app
}

func (app *App) SetHelp(help string) *App {
	app.help = help
	return app
}

func (app *App) SetVersion(version string) *App {
	app.version = version
	return app
}

func (app *App) AddArgument(val string, help string) *App {
	argument := NewArgument(val, help)
	app.arguments = append(app.arguments, argument)
	return app
}

func (app *App) AddFlag(val string) *App {
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

func (app *App) AddOption(val string) *App {
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

func (app *App) AddSubCommand(sub_command *App) *App {
	app.sub_commands = append(app.sub_commands, sub_command)
	return app
}

func (app *App) SetParent(parent *App) { app.parent = parent }

type Flag struct {
	name  string
	long  string
	short string
	help  string
}

type Argument struct {
	name        string
	help        string
	raw         string
	variadic    bool
	is_required bool
}

func NewArgument(val string, help string) *Argument {
	var delimiters []string
	var required bool
	var variadic bool

	// FIXME: Find more robust way for checking
	if strings.ContainsAny(val, "<") {
		delimiters = []string{"<", ">"}
		required = true
	} else if strings.ContainsAny(val, "[") {
		delimiters = []string{"[", "]"}
		required = false
	}

	name := strings.Replace(val, delimiters[0], "", -1)
	name = strings.Replace(name, delimiters[1], "", -1)
	name = strings.Replace(name, "-", "_", -1)

	if strings.HasSuffix(val, "...") {
		variadic = true
		name = strings.Replace(name, "...", "", -1)
	}

	return &Argument{
		name:        name,
		help:        help,
		raw:         val,
		variadic:    variadic,
		is_required: required,
	}
}

type Option struct {
	name  string
	help  string
	short string
	long  string
	args  []*Argument
}
