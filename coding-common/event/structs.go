package event_x

import "slices"

// Event 事件类型基类
type Event[T any] struct {
    //事件类型
    Type string
    //事件携带数据源
    Data T
}

// Clone 克隆事件
func (event *Event[T]) Clone() *Event[T] {
    return &Event[T]{Type: event.Type, Data: event.Data}
}

//func (event *Event) ToString() string {
//    return fmt.Sprintf("Event Type %v", event.Type)
//}

// IEventDispatcher 事件调度接口
type IEventDispatcher[T any] interface {
    // RegisterListener 事件监听
    RegisterListener(eventType string, listener *EventListener[T])
    // RemoveEventListener 移除事件监听
    RemoveEventListener(eventType string, listener *EventListener[T]) bool
    // HasEventListener 是否包含事件
    HasEventListener(eventType string) bool
    // DispatchEvent 事件派发
    DispatchEvent(event Event[T], async bool) bool
}

// EventDispatcher 事件调度器基类
type EventDispatcher[T any] struct {
    m map[string][]*EventListener[T]
}

// RegisterListener 事件调度器添加事件
func (e *EventDispatcher[T]) RegisterListener(eventType string, listener *EventListener[T]) {
    e.m[eventType] = append(e.m[eventType], listener)
}

// RemoveEventListener 事件调度器移除某个监听
func (e *EventDispatcher[T]) RemoveEventListener(eventType string, listener *EventListener[T]) bool {
    s := slices.DeleteFunc(e.m[eventType], func(l *EventListener[T]) bool {
        return listener == l
    })
    e.m[eventType] = s
    return true
}

// HasEventListener 事件调度器是否包含某个类型的监听
func (e *EventDispatcher[T]) HasEventListener(eventType string) bool {
    _, ok := e.m[eventType]
    return ok
}

// DispatchEvent 事件调度器派发事件
func (e *EventDispatcher[T]) DispatchEvent(event Event[T], async bool) bool {
    for _, listener := range e.m[event.Type] {
        if async {
            go listener.listen(event)
        } else {
            listener.listen(event)
        }
    }
    return true
}

// EventListener 监听器
type EventListener[T any] struct {
    listen func(event Event[T])
}
