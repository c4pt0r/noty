package noty

import (
	"testing"
	"time"

	log "github.com/ngaut/logging"
)

func TestNoty(t *testing.T) {
	Register("test_event", func(param interface{}) (interface{}, error) {
		log.Info(param)
		return param, nil
	})

	ret, err := SignalSync("test_event", "hello")
	if err != nil {
		t.Error(err)
	}
	if ret.(string) != "hello" {
		t.Error("sync call error")
	}

	retChan := make(chan interface{})
	errChan := make(chan error)
	SignalAsync("test_event", "hello async", retChan, errChan)

	select {
	case r := <-retChan:
		if r.(string) != "hello async" {
			t.Error("async call error")
		}
	case err := <-errChan:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Error("time out")
	}

	for i := 0; i < 100; i++ {
		SignalAsync("test_event", "hello async", nil, nil)
	}
}
