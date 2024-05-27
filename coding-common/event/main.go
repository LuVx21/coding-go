package event_x

// NewEventDispatcher 创建事件派发器
func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{m: map[string][]*EventListener{}}
}

// NewEventListener 创建监听器
func NewEventListener(h func(event Event)) *EventListener {
    return &EventListener{listen: h}
}

// NewEvent 创建事件
func NewEvent(eventType string, object any) Event {
    return Event{Type: eventType, Data: object}
}
