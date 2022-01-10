package event

type IHandler interface {
	Name() string
	HandlerEvent(e *VmRequest)
}

type AbstractHandler struct {
}

func (h *AbstractHandler) EventHandleFunc(handler IHandler) EventHandleFunc {

	return func(e string, args interface{}) {
		switch e {
		case handler.Name():
			if it, ok := args.(*VmRequest); ok {
				handler.HandlerEvent(it)
			}
			break
		}
	}
}
