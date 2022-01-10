package event

type IVmEventHandler interface {
	EventHandleFunc(e string, args interface{})
}
