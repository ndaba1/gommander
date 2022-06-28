package gommander

import "testing"

func TestOptionsCreation(t *testing.T) {
	opt := NewOption("port").Short('p').Help("The port option").Argument("<port-number>").Required(true)
	optB := newOption("-p --port <port-number>", "The port option", true)

	assertStructEq[*Option](t, opt, &optB, "Option creation methods out of sync")
	assert(t, opt.compare(&optB)) // linter workaround
	assert(t, opt.required, "Failed to set required value on options")

	expL := "-p, --port <port-number> "
	expF := "The port option"
	gotL, gotF := opt.generate()

	assertEq(t, expL, gotL, "The option generate func is faulty")
	assertEq(t, expF, gotF, "The option generate func is faulty")

	expectedArg := NewArgument("<port-number>")
	assertStructEq[*Argument](t, expectedArg, opt.args[0], "Option args created incorrectly")
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
