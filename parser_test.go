package gommander

import (
	"testing"
)

func TestBasicParsing(t *testing.T) {
	cmd := NewCommand("test").Flag("-v --version", "Version flag").Option("-p --port <port-no>", "Port option")
	parser := NewParser(cmd)
	matches, _ := parser.parse([]string{"-v", "-p", "90"})

	if v, _ := matches.GetOptionValue("--port"); v != "90" {
		t.Error("Option arg parsing not working correctly")
	}

	if !matches.ContainsFlag("-v") {
		t.Error("Flag parsing has some errors")
	}

}

func TestStandardParsing(t *testing.T) {
	app := NewCommand("echo")
	app.SubCommand("first").Flag("-v --verbose", "Set verbose").Option("-n --name <value>", "Some name")

	parser := NewParser(app)
	matches, _ := parser.parse([]string{"first", "-v", "-n", "one", "-n", "two"})

	if v := matches.GetMatchedCommand(); v.name != "first" {
		t.Error("Subcommand resolution has some errors")
	}

	if !matches.ContainsFlag("verbose") {
		t.Error("Flag parsing has some errors")
	}

	values := matches.GetAllOptionInstances("name")
	if values[0] != "one" || values[1] != "two" {
		t.Error("Multiple option argument parsing failed")
	}

}

func TestComplexParsing(t *testing.T) {
	app := NewCommand("echo").Version("0.1.0").Help("A test CLI")

	app.SubCommand("image").
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

	if v := matches.GetMatchedCommand(); v.name != "image" {
		t.Error("Subcommand resolution has some errors")
	}

	if v, _ := matches.GetArgValue("<image-name>"); v != "image-one" {
		t.Errorf("Command argument parsing failed. Expected: image-one, got: %v", v)
	}

	if !matches.ContainsFlag("all") {
		t.Error("Flag parsing has some errors")
	}

	if v, _ := matches.GetOptionValue("--port"); v != "800" {
		t.Error("Option arg parsing not working correctly")
	}

	pstnl_args := matches.GetPositionalArgs()
	if pstnl_args[0] != "ng" || pstnl_args[1] != "serve" {
		t.Error("Positional args parsing failed")
	}

}

func TestOptionSyntaxParsing(t *testing.T) {
	app := NewCommand("basic").Option("-p --port <port-number>", "Port option")
	parser := NewParser(app)

	m_1, _ := parser.parse([]string{"-p", "9000"})
	m_2, _ := parser.parse([]string{"--port", "9000"})
	m_3, _ := parser.parse([]string{"--port=9000"})

	a_1, _ := m_1.GetOptionValue("port")
	a_2, _ := m_2.GetOptionValue("port")
	a_3, _ := m_3.GetOptionValue("port")

	if a_1 != a_2 {
		t.Error("Short option parsing and long option parsing out of sync")
	}

	if a_2 != a_3 {
		t.Error("Long option syntax with `=` parsing failed")
	}

}

func TestParserErrors(t *testing.T) {
	app := NewCommand("echo").Version("0.1.0").Help("A test CLI")

	app.SubCommand("image").
		Argument("<image-name>", "Provide an image name").
		Alias("i").
		Flag("--all", "Ran across all variants").
		Help("A first value subcommand").
		AddFlag(
			NewFlag("test").
				Short('t').
				Help("A simple test flag"),
		)

	parser := NewParser(app)
	_, err := parser.parse([]string{"i"})

	msg := "missing required argument: `<image-name>`"
	ctx := "Expected a required value corresponding to: `<image-name>` but none was provided"

	// Test missing required argument
	exp_err := throw_error(MissingRequiredArgument, msg, ctx)
	if !exp_err.compare(&err) {
		t.Error("Missing require argument error thrown incorrectly")
		t.Errorf("Expected error was: %v", exp_err.Error())
		t.Errorf("Found error was: %v", err.Error())
	}
}

func BenchmarkParseEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewParser(NewCommand("empty"))
		parser.parse([]string{})
	}
}
