package event_x

// NewEventDispatcher 创建事件派发器
func NewEventDispatcher[T any]() *EventDispatcher[T] {
    return &EventDispatcher[T]{m: map[string][]*EventListener[T]{}}
}

// NewEventListener 创建监听器
func NewEventListener[T any](h func(event Event[T])) *EventListener[T] {
    return &EventListener[T]{listen: h}
}

// NewEvent 创建事件
func NewEvent[T any](eventType string, data T) Event[T] {
    return Event[T]{Type: eventType, Data: data}
}
