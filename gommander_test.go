package gommander

import (
	"testing"
)

func TestCommandMetadata(t *testing.T) {
	cmd := NewCommand("basic").
		Help("Basic command").
		Alias("b").
		Author("vndaba").
		Version("0.1.0")

	assertEq(t, cmd.GetHelp(), "Basic command", "Cmd help field set incorrectly")
	assertEq(t, cmd.GetAuthor(), "vndaba", "Cmd author set wrongly")
	assertEq(t, cmd.GetVersion(), "0.1.0", "Cmd Version set wrongly")
	assertDeepEq(t, cmd.GetAliases(), []string{"b"}, "Cmd aliases set wrongly")
	assertDeepEq(t, cmd.GetFlags()[0], helpFlag(), "Help flag not added automatically")
}

func TestCommandSettings(t *testing.T) {
	clearCache()
	app := App()
	app.SubCommand("dummy")

	// assertStructArrEq[*Flag](t, app.GetFlags(), []*Flag{helpFlag(), versionFlag()}, "Help and version flags not set correctly")
	assertEq(t, len(app.GetFlags()), 2, "Help and version flags not set correctly")
	assertDeepEq(t, app.GetFlags()[1], versionFlag(), "Version flag not set correctly")

	app.Set(DisableVersionFlag, true)
	app.Set(IncludeHelpSubcommand, true)
	// TODO: Complete the other settings
	app._init()

	assertEq(t, len(app.GetFlags()), 1, "Failed to disable the version flag")
	assertEq(t, app.GetSubCommands()[1].GetName(), "help", "Include Help subcommand setting failed")
}

func TestEventListeners(t *testing.T) {
	app := App()
	app.SubCommand("test")

	app.Override(UnknownCommand, func(ec *EventConfig) {
		assertEq(t, len(ec.GetArgs()), 1, "Incorrect number of args passed along")
		assertEq(t, ec.GetEvent(), UnknownCommand, "Event on EventCfg set incorrectly")
		assertEq(t, ec.GetApp(), app, "app reference passed incorrectly")
		assertEq(t, ec.GetError().message, "no such subcommand found: `new`", "Error message configured wrongly")
		assertEq(t, ec.GetExitCode(), 40, "Wrong exit code found")
	})

	app.Set(IncludeHelpSubcommand, true)
	app.ParseFrom([]string{"my_bin", "new"})
}

func TestRootCmdArgs(t *testing.T) {
	// Test single required arg
	{
		clearCache()
		app := App()

		app.Argument("<file>", "file to open").
			Action(func(pm *ParserMatches) {
				val, _ := pm.GetArgValue("file")
				assertEq(t, val, "happy.png", "Root cmd arg parsing for required args is faulty")
			})

		app.ParseFrom([]string{"my_bin", "happy.png"})
	}

	// Test required arg error
	{
		clearCache()
		app := App()
		app.Argument("<file>", "file to open")

		expectedError := generateError(app, MissingRequiredArgument, []string{"<file>"})
		exec := func() {
			app.ParseFrom([]string{"my_bin"})
		}

		assertStdOut(t, expectedError.GetErrorString(app), exec, "Error throwing for missing required arg on root cmd faulty")
	}

	// Test optional args parsing
	{
		clearCache()
		app := App()
		app.Argument("[file]", "file to open").
			Action(func(pm *ParserMatches) {
				val, _ := pm.GetArgValue("file")
				assertEq(t, val, "happy.png", "Root cmd arg parsing for required args is faulty")
			})

		app.ParseFrom([]string{"my_bin", "happy.png"})
	}

	// Test multiple root args
	{
		app := App()
		app.Argument("<arg1>", "first arg").
			Argument("<arg2>", "second arg").
			Argument("<arg3>", "third arg").
			Action(func(pm *ParserMatches) {
				val1, _ := pm.GetArgValue("arg1")
				val2, _ := pm.GetArgValue("arg2")
				val3, _ := pm.GetArgValue("arg3")

				assertEq(t, val1, "one", "Root cmd multi-arg parsing failed")
				assertEq(t, val2, "two", "Root cmd multi-arg parsing failed")
				assertEq(t, val3, "three", "Root cmd multi-arg parsing failed")
			})

		app.ParseFrom([]string{"my_bin", "one", "two", "three"})
	}

}

func BenchmarkBuildEmptyCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCommand("empty")
	}
}

func BenchmarkBuildEmptyApp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		App()
	}
}

func _compareVariants(b *testing.B, vars ...func(*Command)) {
	for i, fn := range vars {
		name := ""
		switch i {
		case 0:
			name = "constructor"
		case 1:
			name = "builder"
		case 2:
			name = "composite-func"
		}
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				fn(NewCommand("empty"))
			}
		})
	}
}

func BenchmarkBuildWithFlags(b *testing.B) {
	constructor := func(a *Command) {
		a.AddFlag(&Flag{
			Name:     "verbose",
			LongVal:  "--verbose",
			ShortVal: "-V",
			HelpStr:  "Verbosity flag",
		})
	}

	builder := func(a *Command) {
		a.AddFlag(NewFlag("verbose").Help("Verbosity flag").Short('V'))
	}

	composite := func(a *Command) {
		a.Flag("-V --verbose", "Verbosity flag")
	}

	_compareVariants(b, constructor, builder, composite)
}

func BenchmarkBuildWithArgs(b *testing.B) {
	constructor := func(a *Command) {
		a.AddArgument(&Argument{
			Name:       "arg1",
			HelpStr:    "Argument one",
			RawValue:   "<arg1...>",
			ArgType:    str,
			IsVariadic: true,
			IsRequired: true,
		})
	}

	buidler := func(a *Command) {
		a.AddArgument(NewArgument("arg1").Help("Argument one").Required(true).Variadic(true))
	}

	composite := func(a *Command) {
		a.Argument("<arg1...>", "Argument one")
	}

	_compareVariants(b, constructor, buidler, composite)
}
