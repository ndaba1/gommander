package gommander

import "testing"

func TestOptionsCreation(t *testing.T) {
	opt := NewOption("port").Short('p').Help("The port option").Argument("<port-number>").Required(true)
	optB := newOption("-p --port <port-number>", "The port option", true)

	if !opt.compare(&optB) {
		t.Errorf("Option creation methods out of sync: 1: %v  2: %v", opt, optB)
	}

	if !opt.required {
		t.Error("Failed to set required value on options")
	}

	expL := "-p, --port <port-number> "
	expF := "The port option"

	if l, f := opt.generate(); expF != f || expL != l {
		t.Errorf("The option generate func is faulty. Expected (%v, %v), but found (%v, %v)", expL, expF, l, f)
	}

	expectedArg := NewArgument("<port-number>")
	if opt.args[0].name != expectedArg.name {
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
		newOption("-p --port <port-number>", "A port option", true)
	}
}
