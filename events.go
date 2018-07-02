package netco

import "fmt"

// 事件类型 基类
type Event struct {
	// 事件 触发器
	Target IEventDispatcher
	// 事件类型
	Type string
	// 事件数据
	Object interface{}
}

// 事件调度器 基类
type EventDispatcher struct {
	savers []*EventSaver
}

// 调度器存放单元
type EventSaver struct {
	Type      string
	Listeners []*EventListener
}

// 事件监听器
type EventListener struct {
	Handler EventHandler
}

// 监听器函数
type EventHandler func(event Event)

// 事件调度接口
type IEventDispatcher interface {
	AddEventListener(eventType string, listener *EventListener)
	RemoveEventListener(eventType string, listener *EventListener) bool
	HasEventListener(eventType string) bool
	DispatchEvent(event Event) bool
}

// 创建事件派发器
func NewEventDispatcher() *EventDispatcher {
	return new(EventDispatcher)
}

// 创建监听器
func NewEventListener(h EventHandler) *EventListener {
	l := new(EventListener)
	l.Handler = h
	return l
}

// 创建事件
func NewEvent(eventType string, object interface{}) Event {
	return Event{Type:eventType, Object:object}
}

// 克隆事件
func (this *Event) Clone() *Event {
	e := new(Event)
	e.Type = this.Type
	e.Target = this.Target
	return e
}

func (this *Event) ToString() string {
	return fmt.Sprintf("Event Type %v", this.Type)
}

// 事件调度器
// 添加事件
func (this *EventDispatcher) AddEventListener (eventType string, listener *EventListener) {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			saver.Listeners = append(saver.Listeners, listener)
			return
		}
	}
	saver := &EventSaver{Type:eventType, Listeners:[]*EventListener{listener}}
	this.savers = append(this.savers, saver)
}

// 移除事件
func (this *EventDispatcher) RemoveEventListener (eventType string, listner* EventListener) bool {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			for i, l := range saver.Listeners {
				if listner == l {
					saver.Listeners = append(saver.Listeners[:i], saver.Listeners[i+1:]...)
					return true
				}
			}
		}
	}
	return false
}

// 是否包含某个事件
func (this *EventDispatcher) HasEventListener(eventType string) bool {
	for _, saver := range this.savers {
		if saver.Type == eventType {
			return true;
		}
	}
	return false
}

// 调度事件
func (this *EventDispatcher) DispatchEvent(event Event) bool {
	for _, saver := range this.savers {
		if saver.Type == event.Type {
			for _, listener := range saver.Listeners {
				event.Target = this
				listener.Handler(event)
			}
			return true
		}
	}
	return false
}