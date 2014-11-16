package noty

import (
	"fmt"
	"sync"
)

type SignalHandler func(param interface{}) (interface{}, error)

type AsyncCall struct {
	param   interface{}
	retChan chan<- interface{}
	errChan chan<- error
}

var handlers = make(map[string]SignalHandler)
var asyncCallChans = make(map[string]chan AsyncCall)
var lck = sync.RWMutex{}

func Register(eventName string, fn SignalHandler) {
	lck.Lock()
	handlers[eventName] = fn
	c := make(chan AsyncCall)
	asyncCallChans[eventName] = c
	lck.Unlock()

	go func(eventName string, c chan AsyncCall) {
		for param := range c {
			r, err := fn(param.param)
			if err != nil && param.errChan != nil {
				param.errChan <- err
				continue
			}
			if param.retChan != nil {
				param.retChan <- r
			}
		}
	}(eventName, c)
}

func SignalSync(eventName string, param interface{}) (interface{}, error) {
	lck.RLock()
	fn := handlers[eventName]
	lck.RUnlock()

	if fn != nil {
		r, err := fn(param)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	return nil, fmt.Errorf("no handler for event : %s", eventName)
}

func SignalAsync(eventName string, param interface{}, retChan chan<- interface{}, errChan chan<- error) {
	lck.RLock()
	c := asyncCallChans[eventName]
	lck.RUnlock()

	c <- AsyncCall{
		param:   param,
		retChan: retChan,
		errChan: errChan,
	}
}
