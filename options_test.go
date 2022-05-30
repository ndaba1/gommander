package gommander

import "testing"

func TestOptionsCreation(t *testing.T) {
	opt := NewOption("port").Short('p').Help("The port option").Argument("<port-number>").Required(true)
	opt_b := new_option("-p --port <port-number>", "The port option", true)

	if !opt.compare(&opt_b) {
		t.Errorf("Option creation methods out of sync: 1: %v  2: %v", opt, opt_b)
	}

	if !opt.required {
		t.Error("Failed to set required value on options")
	}

	exp_l := "-p, --port <port-number> "
	exp_f := "The port option"

	if l, f := opt.generate(); exp_f != f || exp_l != l {
		t.Errorf("The option generate func is faulty. Expected (%v, %v), but found (%v, %v)", exp_l, exp_f, l, f)
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
