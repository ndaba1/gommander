package gommander

import "testing"

func TestFlagsCreation(t *testing.T) {
	flag := NewFlag("help").Short('h').Help("The help flag")
	flag_2 := new_flag("-h --help", "The help flag")

	if !flag.compare(&flag_2) {
		t.Errorf("Flag creation functions out of sync: 1. %v  2. %v",
			flag, flag_2,
		)
	}

	flag.Global(true)
	if !flag.is_global {
		t.Error("Failed to set flag as global")
	}

	exp_l := "-h, --help"
	exp_f := "The help flag"

	if l, f := flag.generate(); l != exp_l || f != exp_f {
		t.Errorf("Flag generate functioning incorrectly. Expected (%v, %v), but found (%v, %v)", exp_l, exp_f, l, f)
	}
}

func BenchmarkFlagsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewFlag("version").Short('V').Help("A version flag")
	}
}

func BenchmarkNewFlagFn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		new_flag("-V --version", "A version flag")
	}
}
