package gommander

import "testing"

func TestArgsCreation(t *testing.T) {
	arg := NewArgument("<test...>").Help("Test argument")
	arg_b := new_argument("<test>", "Test argument")

	if arg.name != arg_b.name || arg.help != arg_b.help {
		t.Errorf("Arg creation methods out of sync: first is: (%s - %s) and second is: (%s - %s)", arg.name, arg.help, arg_b.name, arg_b.help)
	}

	if !arg.is_required {
		t.Errorf("Failed to make arg required")
	}

	if !arg.variadic {
		t.Errorf("Failed to make arg variadic")
	}

	if arg.name != "test" {
		t.Errorf("Arg name not set correctly")
	}

	if arg.help != "Test argument" {
		t.Errorf("Arg help string wrongly set")
	}

	arg.ValidateWith([]string{"ONE", "TWO"})

	if !arg.ValueIsValid("one") {
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

func BenchmarkNewArgsFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new_argument("<test>", "A test argument")
	}
}
