package event

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetPayload() interface{}
	SetPayload(payload interface{})
	GetDateTime() time.Time
}

type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(evenName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
