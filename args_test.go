package gommander

import "testing"

func TestArgsCreation(t *testing.T) {
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

	if !arg.test_value("one") {
		t.Errorf("Arg validation working incorrectly")
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
