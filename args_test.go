package gommander

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

func TestXxx(t *testing.T) {
	// dummy function to set test_mode env vars
	// Workaround to set test mode once for all other tests
	setGommanderTestMode()
}

func TestArgsCreation(t *testing.T) {
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

}

func TestArgsMetadata(t *testing.T) {
	arg := NewArgument("<basic>").
		Variadic(true).
		Help("Test argument").
		ValidateWith([]string{"ONE", "TWO"})

	if !arg.testValue("one") && arg.testValue("TWO") {
		t.Errorf("Arg validation working incorrectly")
	}
	if arg.getRawValue() != "<basic...>" {
		t.Errorf("Raw value return function working incorrectly")
	}

	expLeading := "<basic...>"
	expFloating := "Test argument"

	if l, f := arg.generate(); l != expLeading || f != expFloating {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", expLeading, expFloating, l, f)
	}
}

func TestOptionalArgs(t *testing.T) {
	arg := NewArgument("[optional]").Default("DEFAULT").Help("Optional value with default")

	if arg.isRequired {
		t.Error("Failed to set argument as optional")
	}
	if !arg.hasDefaultValue() || arg.defaultValue != "DEFAULT" {
		t.Error("Failed to set default value for argument")
	}

	expLeading := "[optional]"
	expFloating := "Optional value with default (default: DEFAULT)"

	if l, f := arg.generate(); l != expLeading || f != expFloating {
		t.Errorf("The arg generate function is problematic. Expected: (%v, %v) but found (%v, %v)", expLeading, expFloating, l, f)
	}
}

func TestArgValidValues(t *testing.T) {
	// valid values validator
	arg := NewArgument("<lang>").
		DisplayAs("language").
		ValidateWith([]string{"ENG", "SPA", "RUS", "FRE"})

	if arg.getRawValue() != "language" {
		t.Error("Failed to set raw arg value using DisplayAs method")
	}
	if arg.testValue("else") {
		t.Error("Values validation not working properly")
	}

	exec := func() {
		arg.Default("NEW")
	}
	expected := fmt.Sprintf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", arg.name, "NEW", arg.validValues)

	if !assertStdOut(expected, exec) {
		t.Error("Argument validation for arguments with valid values is buggy")
	}
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

	if arg.getRawValue() != "int" {
		t.Error("Failed to set raw arg value using DisplayAs method")
	}
	if !arg.testValue("2") {
		t.Error("Strconv validator function working incorrectly")
	}

	exec := func() {
		arg.Default("notInt")
	}
	expected := fmt.Sprintf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", arg.name, "notInt")

	if !assertStdOut(expected, exec) {
		t.Error("Argument validation for arguments with validator functions is buggy")
	}
}

func assertStdOut(expected string, exec func()) bool {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exec()
	_ = w.Close()
	res, _ := io.ReadAll(r)
	output := string(res)

	os.Stdout = stdOut

	return output == expected
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
