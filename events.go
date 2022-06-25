package gommander

import (
	"os"
	"sort"
)

type Event byte
type EventCallback = func(*EventConfig)

const (
	// An event emitted when an argument that is marked as required, either on a command or on an option is not provided. Only one argument is passed along for this event when emitted, the name of the argument in the form <arg>
	MissingRequiredArgument Event = iota
	// Emitted when the help flag is invoked. No arguments passed for this event
	OutputHelp
	// Emitted when the version flag is invoked. No arguments passed for this event
	OutputVersion
	// This event is emitted when a value is passed as a subcommand, but no such subcommand could be resolved. A single argument is passed for this event, the value of the unknown command
	UnknownCommand
	// Emitted when an option-like value is provided, i.e. begins with `-` but no such option or flag was found. A single argument is passed, the value of the unknown option
	UnknownOption
	// A general event emitted when an argument was not expected. Single argument passed, the value of the unresolved argument
	UnresolvedArgument
	// This event occurs when an argument has a set of valid values or a validator function and the value provided does not match either. A Single value is passed along for this event: the value of the invalid argument
	InvalidArgumentValue
	// An event emitted when a required option is not provided. Single argument: the name of the missing option
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

			if !isTestMode() {
				os.Exit(cfg.exitCode)
			}
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
