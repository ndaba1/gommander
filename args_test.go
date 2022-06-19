package gommander

import "testing"

func TestArgsCreation(t *testing.T) {

	/********************* Required arguments tests ********************/

	arg := NewArgument("<test>").Help("Test argument").Variadic(true)
	argB := newArgument("<test...>", "Test argument")

	if !arg.compare(argB) {
		t.Errorf("Arg creation methods out of sync: 1: %v 2: %v", arg, argB)
	}

	if !arg.isRequired {
		t.Errorf("Failed to make arg required")
	}

	if !arg.isVariadic {
		t.Errorf("Failed to make arg variadic")
	}

	if arg.name != "test" {
		t.Errorf("Arg name not set correctly")
	}

	if arg.help != "Test argument" {
		t.Errorf("Arg help string wrongly set")
	}

	arg.ValidateWith([]string{"ONE", "TWO"})

	if !arg.testValue("one") && arg.testValue("TWO") {
		t.Errorf("Arg validation working incorrectly")
	}

	if arg.getRawValue() != "<test...>" {
		t.Errorf("Raw value return function working incorrectly")
	}

	// Other tests
	expL := "<test...>"
	expF := "Test argument"

	if l, f := arg.generate(); l != expL || f != expF {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", expL, expF, l, f)
	}

	/********************* Optional arguments tests ********************/

	arg = NewArgument("[optional]").Default("DEFAULT").Help("Optional value with default")

	if arg.isRequired {
		t.Error("Failed to set argument as optional")
	}

	if !arg.hasDefaultValue() || arg.defaultValue != "DEFAULT" {
		t.Error("Failed to set default value for argument")
	}

	expL = "[optional]"
	expF = "Optional value with default (default: DEFAULT)"

	if l, f := arg.generate(); l != expL || f != expF {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", expL, expF, l, f)
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

func BenchmarkComplexArgBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewArgument("complex").
			Help("Arg with many options").
			Required(true).
			ValidateWith([]string{"ONE", "TEN"}).
			Default("TEN").
			Variadic(false)
	}
}

func BenchmarkNewArgFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newArgument("<test>", "A test argument")
	}
}
