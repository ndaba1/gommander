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
		).Action(greetCb)

	app.Set(gommander.IncludeHelpSubcommand, true)

	app.Parse()
}

func greetCb(matches *gommander.ParserMatches) {
	var completeGreeting strings.Builder

	greeting, _ := matches.GetOptionValue("lang")
	switch greeting {
	case "ENGLISH":
		completeGreeting.WriteString("Hello! ")
	case "SPANISH":
		completeGreeting.WriteString("Hola! ")
	case "SWAHILI":
		completeGreeting.WriteString("Jambo! ")
	}

	// Returns the value and error if any
	name, _ := matches.GetOptionValue("name")
	completeGreeting.WriteString(name)

	fmt.Println(completeGreeting.String())
}
