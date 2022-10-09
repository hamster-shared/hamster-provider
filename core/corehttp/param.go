package corehttp

import "time"

type BootParam struct {
	Option bool `json:"option"`
}

type UnitPriceParam struct {
	UnitPrice uint64 `json:"unitPrice"`
}

type EthereumDeployParam struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Network       string    `json:"network"`       //rinkbey network or mainnet network
	LeaseTerm     int       `json:"leaseTerm"`     //
	ApplicationID uint      `json:"applicationId"` //application id
}
