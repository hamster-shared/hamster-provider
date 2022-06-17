package event

import (
	"fmt"
	"sync"
)

const ResourceOrder_CreateOrderSuccess = "resource.create_order.cmd"
const ResourceOrder_ReNewOrderSuccess = "resource.renew_order.cmd"
const ResourceOrder_WithdrawLockedOrderPriceSuccess = "resource.cancel_order.cmd"
const ResourceOrder_Recover = "resource.recover_order.cmd"
const ResourceOrder_TheGraph = "resource.thegraph"

type EventHandleFunc func(e string, args interface{})

type EventHandler struct {
	ID      string
	Handler EventHandleFunc
}

type IEventBus interface {
	Pub(e string, args interface{})
	Sub(e string, id string, handleFunc EventHandleFunc)
	Unsub(e string, id string)
}

type EventBus struct {
	rwmutex *sync.RWMutex
	items   map[string][]*EventHandler
}

func newEventHandler(id string, handleFunc EventHandleFunc) *EventHandler {
	return &EventHandler{
		id, handleFunc,
	}
}

func newEventBus() IEventBus {
	it := new(EventBus)
	it.init()
	return it
}

func (me *EventBus) init() {
	me.rwmutex = new(sync.RWMutex)
	me.items = make(map[string][]*EventHandler)
}

func (me *EventBus) Pub(e string, args interface{}) {
	me.rwmutex.RLock()
	defer me.rwmutex.RUnlock()

	handlers, ok := me.items[e]
	if ok {
		for _, it := range handlers {
			fmt.Printf("eventbus.Pub, event=%s, handler=%s", e, it.ID)
			it.Handler(e, args)
		}
	}
}

func (me *EventBus) Sub(e string, id string, handleFunc EventHandleFunc) {
	me.rwmutex.Lock()
	defer me.rwmutex.Unlock()

	handler := newEventHandler(id, handleFunc)
	handlers, ok := me.items[e]

	if ok {
		me.items[e] = append(handlers, handler)
	} else {
		me.items[e] = []*EventHandler{handler}
	}
}

func (me *EventBus) Unsub(e string, id string) {
	me.rwmutex.Lock()
	defer me.rwmutex.Unlock()

	handlers, ok := me.items[e]
	if ok {
		for i, it := range handlers {
			if it.ID == id {
				lastI := len(handlers) - 1
				if i != lastI {
					handlers[i], handlers[lastI] = handlers[lastI], handlers[i]
				}
				me.items[e] = handlers[:lastI]
			}
		}
	}
}

var GlobalEventBus = newEventBus()
