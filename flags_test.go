package gommander

import "testing"

func TestFlagsCreation(t *testing.T) {
	flag := NewFlag("help").Short('h').Help("The help flag")
	flagB := newFlag("-h --help", "The help flag")

	assertStructEq[*Flag](t, flag, &flagB, "Flag creation functions are out of sync")

	flag.Global(true)
	assert(t, flag.isGlobal, "Failed to set flag as global")

	expL := "-h, --help"
	expF := "The help flag"
	gotL, gotF := flag.generate()

	assertEq(t, expL, gotL, "Flag generate method functioning incorrectly")
	assertEq(t, expF, gotF, "Flag generate method functioning incorrectly")
}

func BenchmarkFlagsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFlag("version").Short('V').Help("A version flag")
	}
}

func BenchmarkNewFlagFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		newFlag("-V --version", "A version flag")
	}
}
