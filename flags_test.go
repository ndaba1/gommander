package gommander

import "testing"

func TestFlagsCreation(t *testing.T) {
	flag := NewFlag("help").Short('h').Help("The help flag")
	flagB := newFlag("-h --help", "The help flag")

	if !flag.compare(&flagB) {
		t.Errorf("Flag creation functions out of sync: 1. %v  2. %v",
			flag, flagB,
		)
	}

	flag.Global(true)
	if !flag.isGlobal {
		t.Error("Failed to set flag as global")
	}

	expL := "-h, --help"
	expF := "The help flag"

	if l, f := flag.generate(); l != expL || f != expF {
		t.Errorf("Flag generate functioning incorrectly. Expected (%v, %v), but found (%v, %v)", expL, expF, l, f)
	}
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
