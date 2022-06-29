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
	assertArrEq(t, cmd.GetAliases(), []string{"b"}, "Cmd aliases set wrongly")
	assertStructEq[*Flag](t, cmd.GetFlags()[0], helpFlag(), "Help flag not added automatically")
}

func TestCommandSettings(t *testing.T) {
	app := App()
	app.SubCommand("dummy")

	// assertStructArrEq[*Flag](t, app.GetFlags(), []*Flag{helpFlag(), versionFlag()}, "Help and version flags not set correctly")
	assertEq(t, len(app.GetFlags()), 2, "Help and version flags not set correctly")
	assertStructEq[*Flag](t, app.GetFlags()[1], versionFlag(), "Version flag not set correctly")

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

func BenchmarkBuildEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCommand("empty")
	}
}
