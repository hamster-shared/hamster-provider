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
	ChangeResourceStatus(index uint64) error

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
}
