package event_x

import (
	"log"
	"testing"
	"time"
)

const HELLO_WORLD = "helloWorld"

func Test_00(t *testing.T) {
	listener1 := NewEventListener(func(event Event[string]) {
		time.Sleep(time.Second * 2)
		log.Println("监听器1", time.Now(), event.Type, event.Data)
	})
	listener2 := NewEventListener(func(event Event[string]) {
		time.Sleep(time.Second * 2)
		log.Println("监听器2", time.Now(), event.Type, event.Data)
	})

	dispatcher := NewEventDispatcher[string]()
	dispatcher.RegisterListener(HELLO_WORLD, listener1)
	dispatcher.RegisterListener(HELLO_WORLD, listener2)

	// time.Sleep(time.Second * 2)
	// dispatcher.RemoveEventListener(HELLO_WORLD, listener)

	dispatcher.DispatchEvent(NewEvent(HELLO_WORLD, "data"), true)

	time.Sleep(time.Second * 3)
}
