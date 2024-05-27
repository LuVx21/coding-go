package event_x

import "slices"

// Event 事件类型基类
type Event struct {
    //事件类型
    Type string
    //事件携带数据源
    Data any
}

// Clone 克隆事件
func (event *Event) Clone() *Event {
    e := new(Event)
    e.Type = event.Type
    return e
}

//func (event *Event) ToString() string {
//    return fmt.Sprintf("Event Type %v", event.Type)
//}

// IEventDispatcher 事件调度接口
type IEventDispatcher interface {
    // RegisterListener 事件监听
    RegisterListener(eventType string, listener *EventListener)
    // RemoveEventListener 移除事件监听
    RemoveEventListener(eventType string, listener *EventListener) bool
    // HasEventListener 是否包含事件
    HasEventListener(eventType string) bool
    // DispatchEvent 事件派发
    DispatchEvent(event Event, async bool) bool
}

// EventDispatcher 事件调度器基类
type EventDispatcher struct {
    m map[string][]*EventListener
}

// RegisterListener 事件调度器添加事件
func (e *EventDispatcher) RegisterListener(eventType string, listener *EventListener) {
    e.m[eventType] = append(e.m[eventType], listener)
}

// RemoveEventListener 事件调度器移除某个监听
func (e *EventDispatcher) RemoveEventListener(eventType string, listener *EventListener) bool {
    s := slices.DeleteFunc(e.m[eventType], func(l *EventListener) bool {
        return listener == l
    })
    e.m[eventType] = s
    return true
}

// HasEventListener 事件调度器是否包含某个类型的监听
func (e *EventDispatcher) HasEventListener(eventType string) bool {
    _, ok := e.m[eventType]
    return ok
}

// DispatchEvent 事件调度器派发事件
func (e *EventDispatcher) DispatchEvent(event Event, async bool) bool {
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
type EventListener struct {
    listen func(event Event)
}
