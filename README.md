# Gommander (go-commander)

<p align="center">
<img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/ndaba1/gommander/gommander-ci-workflow?label=CI&logo=github%20actions&logoColor=fff">
<img alt="Go report card", src="https://goreportcard.com/badge/github.com/ndaba1/gommander">
<img alt="Go reference", src="https://pkg.go.dev/badge/github.com/ndaba1/gommander.svg">
</p>

A commander package for creating CLIs in Go. This package aims to be a complete solution for command line argument parsing by providing you with an easy-to-use and extensible api, but without compromising speed.

Features of this package include:

- Support for POSIX-compliant flags
- Color printing
- Easily extensible API
- Fast and lightweight parser
- Subcommands nesting
- Automatic help generation

## Index

- [Gommander](#gommander-go-commander)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
  - [Subcommands](#subcommands)
  - [Arguments](#arguments)
  - [Flags](#flags)
  - [Options](#options)
  - [App Settings and Events](#settings-and-events)
  - [App Themes and UI](#themes-and-ui)
  - [Command Callbacks](#command-callbacks)

## Installation

To add the package to your go project, execute the following command:

```bash
go get github.com/ndaba1/gommander
```

## Quick Start

Gettting started with gommander is quite easy. Create an instance of `gommander.App()` which is the internal representation of your CLI program to which you add [subcommands](#subcommands), [arguments](#arguments), [flags](#flags) and [options](#options). The following is a quick example:

```go
package main

import (
  "fmt"

  "github.com/ndaba1/gommander"
)

func main() {
  app := gommander.App();

  app.Author("vndaba").
      Version("0.1.0").
      Help("A demo usage of gommander").
      Name("demo")

  app.Subcommand("greet").
      Argument("<name>", "The name to greet").
      Help("A simple command that greets the provided name").
      Option("-c --count <number>", "The number of times to greet the name")

  app.Parse()
}

```

When the help flag is invoked on the above example, the following output gets printed out:

```bash

A demo usage of gommander

Usage:
    demo [FLAGS] <SUBCOMMAND>

Flags:
    -h, --help        Print out help information
    -v, --version     Print out version information

Subcommands:
    greet             A simple cmd that greets provided name

```

## Subcommands

Commands and subcommands are the gist of the package. The entrypoint of the program is itself a command to which more subcommands may be nested. A new command can be created by the `gommander.NewCommand()` method but you will rarely have to work with this method directly. The convention is to first create an app via the `gommander.App()` method which itself is also a command that has been marked as the entrypoint of the program(the root cmd) then chain further subcommands to it.

```go
//...package declaration and imports
func main() {
    app := gommander.App()

    app.SubCommand("basic")
}
//...
```

The `.SubCommand()` method returns the newly created subcommand for further manipulation and customization such as adding [arguments](#arguments), [flags](#flags), [options](#options) or updating command metadata such as it description, etc.

You can also manipulate the fields of the root command, such the author of the program, the version and even set the name of the program to be different from the name of its binary.(This doesn't change the name of the actual binary, only changes the name display to users when printing out help information).

```go
// ...
func main() {
    app := gommander.App()

    app.Help("A simple CLI app").
        Version("0.1.0").
        Author("vndaba")
}
// ...
```

Subcommand Nesting is also supported by the package. You can chain as many subcommands as you would like. The following is an example:

```go
// ...
func main() {
    app := gommander.App()

    image := app.SubCommand("image").Help("Manage images")

    image.SubCommand("ls").Help("List available images")
    image.SubCommand("build").Help("Build a new image")

    // ...
}
// ...
```

The `.SubCommand()` method is very ergonomic to use and reduces function parameters nesting in your program. However, you could also use the `.AddSubCommand()` method to add a new subcommand. It is similar to the subcommand method, but instead of taking the name of the new command as input, it takes it an instance of an already created command and instead of returning the newly created subcommand, it returns the command to which it is chained as so:

```go
// ...
func main() {
    app := gommander.App()

    app.AddSubCommand(
		gommander.NewCommand("first").Help("A basic subcommand"),
	).AddSubCommand(
		gommander.NewCommand("second").Help("Another basic subcommand"),
	)
}
// ...
```

You could keep chaining subcommands using this method as shown. The `.SubCommand()` method invokes this one internally.

Throughout the package, you will notice a similar pattern to other methods as well. For instance, flags can be created via the `.Flag()` method but there also exists a `.AddFlag()` method which follows the same rules as the subcommand method counterpart. The same also applies for [options](#options) and [arguments](#arguments).

### Subcommand Groups

Subcommands can also be added to groups. This may be done to change how they are printed when showing help. For this, the `.AddToGroup()` method may be used. An example of this is shown in the [subcommands](./examples/subcommands/subcommands.go) example.

## Arguments

Arguments are values passed to a command as input. They can also be passed to an option. For instance, in the demo example in the [quick start section](#quick-start), the greet subcommand takes in a name as input. The program is therefore expected to be run as follows:

```bash
./demo.exe greet John
```

Here, `John` is a value that is an instance of the `<name>` argument. Adding arguments to commands is fairly simple. You could either use the [.Argument()](#argument-method) method or the [.AddArgument()](#addargument-method) one.

### **Argument() method**

```go
// ...
func main() {
    app := gommander.App()

    app.Argument("<basic>", "Some basic argument")
}
// ...
```

Here, we see the value of the argument enclosed between angle brackets. **This means that the argument is required** and therefore, an error will be thrown if one is not passed to the program. Optional arguments are represented by square brackets: `[arg]`. If neither the square or angle brackets are provided, the argument is marked as optional.
The `.Argument()` method takes in the value of the argument and its help string/description. Here are the acceptable forms for the argument value:
| Value | Semantics|
|:-----|:--------|
| `<arg>` | Argument is required |
| `<arg...>` | Argument is required and is variadic|
| `[arg]` | Argument is optional|
| `[arg...]` | Argument is optional but variadic if provided|

### **AddArgument() method**

```go
// ...
func main() {
    app := gommander.App()

    app.AddArgument(
		gommander.NewArgument("language").
			Required(true).
			Variadic(false).
			ValidateWith([]string{"ENGLISH", "SPANISH", "FRENCH", "SWAHILI"}).
			Default("ENGLISH").
			Help("The language to use"),
	)
}
// ...
```

The `.AddArgument()` method, while more verbose, provides more flexibility in defining arguments. It ought to be used when defining more complex arguments. The `gommander.NewArgument()` returns an instance of an Argument to which more methods can be chained. Most of the methods are axiomatic and their functionality can be deduced from their names.
The `.ValidateWith()` method sets valid_values for an argument. If the value passed is not one of those values, a well-described error is thrown by the program and printed out.
The `.Default()` method sets a default value for an argument. If the argument is required but no value was passed, the default value is used.

Arguments can also be passed to options. This is discussed in depth in the [options](#options) section

## Flags

Flags are values, prefixed with either a `-` in their short form, or `--` in their long form. Adding a flag to an instance of a command is simple. It can be achieved in one of two ways:

- The `.Flag()` method:

```go
// ...
func main() {
    app := gommander.App()

    app.Flag("-V --verbose", "A flag for verbosity")
}
// ...
```

- The `.AddFlag()` method:

```go
// ...
func main() {
    app := gommander.App()

    app.AddFlag(
        gommander.NewFlag("verbose").
            Short('V').
            Help("A flag for verbosity").
            Global(true),
    )
}
// ...
```

When a flag is set as global, it will be propagated to all the subcommands in the app.
The parser also supports posix flag syntax, therefore, if a command contains flags, say `-i`, `-t`, `-d`, instead of passing the flags individually to the program, users can combine the flags as `itd`.

## Options

Options are simply flags that take in a value as input. There are also two ways for declaring options:

- The `.Option()` method

```go
// ...
func main() {
    app := gommander.App()

    app.Option("-p --port <port-number>", "The port number to use")
}
// ...
```

- The `.AddOption()` method:

```go
// ...
func main() {
    app := gommander.App()

    app.AddOption(
        gommander.NewOption("port").
            Short('p').
            AddArgument(
                gommander.NewArgument("<port-number>").
                    Default("9000"),
            ),
    )
}
// ...
```

When adding an argument, you can either use the `.Argument()` method or the `.AddArgument()` one.
Support option syntaxes are:

- `-p 80`
- `--port 80`
- `--port=80`

## Settings and Events

The default behavior of the program can be easily modified or even overriden. This can be achieved through settings and events.
The program settings are very straight-forward and can be accessed via the `Command.Set()` method. Settings configured via this method are simple boolean values and include:

```go
// ...
func main() {
    app := gommander.App()

    // .... app logic

    app.Set(gommander.IncludeHelpSubcommand, true).
		Set(gommander.DisableVersionFlag, true).
		Set(gommander.IgnoreAllErrors, false).
		Set(gommander.ShowHelpOnAllErrors, true).
		Set(gommander.ShowCommandAliases, true).
		Set(gommander.OverrideAllDefaultListeners, false)

    app.Parse()
}
// ...
```

The package also has the concept of events which are emitted by the app and can be reacted to by adding new listeners or even overriding the default listeners.

```go
// ...
func main() {
    app := gommander.App()

    // .... app logic

    app.On(gommander.OutputVersion, func(ec *gommander.EventConfig) {
		app := ec.GetApp()

		fmt.Printf("You are version: %v of %v which was authored by: %v", app.GetVersion(), app.GetName(), app.GetAuthor())
	})
}
// ...
```

The `.On()` method is used to register a new callback for a given event. This callback is defined as: `type EventCallback = func(*gommander.EventConfig)`. The `.On()` method adds the new callback after the default ones. If you wish to override the default behavior, use the `.Override()` method to register the callback.

Other event-related methods include:

- `Command.BeforeAll()`
- `Command.BeforeHelp()`
- `Command.AfterAll()`
- `Command.AfterHelp()`

## Themes and UI

Themes control the color palette used by the program. You can define your own theme or use predefined ones. The package uses `github.com/fatih/color` as a dependency for color functionality.
The package uses the concept of designations for theme functionality, i.e. `Keywords` are assigned one color, `Descriptions` another, `Errors` another and so on and so forth.

Using a predefined theme is as shown below:

```go
// ...
func main() {
    app := gommander.App()

    // .... app logic

    app.UsePredefinedTheme(gommander.PlainTheme)
}
// ...
```

You can also create your own theme. You will need to import the `github.com/fatih/color` package

```go
app.Theme(
		gommander.NewTheme(color.FgGreen, color.FgBlue, color.FgRed, color.FgWhite, color.FgWhite),
	)
```

The NewTheme method takes in values of type `ColorAttribute` defined in the `fatih/color` package.

## Command Callbacks

The package only serves one purpose, to parse command line arguments. To define what to do with the parsed arguments, command callbacks are defined. There are simply functions of the type: `func(*gommander.ParserMatches)` that get invoked when a command is matched.
These functions are defined by the `Command.Action()` method.

See an example of this [here](./examples/demo/demo.go).
