package emit

import (
	"reflect"
	"sync"
)

// Handler 监听触发的回调函数
type Handler func(interface{})

type Emitter struct {
	lisMu    sync.RWMutex
	listener map[string][]Handler
}

// New 实例化监听器
func New() *Emitter {
	return &Emitter{
		listener: make(map[string][]Handler),
	}
}

// On 增加事件监听
func (e *Emitter) On(name string, handler ...Handler) *Emitter {
	e.lisMu.Lock()
	defer e.lisMu.Unlock()
	handlers, ok := e.listener[name]
	if !ok {
		handlers = handler
	} else {
		handlers = append(handlers, handler...)
	}
	e.listener[name] = handlers
	return e
}

// Off 取消监听
// Off("evt") 取消所有
// Off("evt", handler1, handler2) 取消指定函数
func (e *Emitter) Off(name string, handler ...Handler) *Emitter {
	e.lisMu.Lock()
	defer e.lisMu.Unlock()
	if len(handler) == 0 {
		delete(e.listener, name)
		return e
	}
	handlers, ok := e.listener[name]
	if !ok || len(handlers) == 0 {
		return e
	}
	nextHandlers := []Handler{}
	for _, existH := range handlers {
		rm := false
		for _, offH := range handler {
			if reflect.ValueOf(existH).Pointer() == reflect.ValueOf(offH).Pointer() {
				rm = true
			}
		}
		if !rm {
			nextHandlers = append(nextHandlers, existH)
		}
	}
	e.listener[name] = nextHandlers
	return e
}

// Emit 分发事件
func (e *Emitter) Emit(name string, v interface{}) *Emitter {
	e.lisMu.RLock()
	handlers, ok := e.listener[name]
	e.lisMu.RUnlock()
	if !ok || len(handlers) == 0 {
		return e
	}
	for _, h := range handlers {
		h(v)
	}
	return e
}

// 全局默认的监听器
var emit = &Emitter{}

func On(name string, handler ...Handler) *Emitter {
	return emit.On(name, handler...)
}

func Off(name string, handler ...Handler) *Emitter {
	return emit.Off(name, handler...)
}

func Emit(name string, v interface{}) *Emitter {
	return emit.Emit(name, v)
}
