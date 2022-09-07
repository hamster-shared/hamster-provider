package chain

import (
	"time"
)

// ReportClient data reporting interface
type ReportClient interface {
	// RegisterResource resource registration
	RegisterResource(ResourceInfo) error
	// RemoveResource resource deletion
	RemoveResource(index uint64) error
	// ModifyResourcePrice modify resource unit price
	ModifyResourcePrice(index uint64, unitPrice int64) error
	// ChangeResourceStatus modify resource status to unused
	ChangeResourceStatus(resourceIndex uint64) error

	// AddResourceDuration add resource rental time
	AddResourceDuration(index uint64, duration int) error

	// Heartbeat protocol heartbeat report
	Heartbeat(agreementindex uint64) error

	// OrderExec
	OrderExec(orderIndex uint64) error

	LoadKeyFromChain() ([]string, error)

	CalculateInstanceOverdue(orderIndex uint64) time.Duration

	GetAgreementIndex(orderIndex uint64) (uint64, error)

	//GetResource get vm resource
	GetResource(resourceIndex uint64) (*ComputingResource, error)

	CalculateResourceOverdue(expireBlock uint64) (time.Duration, error)

	ReceiveIncome() error

	GetAccountInfo() (*AccountInfo, error)

	GetStakingInfo() (*StakingAmount, error)

	GetMarketStackInfo() (*StakingAmount, error)

	StakingAmount(unitPrice int64) error

	WithdrawStakingAmount(unitPrice int64) error

	ReceiveIncomeJudge() bool

	GetGatewayNodes() ([]string, error)

	ProcessApplyFreeResource(index uint64, peerId string) error

	ReleaseApplyFreeResource(index uint64) error

	GetReward() (*MarketIncome, error)

	PayoutReward() error
}
