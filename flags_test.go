package gommander

import "testing"

func TestFlagsCreation(t *testing.T) {
	flag := NewFlag("help").Short('h').Help("The help flag")
	flag_2 := new_flag("-h --help", "The help flag")

	if flag.help != flag_2.help ||
		flag.long != flag_2.long ||
		flag.short != flag_2.short {
		t.Errorf("Flag creation functions out of sync: 1. (%s, %s - %s)  2. (%s, %s - %s)",
			flag.short, flag.long, flag.help,
			flag_2.short, flag_2.long, flag_2.help,
		)
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
