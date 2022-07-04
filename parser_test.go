package gommander

import (
	"testing"
)

func TestParseBasic(t *testing.T) {
	clearCache()
	cmd := NewCommand("test").Flag("-v --version", "Version flag").Option("-p --port <port-no>", "Port option")
	parser := NewParser(cmd)
	matches, _ := parser.parse([]string{"-v", "-p", "90"})

	v, _ := matches.GetOptionValue("--port")
	assertEq(t, v, "90", "Option arg parsing not working correctly")
	assert(t, matches.ContainsFlag("-v"), "Flag parsing has some errors")
}

func TestParseStandard(t *testing.T) {
	clearCache()
	app := NewCommand("echo")
	app.SubCommand("first").Flag("-v --verbose", "Set verbose").Option("-n --name <value>", "Some name")

	parser := NewParser(app)
	matches, _ := parser.parse([]string{"first", "-v", "-n", "one", "-n", "two"})

	assert(t, matches.ContainsFlag("verbose"), "Flag parsing has some errors")
	assertEq(t, matches.GetMatchedCommand().name, "first", "Subcommand resolution has some errors")
	assertDeepEq(t, matches.GetAllOptionInstances("name"), []string{"one", "two"}, "Multiple option argument parsing failed")
}

func TestParseComplex(t *testing.T) {
	clearCache() // workaound for test mode

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

	assert(t, matches.ContainsFlag("all"), "Flag parsing has some errors")
	assertEq(t, matches.GetMatchedCommand().name, "image", "Subcommand resolution has some errors")
	assertDeepEq(t, matches.GetPositionalArgs(), []string{"ng", "serve"}, "Positional args parsing failed")

	v, _ := matches.GetArgValue("<image-name>")
	v2, _ := matches.GetOptionValue("--port")

	assertEq(t, v, "image-one", "Command argument parsing")
	assertEq(t, v2, "800", "Option arg parsing not working correctly")
}

func TestParseOptionSyntaxes(t *testing.T) {
	clearCache()
	app := NewCommand("basic").Option("-p --port <port-number>", "Port option")
	parser := NewParser(app)

	m1, _ := parser.parse([]string{"-p", "9000"})
	m2, _ := parser.parse([]string{"--port", "9000"})
	m3, _ := parser.parse([]string{"--port=9000"})

	a1, _ := m1.GetOptionValue("port")
	a2, _ := m2.GetOptionValue("port")
	a3, _ := m3.GetOptionValue("port")

	assertEq(t, a1, a2, "Short option parsing and long option parsing out of sync")
	assertEq(t, a2, a3, "Long option syntax with `=` parsing failed")
}

func _assertParserError(t *testing.T, app *Command, parserArgs, errorArgs []string, event Event, msg string) {
	clearCache()
	parser := NewParser(app)

	_, err := parser.parse(parserArgs)
	expected := generateError(app, event, errorArgs)

	assertDeepEq(t, *err, expected, msg)
}

func TestParserErrors(t *testing.T) {
	clearCache()
	app := NewCommand("echo").Version("0.1.0").Help("A test CLI")

	app.SubCommand("image").
		Argument("<image-name>", "Provide an image name").
		Alias("i").
		Flag("--all", "Ran across all variants").
		Help("A first value subcommand").
		AddOption(
			NewOption("port").
				Short('p').
				Argument("<int:port-no>").
				Help("option argument").
				Required(true),
		).
		AddFlag(
			NewFlag("test").
				Short('t').
				Help("A simple test flag"),
		)

		// test missing required argument
	_assertParserError(t, app,
		[]string{"i"},
		[]string{"<image-name>"},
		MissingRequiredArgument,
		"Missing required arg error detection failed",
	)

	_assertParserError(t, app,
		[]string{"i", "imageOne"},
		[]string{"--port"},
		MissingRequiredOption,
		"Missing required option error detection failed",
	)

	_assertParserError(t, app,
		[]string{"invalid"},
		[]string{"invalid"},
		UnknownCommand,
		"Unknown command error detection failed",
	)

	_assertParserError(t, app,
		[]string{"i", "imageOne", "--port=90", "-x"},
		[]string{"-x"},
		UnknownOption,
		"Unknown option error detection failed",
	)

	_assertParserError(t, app,
		[]string{"i", "imageOne", "--port=90", "invalid"},
		[]string{"invalid"},
		UnresolvedArgument,
		"Unresolved argument error detection failed",
	)

	_assertParserError(t, app,
		[]string{"i", "imageOne", "--port=hello"},
		[]string{"hello", "`hello` is not a valid integer"},
		InvalidArgumentValue,
		"invalid argument error detection failed",
	)

}

func BenchmarkParseEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser := NewParser(NewCommand("empty"))
		parser.parse([]string{})
	}
}
