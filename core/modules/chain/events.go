package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type EventProviderRegisterResourceSuccess struct {
	Phase            types.Phase
	AccountId        types.AccountID
	Index            types.U64
	PeerId           string
	Cpu              types.U64
	Memory           types.U64
	System           string
	CpuModel         string
	PriceHour        types.U128
	RentDurationHour types.U32
	Topics           []types.Hash
}

type EventResourceOrderCreateOrderSuccess struct {
	Phase         types.Phase
	AccountId     types.AccountID
	OrderIndex    types.U64
	ResourceIndex types.U64
	Duration      types.U32
	DeployType    types.U32
	PublicKey     string
	Topics        []types.Hash
}

type EventResourceOrderOrderExecSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	OrderIndex     types.U64
	ResourceIndex  types.U64
	AgreementIndex types.U64
	Topics         []types.Hash
}

type EventResourceOrderReNewOrderSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	OrderIndex     types.U64
	ResourceIndex  types.U64
	AgreementIndex types.U64
	Topics         []types.Hash
}

type EventResourceOrderWithdrawLockedOrderPriceSuccess struct {
	Phase      types.Phase
	AccountId  types.AccountID
	OrderIndex types.U64
	Topics     []types.Hash
}

type EventMarketMoney struct {
	Phase  types.Phase
	Money  types.U128
	Topics []types.Hash
}

type EventCancelAgreementSuccess struct {
	Phase          types.Phase
	AccountId      types.AccountID
	AgreementIndex types.U64
	OrderIndex     types.U64
	Topics         []types.Hash
}

type MyEventRecords struct {
	types.EventRecords
	Provider_RegisterResourceSuccess              []EventProviderRegisterResourceSuccess //nolint:stylecheck,golint
	ResourceOrder_CreateOrderSuccess              []EventResourceOrderCreateOrderSuccess //nolint:stylecheck,golint
	ResourceOrder_OrderExecSuccess                []EventResourceOrderOrderExecSuccess
	ResourceOrder_ReNewOrderSuccess               []EventResourceOrderReNewOrderSuccess
	ResourceOrder_WithdrawLockedOrderPriceSuccess []EventResourceOrderWithdrawLockedOrderPriceSuccess
	ResourceOrder_CancelAgreementSuccess          []EventCancelAgreementSuccess
}
