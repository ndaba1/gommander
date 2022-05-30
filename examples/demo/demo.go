package main

import (
	"fmt"
	"strings"

	"github.com/ndaba1/gommander"
)

func main() {
	app := gommander.App()

	app.Author("vndaba").
		Version("0.1.0").
		Help("A demo usage of gommander")

	app.SubCommand("greet").
		Help("A simple cmd that greets provided name").
		AddOption(
			gommander.
				NewOption("name").
				Required(true).
				Short('n').
				Help("The name to greet").
				AddArgument(
					gommander.NewArgument("<name>"),
				),
		).
		AddOption(
			gommander.
				NewOption("lang").
				Short('l').
				Help("The language to use").
				Required(true).
				AddArgument(
					gommander.
						NewArgument("<language>").
						ValidateWith([]string{"ENGLISH", "SPANISH", "SWAHILI"}).
						Default("ENGLISH"),
				),
		).Action(greet_cb)

	app.Set(gommander.IncludeHelpSubcommand, true)

	app.Parse()
}

func greet_cb(matches *gommander.ParserMatches) {
	var complete_greeting strings.Builder

	greeting, _ := matches.GetOptionValue("lang")
	switch greeting {
	case "ENGLISH":
		complete_greeting.WriteString("Hello! ")
	case "SPANISH":
		complete_greeting.WriteString("Hola! ")
	case "SWAHILI":
		complete_greeting.WriteString("Jambo! ")
	}

	// Returns the value and error if any
	name, _ := matches.GetOptionValue("name")
	complete_greeting.WriteString(name)

	fmt.Println(complete_greeting.String())
}
