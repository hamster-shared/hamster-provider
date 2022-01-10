package event

import "github.com/hamster-shared/hamster-provider/core/context"

type IEventService interface {
	Create(r *VmRequest)
	Destroy(r *VmRequest)
	Renew(r *VmRequest)
	Recover(r *VmRequest)
}

func NewEventService(coreContext context.CoreContext) IEventService {
	it := new(EventService)
	it.init(coreContext)
	return it
}

type EventService struct {
	items []*VmRequest
}

func (s *EventService) init(coreContext context.CoreContext) {

	createHandler := &CreateVmHandler{CoreContext: coreContext}
	destroyHandler := &DestroyVmHandler{CoreContext: coreContext}
	renewHandler := &RenewVmHandler{CoreContext: coreContext}
	recoverHandler := &RecoverVmHandler{CoreContext: coreContext}

	GlobalEventBus.Sub(createHandler.Name(), "createHandler", createHandler.EventHandleFunc(createHandler))
	GlobalEventBus.Sub(destroyHandler.Name(), "destroyHandler", destroyHandler.EventHandleFunc(destroyHandler))
	GlobalEventBus.Sub(renewHandler.Name(), "renewHandler", renewHandler.EventHandleFunc(renewHandler))
	GlobalEventBus.Sub(recoverHandler.Name(), "recoverHandler", recoverHandler.EventHandleFunc(recoverHandler))
}

func (s *EventService) Create(req *VmRequest) {
	GlobalEventBus.Pub(ResourceOrder_CreateOrderSuccess, req)
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
