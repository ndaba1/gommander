package gommander

import "testing"

func TestBasicListener(t *testing.T) {
	em := new_emitter()

	em.on(OutputHelp, func(ec *EventConfig) {}, 0)

	if len(em.listeners) == 0 {
		t.Error("Failed to add listener")
	}
}

func TestBeforeAllFn(t *testing.T) {
	em := new_emitter()

	em.insert_before_all(func(ec *EventConfig) {})

	for _, v := range em.listeners {
		if v[0].index != -5 {
			t.Errorf("Failed to add before all listener")
		}
	}
}

func TestAfterAllFn(t *testing.T) {
	em := new_emitter()

	em.insert_after_all(func(ec *EventConfig) {})

	for _, v := range em.listeners {
		if v[0].index != 5 {
			t.Errorf("Failed to add after all listener")
		}
	}
}
