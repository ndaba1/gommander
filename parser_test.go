package gommander

import (
	"testing"
)

func TestBasicParsing(t *testing.T) {
	cmd := NewCommand("test").Flag("-v --version", "Version flag").Option("-p --port <port-no>", "Port option")
	parser := NewParser(cmd)
	matches, _ := parser.parse([]string{"-v", "-p", "90"})

	if v, _ := matches.GetOptionArg("--port"); v != "90" {
		t.Error("Option arg parsing not working correctly")
	}

	if !matches.ContainsFlag("-v") {
		t.Error("Flag parsing has some errors")
	}

}

func TestStandardParsing(t *testing.T) {
	app := NewCommand("echo")
	app.Subcommand("first").Flag("-v --verbose", "Set verbose").Option("-n --name <value>", "Some name")

	parser := NewParser(app)
	matches, _ := parser.parse([]string{"first", "-v", "-n", "one", "-n", "two"})

	if matches.GetMatchedCommand().name != "first" {
		t.Error("Subcommand resolution has some errors")
	}

	if !matches.ContainsFlag("verbose") {
		t.Error("Flag parsing has some errors")
	}

	values := matches.GetAllOptionArgs("name")
	if values[0] != "one" || values[1] != "two" {
		t.Error("Multiple option argument parsing failed")
	}

}

func TestComplexParsing(t *testing.T) {
	app := NewCommand("echo").Version("0.1.0").Help("A test CLI")

	app.Subcommand("image").
		Argument("<image-name>", "Provide an image name").
		Alias("i").
		Flag("--all", "Ran across all variants").
		Help("A first value subcommand").
		AddFlag(
			NewFlag("test").
				Short('t').
				Help("A simple test flag"),
		).
		AddOption(
			NewOption("port").
				Short('p').
				Required(true).
				Help("Pass the port number").
				AddArgument(
					NewArgument("port-number").
						Required(true),
				),
		)

	parser := NewParser(app)
	matches, _ := parser.parse([]string{"i", "image-one", "--all", "-p", "800", "--", "ng", "serve"})

	if matches.GetMatchedCommand().name != "image" {
		t.Error("Subcommand resolution has some errors")
	}

	if v, _ := matches.GetArgValue("<image-name>"); v != "image-one" {
		t.Errorf("Command argument parsing failed. Expected: image-one, got: %v", v)
	}

	if !matches.ContainsFlag("all") {
		t.Error("Flag parsing has some errors")
	}

	if v, _ := matches.GetOptionArg("--port"); v != "800" {
		t.Error("Option arg parsing not working correctly")
	}

	pstnl_args := matches.GetPositionalArgs()
	if pstnl_args[0] != "ng" || pstnl_args[1] != "serve" {
		t.Error("Positional args parsing failed")
	}

}
