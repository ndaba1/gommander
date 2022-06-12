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

var EVENTS_SLICE = []Event{
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
	args        []string
	event       Event
	app_ref     *Command
	exit_code   int
	err         GommanderError
	matched_cmd *Command
}

type EventEmitter struct {
	listeners          map[Event][]EventListener
	events_to_override []Event
}

func (c *EventConfig) GetArgs() []string        { return c.args }
func (c *EventConfig) GetEvent() Event          { return c.event }
func (c *EventConfig) GetApp() *Command         { return c.app_ref }
func (c *EventConfig) GetExitCode() int         { return c.exit_code }
func (c *EventConfig) GetError() GommanderError { return c.err }

func new_emitter() EventEmitter {
	return EventEmitter{
		listeners: make(map[Event][]EventListener),
	}
}

func (em *EventEmitter) on(event Event, cb EventCallback, pstn int) {
	if len(em.listeners[event]) == 0 {
		em.listeners[event] = []EventListener{{cb, pstn}}
	} else {
		new_v := em.listeners[event]

		for e := range em.listeners {
			if e == event {
				new_v = append(new_v, EventListener{cb, pstn})
			}
		}

		em.listeners[event] = new_v
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

			os.Exit(cfg.exit_code)
		}
	}
}

func (em *EventEmitter) override(e Event) {
	em.events_to_override = append(em.events_to_override, e)
}

func (em *EventEmitter) rm_default_lstnr(e Event) {
	new_arr := []EventListener{}
	for k, v := range em.listeners {
		if k == e {
			for _, lstnr := range v {
				if lstnr.index != -4 {
					new_arr = append(new_arr, lstnr)
				}
			}
		}
	}
	em.listeners[e] = new_arr
}

func (em *EventEmitter) insert_before_all(cb EventCallback) {
	for _, e := range EVENTS_SLICE {
		em.on(e, cb, -5)
	}
}

func (em *EventEmitter) insert_after_all(cb EventCallback) {
	for _, e := range EVENTS_SLICE {
		em.on(e, cb, 5)
	}
}

func (em *EventEmitter) on_errors(cb EventCallback) {
	for _, e := range EVENTS_SLICE {
		if e == OutputHelp || e == OutputVersion {
			continue
		}
		em.on(e, cb, -4)
	}
}
