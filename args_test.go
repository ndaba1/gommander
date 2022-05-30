package gommander

import "testing"

func TestArgsCreation(t *testing.T) {

	/********************* Required arguments tests ********************/

	arg := NewArgument("<test>").Help("Test argument").Variadic(true)
	arg_b := new_argument("<test...>", "Test argument")

	if !arg.compare(arg_b) {
		t.Errorf("Arg creation methods out of sync: 1: %v 2: %v", arg, arg_b)
	}

	if !arg.is_required {
		t.Errorf("Failed to make arg required")
	}

	if !arg.is_variadic {
		t.Errorf("Failed to make arg variadic")
	}

	if arg.name != "test" {
		t.Errorf("Arg name not set correctly")
	}

	if arg.help != "Test argument" {
		t.Errorf("Arg help string wrongly set")
	}

	arg.ValidateWith([]string{"ONE", "TWO"})

	if !arg.test_value("one") && arg.test_value("TWO") {
		t.Errorf("Arg validation working incorrectly")
	}

	if arg.get_raw_value() != "<test...>" {
		t.Errorf("Raw value return function working incorrectly")
	}

	// Other tests
	exp_l := "<test...>"
	exp_f := "Test argument"

	if l, f := arg.generate(); l != exp_l || f != exp_f {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", exp_l, exp_f, l, f)
	}

	/********************* Optional arguments tests ********************/

	arg = NewArgument("[optional]").Default("DEFAULT").Help("Optional value with default")

	if arg.is_required {
		t.Error("Failed to set argument as optional")
	}

	if !arg.has_default_value() || arg.default_value != "DEFAULT" {
		t.Error("Failed to set default value for argument")
	}

	exp_l = "[optional]"
	exp_f = "Optional value with default"

	if l, f := arg.generate(); l != exp_l || f != exp_f {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", exp_l, exp_f, l, f)
	}

}

func BenchmarkArgsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewArgument("test").
			Help("A test argument").
			Required(true).
			ValidateWith([]string{"ONE", "TWO", "THREE"}).
			Variadic(true)
	}
}

func BenchmarkNewArgFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new_argument("<test>", "A test argument")
	}
}
