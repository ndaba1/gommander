package gommander

import "testing"

func TestCommandMetadata(t *testing.T) {
	cmd := NewCommand("basic").Help("Basic command")

	if cmd.GetHelp() != "Basic command" {
		t.Error("Cmd help field set incorrectly")
	}

	if len(cmd.GetFlags()) != 1 {
		t.Error("Help flag not set on command")
	}
}

func TestCommandSettings(t *testing.T) {
	app := App()

	if len(app.flags) != 2 {
		t.Error("Help and version flags not set correctly")
	}

	// Test DisableVersionFlag
	app.Set(DisableVersionFlag, true)
	app._init()

	if len(app.flags) != 1 {
		t.Error("Failed to disable version flag")
	}

}

func BenchmarkBuildEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCommand("empty")
	}
}
