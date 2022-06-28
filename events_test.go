package gommander

import "testing"

func TestListenerCreation(t *testing.T) {
	em := newEmitter()

	em.on(OutputHelp, func(ec *EventConfig) {}, 0)
	em.on(OutputVersion, func(ec *EventConfig) {}, 0)

	assert(t, len(em.listeners) == 2, "Failed to add listeners correctly")
}

func TestBeforeAllFn(t *testing.T) {
	em := newEmitter()

	em.insertBeforeAll(func(ec *EventConfig) {})

	for _, v := range em.listeners {
		assert(t, v[0].index == -5, "Failed to add before all listener")
	}
}

func TestAfterAllFn(t *testing.T) {
	em := newEmitter()

	em.insertAfterAll(func(ec *EventConfig) {})

	for _, v := range em.listeners {
		assert(t, v[0].index == 5, "Failed to add after all listener")
	}
}

func TestEmitterFunctionality(t *testing.T) {
	em := newEmitter()

	// Add some basic listeners
	em.onErrors(func(ec *EventConfig) {})

	for _, v := range em.listeners {
		if v[0].index != -4 {
			t.Error("Failed to add default listener")
		}
	}
}

func BenchmarkLstnrCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		em := newEmitter()
		em.on(OutputHelp, func(ec *EventConfig) {}, 0)
	}
}

func BenchmarkBatchLstnrs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		em := newEmitter()
		em.insertBeforeAll(func(ec *EventConfig) {})
	}
}
