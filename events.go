package netco

import (
	"fmt"
	"sync"
)

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
	listenerId int
	listenerMutex *sync.Mutex
}

// 调度器存放单元
type EventSaver struct {
	Type      string
	Listeners []*EventListener
}

// 事件监听器
type EventListener struct {
	Handler EventHandler
	id int
}

// 监听器函数
type EventHandler func(event Event)

// 事件调度接口
type IEventDispatcher interface {
	AddEventListener(eventType string, listener *EventListener) int
	RemoveEventListener(id int) bool
	HasEventListener(eventType string) bool
	DispatchEvent(event Event) bool
}

// 创建事件派发器
func NewEventDispatcher() *EventDispatcher {
	ed := new(EventDispatcher)
	ed.listenerId = 0
	ed.listenerMutex = new(sync.Mutex)
	return ed
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
func (this *EventDispatcher) AddEventListener (eventType string, listener *EventListener) int {
	this.listenerMutex.Lock()
	defer this.listenerMutex.Unlock()
	listener.id = this.listenerId
	this.listenerId++
	for _, saver := range this.savers {
		if saver.Type == eventType {
			saver.Listeners = append(saver.Listeners, listener)
			return listener.id
		}
	}
	saver := &EventSaver{Type:eventType, Listeners:[]*EventListener{listener}}
	this.savers = append(this.savers, saver)
	return listener.id
}

// 移除事件
func (this *EventDispatcher) RemoveEventListener (id int) bool {
	this.listenerMutex.Lock()
	defer this.listenerMutex.Unlock()
	for _, saver := range this.savers {
		for i, l := range saver.Listeners {
			if id == l.id {
				saver.Listeners = append(saver.Listeners[:i], saver.Listeners[i+1:]...)
				return true
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
	var result = false
	var targets []*EventListener

	this.listenerMutex.Lock()
	for _, saver := range this.savers {
		if saver.Type == event.Type {
			for _, listener := range saver.Listeners {
				targets = append(targets, listener)
			}
			result = true
			break
		}
	}
	this.listenerMutex.Unlock()
	for _, v := range targets {
		event.Target = this

		v.Handler(event)
	}

	if !result {
		logger.Debug("没有找到handler", event)
	}

	return result
}