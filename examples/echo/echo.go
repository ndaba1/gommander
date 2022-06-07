package main

import (
	"fmt"

	"github.com/ndaba1/gommander"
)

func main() {
	app := gommander.App()

	app.Name("echo").
		Version("0.1.0").
		Help("A simple echo clone").
		Author("vndaba")

	// Creating arguments for the program
	app.AddArgument(
		gommander.
			NewArgument("text").
			Help("The text to echo out").
			Default("This is a default value that gets printed out if no value is passed").
			Required(true).
			Variadic(true),
	)

	// An alternate, less verbose option syntax. Duplicate args will be ignored.
	app.Argument("<text...>", "The text to echo out")

	// Creating flags
	app.AddFlag(
		gommander.
			NewFlag("newline").
			Short('n').
			Help("Whether or not to add a newline"),
	)

	// An alternate way for creating flags. Duplicate values also ignored
	app.Flag("-n --newline", "Whether or not to add a newline")

	// Declaring a callback for the command
	app.Action(echoCb)

	app.Parse()
}

func echoCb(pm *gommander.ParserMatches) {
	// returns the value of an arg or an error if no such value is found
	text, err := pm.GetArgValue("text")
	if err != nil {
		fmt.Println(err.Error())
	}

	if pm.ContainsFlag("newline") {
		fmt.Println(text)
	} else {
		fmt.Print(text)
	}
}
