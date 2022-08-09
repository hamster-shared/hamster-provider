package chain

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

// ResourceInfo 资源信息
type ResourceInfo struct {
	PeerId        string    `json:"peerId"`
	Cpu           uint64    `json:"cpu"`
	Memory        uint64    `json:"memory"`
	System        string    `json:"system"`
	Image         string    `json:"image"`
	CpuModel      string    `json:"cpuModel"`
	VmType        string    `json:"vmType"`
	Creator       string    `json:"creator"`
	ExpireTime    time.Time `json:"expireTime"`
	User          string    `json:"user"`
	Status        int       `json:"status"`
	Price         uint64    `json:"price"`
	ResourceIndex uint64    `json:"resource_index"`
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
	Index     types.U64       `json:"index"`
	AccountId types.AccountID `json:"accountId"`
	PeerId    types.Text      `json:"peerId"`
	Config    struct {
		Cpu      types.U64  `json:"cpu"`
		Memory   types.U64  `json:"memory"`
		System   types.Text `json:"system"`
		CpuModel types.Text `json:"cpuModel"`
	} `json:"config"`
	RentalStatistics struct {
		RentalCount    types.U32 `json:"rentalCount"`
		RentalDuration types.U32 `json:"rentalDuration"`
		FaultCount     types.U32 `json:"faultCount"`
		FaultDuration  types.U32 `json:"faultDuration"`
	} `json:"rentalStatistics"`
	RentalInfo struct {
		RentDuration types.U32 `json:"rentDuration"`
		EndOfRent    types.U32 `json:"endOfRent"`
	} `json:"rentalInfo"`
	Status Status `json:"status"`
}

type Status struct {
	IsInuse   bool `json:"isInuse"`
	IsLocked  bool `json:"isLocked"`
	IsUnused  bool `json:"isUnused"`
	IsOffline bool `json:"isOffline"`
}

type StakingAmount struct {
	Amount       types.U128
	ActiveAmount types.U128
	LockAmount   types.U128
}

type AccountInfoCustom struct {
	Nonce       types.U32
	Consumers   types.U32
	Providers   types.U32
	Sufficients types.U32
	Data        struct {
		Free       types.U128
		Reserved   types.U128
		MiscFrozen types.U128
		FreeFrozen types.U128
	}
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

type ComputingOrder struct {
	Index      types.U64
	TenantInfo struct {
		AccountId types.AccountID
		PublicKey types.Text
	}
	ResourceIndex  types.U64
	Create         types.U32
	RentDuration   types.U32
	Time           Duration
	Status         OrderStatus
	AgreementIndex types.OptionU64
}

type Duration struct {
	Secs  types.U64
	Nanos types.U32
}

type OrderStatus struct {
	IsPending  bool
	IsFinished bool
	IsCanceled bool
}

type AccountInfo struct {
	Address string
	Amount  types.U128
}

func (m *OrderStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	fmt.Println(b)

	if err != nil {
		return err
	}

	if b == 0 {
		m.IsPending = true
	} else if b == 1 {
		m.IsFinished = true
	} else if b == 2 {
		m.IsCanceled = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *OrderStatus) Encode(encoder scale.Encoder) error {
	var err1 error
	if m.IsPending {
		err1 = encoder.PushByte(0)
	} else if m.IsFinished {
		err1 = encoder.PushByte(1)
	} else if m.IsCanceled {
		err1 = encoder.PushByte(2)
	}
	if err1 != nil {
		return err1
	}
	return nil
}

type MarketUser struct {
	StakedAmount types.U128
}

const Provider_MarketUserStatus = types.U8(0)

type MarketIncome struct {
	LastEraIndex types.U32
	TotalIncome  types.U128
}
