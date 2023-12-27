package events

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

type HandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type DispatcherInterface interface {
	Register(eventName string, handler HandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(evenName string, handler HandlerInterface) error
	Has(eventName string, handler HandlerInterface) (bool, error)
	Clear()
}
