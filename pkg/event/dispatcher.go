package event

import (
	"errors"
	"sync"
)

var (
	ErrEventNotRegistered       = errors.New("event not registered")
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrHandlerNotFound          = errors.New("handler not found")
)

type Dispatcher struct {
	handlers map[string][]HandlerInterface
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]HandlerInterface),
	}
}

func (d *Dispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		var wg *sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
		return nil
	}
	return ErrEventNotRegistered
}

func (d *Dispatcher) Register(event EventInterface, handler HandlerInterface) error {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for _, registeredHandler := range handlers {
			if handler == registeredHandler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	d.handlers[event.GetName()] = append(d.handlers[event.GetName()], handler)
	return nil
}

func (d *Dispatcher) Has(event EventInterface, handler HandlerInterface) (bool, error) {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for _, registeredHandler := range handlers {
			if handler == registeredHandler {
				return true, nil
			}
		}
		return false, nil
	}
	return false, ErrEventNotRegistered
}

func (d *Dispatcher) Remove(event EventInterface, handler HandlerInterface) error {
	if handlers, ok := d.handlers[event.GetName()]; ok {
		for i, registeredHandler := range handlers {
			if handler == registeredHandler {
				d.handlers[event.GetName()] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}
		return ErrHandlerNotFound
	}
	return ErrEventNotRegistered
}

func (d *Dispatcher) Clear() {
	d.handlers = make(map[string][]HandlerInterface)
}
