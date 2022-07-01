package gommander

import "testing"

func TestFlagsCreation(t *testing.T) {
	flag := NewFlag("help").Short('h').Help("The help flag").Global(true)
	flagB := newFlag("-h --help", "The help flag")

	assertStructEq[*Flag](t, flag, &flagB, "Flag creation functions are out of sync")
	assert(t, flag.compare(&flagB)) // linter workaround
	assert(t, flag.IsGlobal, "Failed to set flag as global")

	expL := "-h, --help"
	expF := "The help flag"
	gotL, gotF := flag.generate()

	assertEq(t, expL, gotL, "Flag generate method functioning incorrectly")
	assertEq(t, expF, gotF, "Flag generate method functioning incorrectly")
}

func BenchmarkFlagBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFlag("version").Short('V').Help("A version flag")
	}
}

func BenchmarkFlagFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newFlag("-V --version", "A version flag")
	}
}

func BenchmarkFlagConstructor(b *testing.B) {
	fn := func(f Flag) {}
	for i := 0; i < b.N; i++ {
		fn(Flag{
			Name:     "version",
			ShortVal: "-v",
			LongVal:  "--version",
			HelpStr:  "A version flag",
			IsGlobal: true,
		})
	}
}
