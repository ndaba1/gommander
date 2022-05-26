package gommander

import (
	"sort"
)

type Event byte
type EventCallback = func(*EventConfig)

const (
	MissingRequiredArgument Event = iota
	OutputHelp
	OutputVersion
	UnknownCommand
	UnknownOption
	UnresolvedArgument
	InvalidArgumentValue
)

var EVENTS_SLICE = []Event{
	MissingRequiredArgument,
	OutputHelp, OutputVersion,
	UnknownCommand, UnknownOption,
	UnresolvedArgument, InvalidArgumentValue,
}

type EventListener struct {
	cb    EventCallback
	index int
}

type EventConfig struct {
	args      []string
	event     Event
	app_ref   *Command
	exit_code int
	err       string
}

type EventEmitter struct {
	listeners map[Event][]EventListener
}

func (c *EventConfig) GetArgs() *[]string { return &c.args }
func (c *EventConfig) GetEvent() Event    { return c.event }
func (c *EventConfig) GetApp() *Command   { return c.app_ref }
func (c *EventConfig) GetExitCode() int   { return c.exit_code }
func (c *EventConfig) GetError() string   { return c.err }

func new_emitter() EventEmitter {
	return EventEmitter{
		listeners: make(map[Event][]EventListener),
	}
}

func (em *EventEmitter) on(event Event, cb EventCallback, pstn int) {
	if len(em.listeners) == 0 {
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
		}
	}
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
