package gommander

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_(t *testing.T) {
	// dummy function to set test_mode env vars
	// Workaround to set test mode once for all other tests
	setGommanderTestMode()
}

func TestArgsCreation(t *testing.T) {
	arg := NewArgument("<test>").Help("Test argument").Variadic(true)
	argB := newArgument("<test...>", "Test argument")

	assertStructEq[*Argument](t, arg, argB, "Arg creation methods out of sync")
	assert(t, arg.compare(argB)) // workaround to suppress linter warnings
	assert(t, arg.isRequired, "Failed to make arg required")
	assert(t, arg.isVariadic, "Failed to make arg variadic")
	assertEq(t, arg.name, "test", "Arg name not set correctly")
	assertEq(t, arg.help, "Test argument", "Arg help string wrongly set")
}

func TestArgsMetadata(t *testing.T) {
	arg := NewArgument("<basic>").
		Variadic(true).
		Help("Test argument").
		ValidateWith([]string{"ONE", "TWO"})

	assertNe(t, arg.name, "<basic>", "Enclosures not stripped from name")
	assert(t, arg.testValue("one"), "Arg validation working incorrectly")
	assert(t, arg.testValue("TWO"), "Arg validation working incorrectly")
	assertEq(t, arg.getRawValue(), "<basic...>", "Raw value return function working incorrectly")

	expLeading := "<basic...>"
	expFloating := "Test argument"
	gotLeading, gotFloating := arg.generate()

	assertEq(t, expLeading, gotLeading, "The arg generate function is problematic")
	assertEq(t, expFloating, gotFloating, "The arg generate function is problematic")
}

func TestOptionalArgs(t *testing.T) {
	arg := NewArgument("[optional]").Default("DEFAULT").Help("Optional value with default")

	assert(t, !arg.isRequired, "Failed to set argument as optional")
	assert(t, arg.hasDefaultValue(), "Failed to set default value")
	assert(t, arg.defaultValue == "DEFAULT", "Failed to set default value")

	expLeading := "[optional]"
	expFloating := "Optional value with default (default: DEFAULT)"
	gotLeading, gotFloating := arg.generate()

	assertEq(t, expLeading, gotLeading, "The arg generate function is problematic")
	assertEq(t, expFloating, gotFloating, "The arg generate function is problematic")
}

func TestArgValidValues(t *testing.T) {
	// valid values validator
	arg := NewArgument("<lang>").
		DisplayAs("language").
		ValidateWith([]string{"ENG", "SPA", "RUS", "FRE"})

	assert(t, !arg.testValue("else"), "Values validation not working properly")
	assertEq(t, arg.getRawValue(), "language", "Failed to set raw arg value using DisplayAs method")

	exec := func() {
		arg.Default("NEW")
	}
	expected := fmt.Sprintf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", arg.name, "NEW", arg.validValues)
	assertStdOut(t, expected, exec, "Argument validation for arguments with valid values is buggy")
}

func TestArgValidatorFunc(t *testing.T) {
	arg := NewArgument("<age>").
		DisplayAs("int").
		ValidatorFunc(func(s string) error {
			_, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			return nil
		})

	assertEq(t, arg.getRawValue(), "int", "Failed to set raw arg value using DisplayAs method")
	assert(t, arg.testValue("2"), "Strconv validator function working incorrectly")

	exec := func() {
		arg.Default("notInt")
	}
	expected := fmt.Sprintf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", arg.name, "notInt")
	assertStdOut(t, expected, exec, "Argument validation for arguments with validator functions is buggy")
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
			DisplayAs("Something").
			Required(true).
			ValidateWith([]string{"ONE", "TEN"}).
			Default("TEN").
			Variadic(false).
			ValidatorFunc(func(s string) error {
				_, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				return nil
			})
	}
}

func BenchmarkArgConstructor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newArgument("<test>", "A test argument")
	}
}
