package gommander

import (
	"os"
	"sort"
)

type Event byte
type EventCallback = func(*EventConfig)

const (
	// Only one argument passed along, the name of the argument in the form <arg>
	MissingRequiredArgument Event = iota
	// No arguments passed for this event
	OutputHelp
	// No arguments passed for this event
	OutputVersion
	// A single argument is passed for this event, the value of the unknown command
	UnknownCommand
	// A single argument is passed, the value of the unknown option
	UnknownOption
	// Single argument passed, the value of the unresolved argument
	UnresolvedArgument
	// Single argument: the value of the invalid argument
	InvalidArgumentValue
	// Single argument: the name of the missing option
	MissingRequiredOption
)

var eventsSlice = []Event{
	MissingRequiredArgument,
	OutputHelp, OutputVersion,
	UnknownCommand, UnknownOption,
	UnresolvedArgument, InvalidArgumentValue,
	MissingRequiredOption,
}

type EventListener struct {
	cb    EventCallback
	index int
}

type EventConfig struct {
	args       []string
	event      Event
	appRef     *Command
	exitCode   int
	err        Error
	matchedCmd *Command
}

type EventEmitter struct {
	listeners        map[Event][]EventListener
	eventsToOverride []Event
}

func (c *EventConfig) GetArgs() []string { return c.args }
func (c *EventConfig) GetEvent() Event   { return c.event }
func (c *EventConfig) GetApp() *Command  { return c.appRef }
func (c *EventConfig) GetExitCode() int  { return c.exitCode }
func (c *EventConfig) GetError() Error   { return c.err }

func newEmitter() EventEmitter {
	return EventEmitter{
		listeners: make(map[Event][]EventListener),
	}
}

func (em *EventEmitter) on(event Event, cb EventCallback, pstn int) {
	if len(em.listeners[event]) == 0 {
		em.listeners[event] = []EventListener{{cb, pstn}}
	} else {
		newLsntr := em.listeners[event]

		for e := range em.listeners {
			if e == event {
				newLsntr = append(newLsntr, EventListener{cb, pstn})
			}
		}

		em.listeners[event] = newLsntr
	}

}

func (em *EventEmitter) emit(cfg EventConfig) {
	event := cfg.GetEvent()

	for e, v := range em.listeners {
		if e == event {
			sort.SliceStable(v, func(i, j int) bool {
				return v[j].index > v[i].index
			})

			for _, lstnr := range v {
				lstnr.cb(&cfg)
			}

			os.Exit(cfg.exitCode)
		}
	}
}

func (em *EventEmitter) override(e Event) {
	em.eventsToOverride = append(em.eventsToOverride, e)
}

func (em *EventEmitter) rmDefaultLstnr(e Event) {
	newArr := []EventListener{}
	for k, v := range em.listeners {
		if k == e {
			for _, lstnr := range v {
				if lstnr.index != -4 {
					newArr = append(newArr, lstnr)
				}
			}
		}
	}
	em.listeners[e] = newArr
}

func (em *EventEmitter) insertBeforeAll(cb EventCallback) {
	for _, e := range eventsSlice {
		em.on(e, cb, -5)
	}
}

func (em *EventEmitter) insertAfterAll(cb EventCallback) {
	for _, e := range eventsSlice {
		em.on(e, cb, 5)
	}
}

func (em *EventEmitter) onErrors(cb EventCallback) {
	for _, e := range eventsSlice {
		if e == OutputHelp || e == OutputVersion {
			continue
		}
		em.on(e, cb, -4)
	}
}
