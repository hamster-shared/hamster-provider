package chain

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

// ResourceInfo resource information
type ResourceInfo struct {
	PeerId     string    `json:"peerId"`
	Cpu        uint64    `json:"cpu"`
	Memory     uint64    `json:"memory"`
	System     string    `json:"system"`
	Image      string    `json:"image"`
	CpuModel   string    `json:"cpuModel"`
	VmType     string    `json:"vmType"`
	Creator    string    `json:"creator"`
	ExpireTime time.Time `json:"expireTime"`
	User       string    `json:"user"`
	Status     int       `json:"status"`
	Price      int64     `json:"price"`
}

type RentalAgreement struct {
	Index      types.U64
	Provider   types.AccountID
	TenantInfo struct {
		AccountId types.AccountID
		PublicKey string
	}
	PeerId        string
	ResourceIndex types.U64
	Config        struct {
		Cpu      types.U64
		Memory   types.U64
		System   string
		CpuModel string
	}
	RentalInfo struct {
		RentUnitPrice types.U128
		RentDuration  types.U32
		EndOfRent     types.U32
	}
	Price         types.U128
	LockPrice     types.U128
	PenaltyAmount types.U128
	ReceiveAmount types.U128
	Start         types.U32
	End           types.U32
	Calculation   types.U32
}

type ComputingResource struct {
	Index     types.U64
	AccountId types.AccountID
	PeerId    types.Text
	Config    struct {
		Cpu      types.U64
		Memory   types.U64
		System   types.Text
		CpuModel types.Text
	}
	RentalStatistics struct {
		RentalCount    types.U32
		RentalDuration types.U32
		FaultCount     types.U32
		FaultDuration  types.U32
	}
	RentalInfo struct {
		RentUnitPrice types.U128
		RentDuration  types.U32
		EndOfRent     types.U32
	}
	Status Status
}

type ComputingOrder struct {
	Index      types.U64
	TenantInfo struct {
		AccountId types.AccountID
		PublicKey types.Text
	}
	Price          types.U128
	ResourceIndex  types.U64
	Create         types.U32
	RentDuration   types.U32
	Status         Status
	AgreementIndex types.U64
}

type Status struct {
	IsInuse   bool
	IsLocked  bool
	IsUnused  bool
	IsOffline bool
}

func (m *Status) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	fmt.Println(b)

	if err != nil {
		return err
	}

	if b == 0 {
		m.IsInuse = true
	} else if b == 1 {
		m.IsLocked = true
	} else if b == 2 {
		m.IsUnused = true
	} else if b == 3 {
		m.IsOffline = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *Status) Encode(encoder scale.Encoder) error {
	var err1 error
	if m.IsInuse {
		err1 = encoder.PushByte(0)
	} else if m.IsLocked {
		err1 = encoder.PushByte(1)
	} else if m.IsUnused {
		err1 = encoder.PushByte(2)
	} else if m.IsOffline {
		err1 = encoder.PushByte(3)
	}
	if err1 != nil {
		return err1
	}
	return nil
}
