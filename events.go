package gommander

type EventListener = func(EventConfig)

type EventConfig struct {
}

type EventEmitter struct {
	listeners []EventListener
}

func NewEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: []EventListener{},
	}
}
