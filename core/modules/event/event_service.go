package event

type IEventService interface {
	Create(r *VmRequest)
	Destroy(r *VmRequest)
	Renew(r *VmRequest)
	Recover(r *VmRequest)
}

func NewEventService(coreContext *EventContext) IEventService {
	it := new(EventService)
	it.init(coreContext)
	return it
}

type EventService struct {
	items []*VmRequest
}

func (s *EventService) init(coreContext *EventContext) {

	createHandler := &CreateVmHandler{CoreContext: coreContext}
	destroyHandler := &DestroyVmHandler{CoreContext: coreContext}
	renewHandler := &RenewVmHandler{CoreContext: coreContext}
	recoverHandler := &RecoverVmHandler{CoreContext: coreContext}
	thegraphHandler := &TheGraphHandler{CoreContext: coreContext}

	GlobalEventBus.Sub(createHandler.Name(), "createHandler", createHandler.EventHandleFunc(createHandler))
	GlobalEventBus.Sub(destroyHandler.Name(), "destroyHandler", destroyHandler.EventHandleFunc(destroyHandler))
	GlobalEventBus.Sub(renewHandler.Name(), "renewHandler", renewHandler.EventHandleFunc(renewHandler))
	GlobalEventBus.Sub(recoverHandler.Name(), "recoverHandler", recoverHandler.EventHandleFunc(recoverHandler))
	GlobalEventBus.Sub(thegraphHandler.Name(), "thegraphHandler", thegraphHandler.EventHandleFunc(thegraphHandler))
}

func (s *EventService) Create(req *VmRequest) {
	//if req.Tag == OPCreatedVm {
	//	GlobalEventBus.Pub(ResourceOrder_CreateOrderSuccess, req)
	//} else if req.Tag == OPFreeResourceApply {
	GlobalEventBus.Pub(ResourceOrder_TheGraph, req)
	//}

}

func (s *EventService) Destroy(req *VmRequest) {
	GlobalEventBus.Pub(ResourceOrder_WithdrawLockedOrderPriceSuccess, req)
}

func (s *EventService) Renew(req *VmRequest) {
	GlobalEventBus.Pub(ResourceOrder_ReNewOrderSuccess, req)
}

func (s *EventService) Recover(req *VmRequest) {
	GlobalEventBus.Pub(ResourceOrder_Recover, req)
}
