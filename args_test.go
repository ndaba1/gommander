package gommander

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_init(t *testing.T) {
	setGommanderTestMode()
}

func TestArgCreation(t *testing.T) {
	arg := NewArgument("<test>").Help("Test argument").Variadic(true)
	argB := newArgument("<test...>", "Test argument")

	assertDeepEq(t, arg, argB, "Arg creation methods out of sync")
	assert(t, arg.IsRequired, "Failed to make arg required")
	assert(t, arg.IsVariadic, "Failed to make arg variadic")
	assertEq(t, arg.Name, "test", "Arg name not set correctly")
	assertEq(t, arg.HelpStr, "Test argument", "Arg help string wrongly set")
}

func TestArgMetadata(t *testing.T) {
	arg := NewArgument("<basic>").
		Variadic(true).
		Help("Test argument").
		ValidateWith([]string{"ONE", "TWO"}).
		Type(str)

	assertNe(t, arg.Name, "<basic>", "Enclosures not stripped from name")
	assert(t, arg.testValue("one"), "Arg validation working incorrectly")
	assert(t, arg.testValue("TWO"), "Arg validation working incorrectly")
	assertEq(t, arg.getRawValue(), "<basic...>", "Raw value return function working incorrectly")
	assertEq(t, arg.ArgType, str, "Arg builder Type method not working correctly")
	assertEq(t, len(arg.ValidatorFns), 0, "Unnecessary validator func added for strings")

	expLeading := "<basic...>"
	expFloating := "Test argument"
	gotLeading, gotFloating := arg.generate(App())

	assertEq(t, expLeading, gotLeading, "The arg generate function is problematic")
	assertEq(t, expFloating, gotFloating, "The arg generate function is problematic")
}

func TestOptionalArgs(t *testing.T) {
	arg := NewArgument("[optional]").Default("DEFAULT").Help("Optional value with default")

	assert(t, !arg.IsRequired, "Failed to set argument as optional")
	assert(t, arg.hasDefaultValue(), "Failed to set default value")
	assert(t, arg.DefaultValue == "DEFAULT", "Failed to set default value")

	expLeading := "[optional]"
	expFloating := "Optional value with default (default: DEFAULT)"
	gotLeading, gotFloating := arg.generate(App())

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
	expected := fmt.Sprintf("error occurred when setting default value for argument: %v \n.  the passed value %v does not match the valid values: %v", arg.Name, "NEW", arg.ValidValues)
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
	expected := fmt.Sprintf("you tried to set a default value for argument: %v, but the validator function returned an error for values: %v", arg.Name, "notInt")
	assertStdOut(t, expected, exec, "Argument validation for arguments with validator functions is buggy")
}

func TestArgRegexValidator(t *testing.T) {
	clearCache()
	{
		arg := NewArgument("version").ValidatorRegex(`^v[\d\.]`)
		assert(t, arg.testValue("v0.1.0"), "Regex validation is buggy")
		assert(t, !arg.testValue("1.1"), "Regex validation working incorrectly")
		assert(t, !arg.testValue("vA.1"), "Regex validation working incorrectly")
	}
}

func TestArgTypeValidation(t *testing.T) {
	{
		arg := NewArgument("<int:count>")
		assert(t, arg.testValue("2"), "Integer arg validation faulty")
		assert(t, arg.testValue("-2"), "Integer arg validation faulty")
		assertEq(t, arg.testValue("2.0"), false, "Integer arg validation faulty against float")
		assertEq(t, arg.testValue("two"), false, "Integer arg validation faulty against string")
	}
	{
		arg := NewArgument("<uint:count>")
		assert(t, arg.testValue("4"), "Uinteger arg validation is faulty")
		assert(t, !arg.testValue("-4"), "Uinteger arg validation is faulty")
	}
	{
		arg := NewArgument("<float:count>")
		assert(t, arg.testValue("2.0"), "Float arg validation faulty")
		assertEq(t, arg.testValue("two"), false, "Float arg validation faulty against string")
	}
	{
		arg := NewArgument("<bool:count>")
		assert(t, arg.testValue("true"), "Boolean arg validation faulty")
		assertEq(t, arg.testValue("2"), false, "Boolean arg validation faulty against int")
		assertEq(t, arg.testValue("2.0"), false, "Boolean arg validation faulty against float")
	}
	{
		arg := NewArgument("<file:path>")
		assert(t, !arg.testValue("fake.png"), "Filename arg validation faulty")
		assert(t, arg.testValue("go.mod"), "Filename arg validation faulty")
	}
	{
		exec := func() {
			NewArgument("<fake:arg>")
		}
		expected := "found unknown argument type: `fake` for argument: `<arg>`\n"

		assertStdOut(t, expected, exec, "Unknown arg types pass on undetected")
	}
}

func BenchmarkArgFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newArgument("<test>", "A test argument")
	}
}

func BenchmarkArgBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewArgument("test").
			Help("A test argument").
			Required(true).
			ValidateWith([]string{"ONE", "TWO", "THREE"}).
			Variadic(true)
	}
}

func BenchmarkArgBuilderFull(b *testing.B) {
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
	fn := func(a Argument) {}
	for i := 0; i < b.N; i++ {
		fn(Argument{
			Name:         "test",
			HelpStr:      "A test argument",
			IsRequired:   true,
			ValidValues:  []string{"ONE", "TWO", "THREE"},
			IsVariadic:   true,
			RawValue:     "<test...>",
			DefaultValue: "TEN",
			ArgType:      integer,
			ValidatorFns: [](func(s string) error){
				func(s string) error {
					_, err := strconv.Atoi(s)
					if err != nil {
						return err
					}
					return nil
				},
			},
		})
	}
}
