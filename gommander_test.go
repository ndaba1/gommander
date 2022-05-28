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

func BenchmarkBuildEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewCommand("empty")
	}
}
