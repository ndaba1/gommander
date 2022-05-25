package gommander

import "testing"

func TestOptionsCreation(t *testing.T) {
	opt := NewOption("port").Short('p').Help("The port option").Argument("<port-number>").Required(true)
	opt_b := new_option("-p --port <port-number>", "The port option", true)

	if opt.name != opt_b.name || opt.short != opt_b.short || opt.help != opt_b.help {
		t.Error("Option creation methods out of sync")
	}

	if !opt.required || !opt_b.required {
		t.Error("Failed to set required value on options")
	}

	expected_arg := NewArgument("<port-number>")
	if opt.args[0].name != expected_arg.name {
		t.Error("Option args created incorrectly")
	}
}

func BenchmarkOptBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewOption("port").
			Short('p').
			Help("A port option").
			Required(true)
	}
}

func BenchmarkOptBuilderwArg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewOption("port").
			Short('p').
			Help("A port option").
			Required(true).
			Argument("<port-number>")
	}
}

func BenchmarkOptArgBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewOption("port").
			Short('p').
			Help("A port option").
			Required(true).
			AddArgument(
				NewArgument("port-number").
					Required(true),
			)
	}
}

func BenchmarkNewOptFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new_option("-p --port <port-number>", "A port option", true)
	}
}
